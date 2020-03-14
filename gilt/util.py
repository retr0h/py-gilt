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

from __future__ import print_function

import contextlib
import errno
import os
import sh
import shutil

import colorama

colorama.init(autoreset=True)


def print_info(msg):
    """Print the given message to STDOUT. """
    print(msg)


def print_warn(msg):
    """Print the given message to STDOUT in YELLOW. """
    print("{}{}".format(colorama.Fore.YELLOW, msg))


def run_command(cmd, debug=False):
    """Execute the given command and return None.

    :param cmd: A `sh.Command` object to execute.
    :param debug: An optional bool to toggle debug output.
    :return: None
    """
    if debug:
        msg = "  PWD: {}".format(os.getcwd())
        print_warn(msg)
        msg = "  COMMAND: {}".format(cmd)
        print_warn(msg)
    cmd()


def build_sh_cmd(cmd, cwd=None):
    """Build a `sh.Command` from a string.

    :param cmd: String with the command to convert.
    :param cwd: Optional path to use as working directory.
    :return: `sh.Command`
    """
    args = cmd.split()
    return getattr(sh, args[0]).bake(_cwd=cwd, *args[1:])


@contextlib.contextmanager
def saved_cwd():
    """Context manager to restore previous working directory. """
    saved = os.getcwd()
    try:
        yield
    finally:
        os.chdir(saved)


def copy(src, dst):
    """Handle the copying of a file or directory.

    The destination basedir _must_ exist.

    :param src: A string containing the path of the source to copy.  If the
     source ends with a '/', will become a recursive directory copy of source.
    :param dst: A string containing the path to the destination.  If the
     destination ends with a '/', will copy into the target directory.
    :return: None
    """
    try:
        shutil.copytree(src, dst)
    except OSError as exc:
        if exc.errno == errno.ENOTDIR:
            shutil.copy(src, dst)
        else:
            raise
