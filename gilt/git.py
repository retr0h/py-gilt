# Copyright (c) 2016 Cisco Systems, Inc.
#
#  Permission is hereby granted, free of charge, to any person obtaining a copy
#  of this software and associated documentation files (the "Software"), to
#  deal in the Software without restriction, including without limitation the
#  rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
#  sell copies of the Software, and to permit persons to whom the Software is
#  furnished to do so, subject to the following conditions:
#
#  The above copyright notice and this permission notice shall be included in
#  all copies or substantial portions of the Software.
#
#  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
#  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
#  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
#  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
#  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
#  FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
#  DEALINGS IN THE SOFTWARE.

import glob
import os
import shutil

import sh
from gilt import util


def clone(name, repository, destination, debug=False):
    """Clone the specified repository into a temporary directory and return None.

    :param name: A string containing the name of the repository being cloned.
    :param repository: A string containing the repository to clone.
    :param destination: A string containing the directory to clone the
     repository into.
    :param debug: An optional bool to toggle debug output.
    :return: None
    """
    msg = "  - cloning {} to {}".format(name, destination)
    util.print_info(msg)
    cmd = sh.git.bake("clone", repository, destination)
    util.run_command(cmd, debug=debug)


def sync(name, destination, version, debug=False):
    os.chdir(destination)
    msg = "  - syncing {} with {} of origin".format(name, version)
    util.print_info(msg)
    _get_version(version, clean=False, debug=debug)


def remote_add(destination, name, url, debug=False):
    os.chdir(destination)
    msg = "  - adding {} remote with url {}".format(name, url)
    util.print_info(msg)
    try:
        cmd = sh.git.bake("remote", "remove", name)
        util.run_command(cmd, debug=debug)
    except sh.ErrorReturnCode:
        pass
    cmd = sh.git.bake("remote", "add", name, url)
    util.run_command(cmd, debug=debug)
    cmd = sh.git.bake("fetch", name)
    util.run_command(cmd, debug=debug)


def extract(repository, destination, version, debug=False):
    """Extract the specified repository/version into the directory and return None.

    :param repository: A string containing the path to the repository to be
     extracted.
    :param destination: A string containing the directory to clone the
     repository into.  Relative to the directory ``gilt`` is running
     in. Must end with a '/'.
    :param version: A string containing the branch/tag/sha to be exported.
    :param debug: An optional bool to toggle debug output.
    :return: None
    """
    with util.saved_cwd():
        if os.path.isdir(destination):
            shutil.rmtree(destination)

        os.chdir(repository)
        _get_version(version, debug=debug)
        cmd = sh.git.bake(
            "checkout-index", force=True, all=True, prefix=destination
        )
        util.run_command(cmd, debug=debug)
        msg = "  - extracting ({}) {} to {}".format(
            version, repository, destination
        )
        util.print_info(msg)


def overlay(repository, files, version, debug=False):
    """Overlay files from repository/version into the directory and return None.

    :param repository: A string containing the path to the repository to be
     extracted.
    :param files: A list of `FileConfig` objects.
    :param version: A string containing the branch/tag/sha to be exported.
    :param debug: An optional bool to toggle debug output.
    :return: None
    """
    with util.saved_cwd():
        os.chdir(repository)
        _get_version(version, debug=debug)

        for fc in files:
            if "*" in fc.src:
                for filename in glob.glob(fc.src):
                    util.copy(filename, fc.dst)
                    msg = "  - copied ({}) {} to {}".format(
                        version, filename, fc.dst
                    )
                    util.print_info(msg)
            else:
                if os.path.isdir(fc.dst) and os.path.isdir(fc.src):
                    shutil.rmtree(fc.dst)
                util.copy(fc.src, fc.dst)
                msg = "  - copied ({}) {} to {}".format(
                    version, fc.src, fc.dst
                )
                util.print_info(msg)


def _get_version(version, clean=True, debug=False):
    """Handle switching to the specified version and return None.

    1. Fetch the origin.
    2. Checkout the specified version.
    3. Clean the repository before we begin.
    4. Pull the origin when a branch; _not_ a commit id.

    :param version: A string containing the branch/tag/sha to be exported.
    :param clean: An optional bool to toggle running `git clean` before `git pull`
    :param debug: An optional bool to toggle debug output.
    :return: None
    """
    if not any(
        (
            _has_branch(version, debug),
            _has_tag(version, debug),
            _has_commit(version, debug),
        )
    ):
        cmd = sh.git.bake("fetch")
        util.run_command(cmd, debug=debug)
    cmd = sh.git.bake("checkout", version)
    util.run_command(cmd, debug=debug)
    if clean:
        cmd = sh.git.bake("clean", "-d", "-x", "-f")
        util.run_command(cmd, debug=debug)
    if _has_branch(version, debug):
        try:
            cmd = sh.git.bake("pull", rebase=True, ff_only=True)
            util.run_command(cmd, debug=debug)
        except sh.ErrorReturnCode:
            msg = "  - pulling failed, local changes exist?"
            util.print_warn(msg)


def _has_commit(version, debug=False):
    """Determine a version is a local git commit sha or not.

    :param version: A string containing the branch/tag/sha to be determined.
    :param debug: An optional bool to toggle debug output.
    :return: bool
    """
    if _has_tag(version, debug) or _has_branch(version, debug):
        return False
    cmd = sh.git.bake("cat-file", "-e", version)
    try:
        util.run_command(cmd, debug=debug)
        return True
    except sh.ErrorReturnCode:
        return False


def _has_tag(version, debug=False):
    """Determine a version is a local git tag name or not.

    :param version: A string containing the branch/tag/sha to be determined.
    :param debug: An optional bool to toggle debug output.
    :return: bool
    """
    cmd = sh.git.bake(
        "show-ref", "--verify", "--quiet", "refs/tags/{}".format(version)
    )
    try:
        util.run_command(cmd, debug=debug)
        return True
    except sh.ErrorReturnCode:
        return False


def _has_branch(version, debug=False):
    """Determine a version is a local git branch name or not.

    :param version: A string containing the branch/tag/sha to be determined.
    :param debug: An optional bool to toggle debug output.
    :return: bool
    """
    cmd = sh.git.bake(
        "show-ref", "--verify", "--quiet", "refs/heads/{}".format(version)
    )
    try:
        util.run_command(cmd, debug=debug)
        return True
    except sh.ErrorReturnCode:
        return False
