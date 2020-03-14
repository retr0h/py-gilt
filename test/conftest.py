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
import random
import shutil
import string

import pytest

pytest_plugins = ["helpers_namespace"]


def pytest_addoption(parser):
    parser.addoption("--runslow", action="store_true", help="run slow tests")


def random_string(len=5):
    return "".join(random.choice(string.ascii_uppercase) for _ in range(len))


@pytest.fixture()
def temp_dir(tmpdir, request):
    _cwd = os.getcwd()
    d = tmpdir.mkdir(random_string())
    os.chdir(d.strpath)

    def cleanup():
        shutil.rmtree(d.strpath)
        os.chdir(_cwd)

    request.addfinalizer(cleanup)

    return d


@pytest.fixture()
def gilt_config_file(temp_dir, request):
    fixture = request.param
    d = temp_dir
    c = d.join(os.extsep.join(("gilt", "yml")))
    c.write(request.getfixturevalue(fixture))

    return c.strpath


@pytest.fixture()
def gilt_data():
    return [
        {
            "git": "https://github.com/retr0h/ansible-etcd.git",
            "version": "master",
            "dst": "roles/retr0h.ansible-etcd/",
        },
        {
            "git": "https://github.com/lorin/openstack-ansible-modules.git",
            "version": "master",
            "files": [{"src": "*_manage", "dst": "library/"}],
        },
    ]


@pytest.helpers.register
def os_split(s):
    rest, tail = os.path.split(s)
    if rest in ("", os.path.sep):
        return (tail,)

    return os_split(rest) + (tail,)
