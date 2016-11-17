# vim: tabstop=4 shiftwidth=4 softtabstop=4

# Copyright (c) 2016 Cisco Systems, Inc.
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

import os

import sh

from gilt import util


def test_print_info(capsys):
    util.print_info('foo')

    result, _ = capsys.readouterr()
    assert 'foo\n' == result


def test_print_warn(capsys):
    util.print_warn('foo')

    result, _ = capsys.readouterr()
    assert 'foo' in result


def test_run_command(capsys):
    cmd = sh.git.bake(version=True)
    util.run_command(cmd)

    result, _ = capsys.readouterr()
    assert '' == result


def test_run_command_with_debug(temp_dir, capsys):
    cmd = sh.git.bake(version=True)
    util.run_command(cmd, debug=True)

    result, _ = capsys.readouterr()
    x = 'COMMAND: {} --version'.format(sh.git)
    assert x in result
    x = 'PWD: {}'.format(temp_dir)
    assert x in result


def test_saved_cwd_contextmanager(temp_dir):
    workdir = os.path.join(temp_dir.strpath, 'workdir')

    os.mkdir(workdir)

    with util.saved_cwd():
        os.chdir(workdir)
        assert workdir == os.getcwd()

    assert temp_dir.strpath == os.getcwd()


def test_copy_file(temp_dir):
    dst_dir = os.path.join(temp_dir.strpath, 'dst')

    os.mkdir(dst_dir)

    src = os.path.join(temp_dir.strpath, 'foo')
    open(src, 'a').close()

    util.copy(src, dst_dir)

    dst = os.path.join(dst_dir, 'foo')
    assert os.path.exists(dst)


def test_copy_dir(temp_dir):
    src_dir = os.path.join(temp_dir.strpath, 'src')
    dst_dir = os.path.join(temp_dir.strpath, 'dst')

    os.mkdir(src_dir)
    os.mkdir(dst_dir)

    d = os.path.join(dst_dir, 'src')
    util.copy(src_dir, d)

    assert os.path.exists(d)
