# vim: tabstop=4 shiftwidth=4 softtabstop=4

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

import fasteners
import glob
import os
import shutil

import sh

from gilt import util


def _git_env(repository):
    """
    Git environment variables for specific repository.

    :param repository: path to a git repository.
    :return: dict of envionment variables.
    """
    return {
        "GIT_WORK_TREE": repository,
        "GIT_DIR": os.path.join(repository, ".git"),
        "GIT_OBJECT_DIRECTORY": os.path.join(repository, ".git", "objects")
    }


def do_overlay(config, debug=False):
    """
    Actual methods the overlays git based repo.

    :param config: A dict contains configuration of the git resource.
    :param debug: An optional bool to toggle debug output.
    :return: None
    """
    with fasteners.InterProcessLock(config.lock_file):
        with util.named_lock(config.name):
            if not os.path.exists(config.src):
                clone(config.name, config.git, config.src, debug=debug)
            if config.dst:
                extract(config.src, config.dst, config.version, debug=debug)
            else:
                overlay(config.src, config.files, config.version, debug=debug)


def clone(name, repository, destination, debug=False):
    """
    Clone the specified repository into a temporary directory and return None.

    :param name: A string containing the name of the repository being cloned.
    :param repository: A string containing the repository to clone.
    :param destination: A string containing the directory to clone the
     repository into.
    :param debug: An optional bool to toggle debug output.
    :return: None
    """
    msg = '  - cloning {} to {}'.format(name, destination)
    util.print_info(msg)
    cmd = sh.git.bake('clone', repository, destination)
    util.run_command(cmd, debug=debug)


def extract(repository, destination, version, debug=False):
    """
    Extract the specified repository/version into the given directory and
    return None.

    :param repository: A string containing the path to the repository to be
     extracted.
    :param destination: A string containing the directory to clone the
     repository into.  Relative to the directory ``gilt`` is running
     in. Must end with a '/'.
    :param version: A string containing the branch/tag/sha to be exported.
    :param debug: An optional bool to toggle debug output.
    :return: None
    """
    git_env = _git_env(repository)
    _get_version(version, debug, git_env)
    cmd = sh.git.bake(
        'checkout-index',
        force=True,
        all=True,
        prefix=destination,
        _env=git_env)
    util.run_command(cmd, debug=debug)
    msg = '  - extracting ({}) {} to {}'.format(version, repository,
                                                destination)
    util.print_info(msg)


def overlay(repository, files, version, debug=False):
    """
    Overlay files from the specified repository/version into the given
    directory and return None.

    :param repository: A string containing the path to the repository to be
     extracted.
    :param files: A list of `FileConfig` objects.
    :param version: A string containing the branch/tag/sha to be exported.
    :param debug: An optional bool to toggle debug output.
    :return: None
    """
    git_env = _git_env(repository)
    _get_version(version, debug, git_env)

    for fc in files:
        if '*' in fc.src:
            for filename in glob.glob(fc.src):
                util.copy(filename, fc.dst)
                msg = '  - copied ({}) {} to {}'.format(
                    version, filename, fc.dst)
                util.print_info(msg)
        else:
            if os.path.isdir(fc.dst) and os.path.isdir(fc.src):
                shutil.rmtree(fc.dst)
            util.copy(fc.src, fc.dst)
            msg = '  - copied ({}) {} to {}'.format(version, fc.src, fc.dst)
            util.print_info(msg)


def _get_version(version, debug=False, env=None):
    """
    Handle switching to the specified version and return None.

    1. Fetch the origin.
    2. Checkout the specified version.
    3. Clean the repository before we begin.
    4. Pull the origin when a branch; _not_ a commit id.

    :param version: A string containing the branch/tag/sha to be exported.
    :param debug: An optional bool to toggle debug output.
    :return: None
    """
    if not (_has_branch(version, debug, env) or
            _has_tag(version, debug, env) or _has_commit(version, debug, env)):
        cmd = sh.git.bake('fetch', _env=env)
        util.run_command(cmd, debug=debug)
    cmd = sh.git.bake('checkout', version, _env=env)
    util.run_command(cmd, debug=debug)
    cmd = sh.git.bake('clean', '-d', '-x', '-f', _env=env)
    util.run_command(cmd, debug=debug)
    if _has_branch(version, debug):
        cmd = sh.git.bake('pull', rebase=True, ff_only=True, _env=env)
        util.run_command(cmd, debug=debug)


def _has_commit(version, debug=False, env=None):
    """
    Determine a version is a local git commit sha or not.

    :param version: A string containing the branch/tag/sha to be determined.
    :param debug: An optional bool to toggle debug output.
    :return: bool
    """
    if _has_branch(version, debug, env) or _has_tag(version, debug, env):
        return False
    cmd = sh.git.bake('cat-file', '-e', version, _env=env)
    try:
        util.run_command(cmd, debug=debug)
        return True
    except sh.ErrorReturnCode:
        return False


def _has_tag(version, debug=False, env=None):
    """
    Determine a version is a local git tag name or not.

    :param version: A string containing the branch/tag/sha to be determined.
    :param debug: An optional bool to toggle debug output.
    :return: bool
    """
    cmd = sh.git.bake(
        'show-ref',
        '--verify',
        '--quiet',
        "refs/tags/{}".format(version),
        _env=env)
    try:
        util.run_command(cmd, debug=debug)
        return True
    except sh.ErrorReturnCode:
        return False


def _has_branch(version, debug=False, env=None):
    """
    Determine a version is a local git branch name or not.

    :param version: A string containing the branch/tag/sha to be determined.
    :param debug: An optional bool to toggle debug output.
    :return: bool
    """
    cmd = sh.git.bake(
        'show-ref',
        '--verify',
        '--quiet',
        "refs/heads/{}".format(version),
        _env=env)
    try:
        util.run_command(cmd, debug=debug)
        return True
    except sh.ErrorReturnCode:
        return False
