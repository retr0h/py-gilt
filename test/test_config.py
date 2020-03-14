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
    "gilt_config_file", ["gilt_data"], indirect=["gilt_config_file"]
)
def test_config(gilt_config_file):
    result = config.config(gilt_config_file)
    os_split = pytest.helpers.os_split
    gilt_root = os.path.basename(config.BASE_WORKING_DIR)

    r = result[0]
    assert 'https://github.com/retr0h/ansible-etcd.git' == r.git
    assert 'master' == r.version
    assert 'retr0h.ansible-etcd' == r.name
    assert (gilt_root, 'clone', 'github.com',
            'retr0h.ansible-etcd') == os_split(r.src)[-4:]
    assert (gilt_root, 'lock', 'github.com',
            'retr0h.ansible-etcd') == os_split(r.lock_file)[-4:]
    assert ('roles', 'retr0h.ansible-etcd', '') == os_split(r.dst)[-3:]
    assert [] == r.files

    r = result[1]
    assert "https://github.com/lorin/openstack-ansible-modules.git" == r.git
    assert "master" == r.version
    assert "lorin.openstack-ansible-modules" == r.name
    assert "lorin.openstack-ansible-modules" == os_split(r.src)[-1]
    assert r.dst is None

    f = r.files[0]
    x = (gilt_root, 'clone', 'github.com',
         'lorin.openstack-ansible-modules', '*_manage')
    assert x == os_split(f.src)[-5:]
    assert ('library', '') == os_split(f.dst)[-2:]


gilt_repos = [
    (
        "example.com:/owner/repo.git",
        {"name": "repo", "owner": "owner", "hostname": "example.com"},
    ),
    (
        "example.com:owner/repo.git",
        {"name": "repo", "owner": "owner", "hostname": "example.com"},
    ),
    (
        "example.com:repo.git",
        {"name": "repo", "owner": "", "hostname": "example.com"},
    ),
    (
        "git+https://example.com/owner/repo.git",
        {"name": "repo", "owner": "owner", "hostname": "example.com"},
    ),
    (
        "git+ssh://example.com/owner/repo.git",
        {"name": "repo", "owner": "owner", "hostname": "example.com"},
    ),
    (
        "git+ssh://git@git.example.com/~philip.sd6/test.repo.git",
        {
            "name": "test.repo",
            "owner": "philip.sd6",
            "hostname": "git.example.com",
        },
    ),
    (
        "git://example.com/owner/repo.git",
        {"name": "repo", "owner": "owner", "hostname": "example.com"},
    ),
    (
        "http://example.com/owner/repo",
        {"name": "repo", "owner": "owner", "hostname": "example.com"},
    ),
    (
        "http://example.com/owner/repo.git",
        {"name": "repo", "owner": "owner", "hostname": "example.com"},
    ),
    (
        "http://example.com/repo",
        {"name": "repo", "owner": "", "hostname": "example.com"},
    ),
    (
        "http://example.com:29418/owner/repo.git",
        {"name": "repo", "owner": "owner", "hostname": "example.com"},
    ),
    (
        "http://user@example.com/user/repo",
        {"name": "repo", "owner": "user", "hostname": "example.com"},
    ),
    (
        "http://user@example.com:29418/user/repo",
        {"name": "repo", "owner": "user", "hostname": "example.com"},
    ),
    (
        "https://example.com/git/scm/project/my-sample.repo.git",
        {
            "name": "my-sample.repo",
            "owner": "project",
            "hostname": "example.com",
        },
    ),
    (
        "https://example.com/git/scm/~philip.sd6/my-sample.repo.git",
        {
            "name": "my-sample.repo",
            "owner": "philip.sd6",
            "hostname": "example.com",
        },
    ),
    (
        "https://example.com/owner/repo",
        {"name": "repo", "owner": "owner", "hostname": "example.com"},
    ),
    (
        "https://example.com/owner/repo.git",
        {"name": "repo", "owner": "owner", "hostname": "example.com"},
    ),
    (
        "https://example.com/repo",
        {"name": "repo", "owner": "", "hostname": "example.com"},
    ),
    (
        "https://example.com:29418/owner/repo.git",
        {"name": "repo", "owner": "owner", "hostname": "example.com"},
    ),
    (
        "https://github.com/sphinx-doc/sphinx.git",
        {"name": "sphinx", "owner": "sphinx-doc", "hostname": "github.com"},
    ),
    (
        "https://github.com/tterranigma/Stouts.openvpn",
        {
            "name": "Stouts.openvpn",
            "owner": "tterranigma",
            "hostname": "github.com",
        },
    ),
    (
        "https://github.com/tterranigma/Stouts.openvpn.git",
        {
            "name": "Stouts.openvpn",
            "owner": "tterranigma",
            "hostname": "github.com",
        },
    ),
    (
        "https://user@example.com/user/repo",
        {"name": "repo", "owner": "user", "hostname": "example.com"},
    ),
    (
        "https://user@example.com:29418/user/repo",
        {"name": "repo", "owner": "user", "hostname": "example.com"},
    ),
    (
        "rsync://example.com/owner/repo.git",
        {"name": "repo", "owner": "owner", "hostname": "example.com"},
    ),
    (
        "ssh://example.com/owner/repo.git",
        {"name": "repo", "owner": "owner", "hostname": "example.com"},
    ),
    (
        "ssh://example.com:29418/owner/repo.git",
        {"name": "repo", "owner": "owner", "hostname": "example.com"},
    ),
    (
        "ssh://user@example.com/owner/repo.git",
        {"name": "repo", "owner": "owner", "hostname": "example.com"},
    ),
    (
        "ssh://user@example.com:29418/owner/repo.git",
        {"name": "repo", "owner": "owner", "hostname": "example.com"},
    ),
    (
        "user@example.com:/owner/repo.git",
        {"name": "repo", "owner": "owner", "hostname": "example.com"},
    ),
    (
        "user@example.com:owner/repo.git",
        {"name": "repo", "owner": "owner", "hostname": "example.com"},
    ),
    (
        "user@example.com:repo.git",
        {"name": "repo", "owner": "", "hostname": "example.com"},
    ),
    (
        "user@foo-example.com:owner/repo.git",
        {"name": "repo", "owner": "owner", "hostname": "foo-example.com"},
    ),
]


@pytest.mark.parametrize('uri,expected', gilt_repos)
def test_config_repo(uri, expected):
    parsedrepo = config._parse_repo_uri(uri)
    assert parsedrepo.hostname == expected['hostname']
    assert parsedrepo.owner == expected['owner']
    assert parsedrepo.name == expected['name']


@pytest.fixture()
def missing_git_key_data():
    return [{"foo": "https://github.com/retr0h/ansible-etcd.git"}]


@pytest.mark.parametrize(
    "gilt_config_file", ["missing_git_key_data"], indirect=["gilt_config_file"]
)
def test_config_missing_git_key(gilt_config_file):
    with pytest.raises(KeyError):
        config.config(gilt_config_file)


@pytest.fixture()
def missing_version_key_data():
    return [
        {"git": "https://github.com/retr0h/ansible-etcd.git", "foo": "master"}
    ]


@pytest.mark.parametrize(
    "gilt_config_file",
    ["missing_version_key_data"],
    indirect=["gilt_config_file"],
)
def test_config_missing_version_key(gilt_config_file):
    with pytest.raises(KeyError):
        config.config(gilt_config_file)


@pytest.fixture()
def missing_dst_key_data():
    return [
        {
            "git": "https://github.com/retr0h/ansible-etcd.git",
            "version": "master",
            "foo": "roles/retr0h.ansible-etcd/",
        }
    ]


@pytest.mark.parametrize(
    "gilt_config_file", ["missing_dst_key_data"], indirect=["gilt_config_file"]
)
def test_config_missing_dst_key(gilt_config_file):
    with pytest.raises(KeyError):
        config.config(gilt_config_file)


@pytest.fixture()
def missing_files_src_key_data():
    return [
        {
            "git": "https://github.com/lorin/openstack-ansible-modules.git",
            "version": "master",
            "files": [{"foo": "*_manage", "dst": "library/"}],
        }
    ]


@pytest.mark.parametrize(
    "gilt_config_file",
    ["missing_files_src_key_data"],
    indirect=["gilt_config_file"],
)
def test_config_missing_files_src_key(gilt_config_file):
    with pytest.raises(KeyError):
        config.config(gilt_config_file)


@pytest.fixture()
def missing_files_dst_key_data():
    return [
        {
            "git": "https://github.com/lorin/openstack-ansible-modules.git",
            "version": "master",
            "files": [{"src": "*_manage", "foo": "library/"}],
        }
    ]


@pytest.mark.parametrize(
    "gilt_config_file",
    ["missing_files_dst_key_data"],
    indirect=["gilt_config_file"],
)
def test_config_missing_files_dst_key(gilt_config_file):
    with pytest.raises(KeyError):
        config.config(gilt_config_file)


@pytest.mark.parametrize(
    "gilt_config_file", ["gilt_data"], indirect=["gilt_config_file"]
)
def test_get_config_generator(gilt_config_file):
    result = [i for i in config._get_config_generator(gilt_config_file)]

    assert isinstance(result, list)
    assert isinstance(result[0], dict)


def test_get_files_generator(temp_dir):
    files_list = [{"src": "foo", "dst": "bar/"}]
    result = [i for i in config._get_files_generator("/tmp/dir", files_list)]

    assert isinstance(result, list)
    assert isinstance(result[0], dict)


@pytest.mark.parametrize(
    "gilt_config_file", ["gilt_data"], indirect=["gilt_config_file"]
)
def test_get_config(gilt_config_file):
    result = config._get_config(gilt_config_file)

    assert isinstance(result, list)
    assert isinstance(result[0], dict)


@pytest.fixture()
def invalid_gilt_data():
    return "{"


@pytest.mark.parametrize(
    "gilt_config_file", ["invalid_gilt_data"], indirect=["gilt_config_file"]
)
def test_get_config_handles_parse_error(gilt_config_file):
    with pytest.raises(config.ParseError):
        config._get_config(gilt_config_file)


def test_get_dst_dir(temp_dir):
    os.chdir(temp_dir.strpath)
    result = config._get_dst_dir("roles/foo")

    assert os.path.join(temp_dir.strpath, "roles", "foo") == result


def test_get_clone_dir():
    parts = pytest.helpers.os_split(config._get_clone_dir())
    gilt_root = os.path.basename(config.BASE_WORKING_DIR)
    assert (gilt_root, "clone") == parts[-2:]


def test_makedirs(temp_dir):
    config._makedirs("foo/")

    d = os.path.join(temp_dir.strpath, "foo")
    assert os.path.isdir(d)

    curmask = os.umask(0)
    os.umask(curmask)

    expected = 0o777 & ~curmask
    assert expected == (os.stat(d).st_mode & 0o777)


def test_makedirs_nested_directory(temp_dir):
    config._makedirs("foo/bar/")

    d = os.path.join(temp_dir.strpath, "foo", "bar")
    assert os.path.isdir(d)


def test_makedirs_basedir(temp_dir):
    config._makedirs("foo/filename.py")

    d = os.path.join(temp_dir.strpath, "foo")
    assert os.path.isdir(d)


def test_makedirs_nested_basedir(temp_dir):
    config._makedirs("foo/bar/filename.py")

    d = os.path.join(temp_dir.strpath, "foo", "bar")
    assert os.path.isdir(d)


def test_makedirs_passes_if_exists(temp_dir):
    d = os.path.join(temp_dir.strpath, "foo")
    os.mkdir(d)

    config._makedirs("foo/")


def test_makedirs_raises(temp_dir):
    with pytest.raises(OSError):
        config._makedirs("")
