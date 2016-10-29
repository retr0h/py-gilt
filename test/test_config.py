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

import pytest

from gilt import config


@pytest.mark.parametrize(
    'gilt_config_file', ['gilt_data'], indirect=['gilt_config_file'])
def test_config(gilt_config_file):
    result = config.config(gilt_config_file)
    os_split = pytest.helpers.os_split

    r = result[0]
    assert 'https://github.com/retr0h/ansible-etcd.git' == r.git
    assert 'master' == r.version
    assert 'retr0h.ansible-etcd' == r.name
    assert ('.gilt', 'clone', 'retr0h.ansible-etcd') == os_split(r.src)[-3:]
    assert ('roles', 'retr0h.ansible-etcd', '') == os_split(r.dst)[-3:]
    assert [] == r.files

    r = result[1]
    assert 'https://github.com/lorin/openstack-ansible-modules.git' == r.git
    assert 'master' == r.version
    assert 'lorin.openstack-ansible-modules' == r.name
    assert 'lorin.openstack-ansible-modules' == os_split(r.src)[-1]
    assert r.dst is None
    f = r.files[0]
    x = ('.gilt', 'clone', 'lorin.openstack-ansible-modules', '*_manage')
    assert x == os_split(f.src)[-4:]
    assert ('library', '') == os_split(f.dst)[-2:]


@pytest.mark.parametrize(
    'gilt_config_file', ['gilt_data'], indirect=['gilt_config_file'])
def test_get_config_generator(gilt_config_file):
    result = [i for i in config._get_config_generator(gilt_config_file)]

    assert isinstance(result, list)
    assert isinstance(result[0], dict)


def test_get_files_generator(temp_dir):
    files_list = [{'src': 'foo', 'dest': 'bar'}]
    result = [i for i in config._get_files_generator('/tmp/dir', files_list)]

    assert isinstance(result, list)
    assert isinstance(result[0], dict)


@pytest.mark.parametrize(
    'gilt_config_file', ['gilt_data'], indirect=['gilt_config_file'])
def test_get_config(gilt_config_file):
    result = config._get_config(gilt_config_file)

    assert isinstance(result, list)
    assert isinstance(result[0], dict)


@pytest.fixture()
def invalid_gilt_data():
    return '{'


@pytest.mark.parametrize(
    'gilt_config_file', ['invalid_gilt_data'], indirect=['gilt_config_file'])
def test_get_config_handles_parse_error(gilt_config_file):
    with pytest.raises(config.ParseError):
        config._get_config(gilt_config_file)


def test_get_dst_dir(temp_dir):
    os.chdir(temp_dir.strpath)
    result = config._get_dst_dir({'dst': 'role/foo'})

    assert os.path.join(temp_dir.strpath, 'role', 'foo', '') == result


def test_get_dst_dir_with_missing_dst(temp_dir):
    result = config._get_dst_dir({'foo': 'bar'})

    assert result is None


def test_get_clone_dir():
    parts = pytest.helpers.os_split(config._get_clone_dir())

    assert ('.gilt', 'clone') == parts[-2:]
