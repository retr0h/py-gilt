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

import glob
import os

import pytest
import sh

from gilt import git

slow = pytest.mark.skipif(
    not pytest.config.getoption("--runslow"),
    reason="need --runslow option to run")


@slow
def test_clone(temp_dir):
    name = 'retr0h.ansible-etcd'
    repo = 'https://github.com/retr0h/ansible-etcd.git'
    destination = os.path.join(temp_dir.strpath, name)

    git.clone(name, repo, destination)

    assert os.path.exists(destination)


@slow
def test_extract(temp_dir):
    name = 'retr0h.ansible-etcd'
    repo = 'https://github.com/retr0h/ansible-etcd.git'
    branch = 'gilt'
    clone_dir = os.path.join(temp_dir.strpath, 'clone_dir')
    extract_dir = os.path.join(temp_dir.strpath, 'extract_dir')

    os.mkdir(clone_dir)
    os.mkdir(extract_dir)

    clone_destination = os.path.join(clone_dir, name)
    extract_destination = os.path.join(extract_dir, name, '')

    git.clone(name, repo, clone_destination)
    git.extract(clone_destination, extract_destination, branch)

    assert os.path.exists(extract_destination)
    git_dir = os.path.join(extract_destination, os.extsep.join(('', 'git')))
    assert not os.path.exists(git_dir)
    giltfile = os.path.join(extract_destination, 'giltfile')
    assert os.path.exists(giltfile)


@slow
def test_overlay(mocker, temp_dir):
    name = 'lorin.openstack-ansible-modules'
    repo = 'https://github.com/lorin/openstack-ansible-modules.git'
    branch = 'master'
    clone_dir = os.path.join(temp_dir.strpath, name)
    dst_dir = os.path.join(temp_dir.strpath, 'dst', '')

    os.mkdir(clone_dir)
    os.mkdir(dst_dir)

    # yapf: disable
    files = [
        mocker.Mock(src=os.path.join(clone_dir, '*_manage'), dst=dst_dir),
        mocker.Mock(src=os.path.join(clone_dir, 'nova_quota'), dst=dst_dir),
        mocker.Mock(src=os.path.join(clone_dir, 'neutron_router'),
                    dst=os.path.join(dst_dir, 'neutron_router.py')),
        mocker.Mock(
            src=os.path.join(clone_dir, 'tests'),
            dst=os.path.join(dst_dir, 'tests'))
    ]
    # yapf: enable
    git.clone(name, repo, clone_dir)
    git.overlay(clone_dir, files, branch)

    assert 5 == len(glob.glob('{}/*_manage'.format(dst_dir)))
    assert 1 == len(glob.glob('{}/nova_quota'.format(dst_dir)))
    assert 1 == len(glob.glob('{}/neutron_router.py'.format(dst_dir)))
    assert 2 == len(glob.glob('{}/*'.format(os.path.join(dst_dir, 'tests'))))


@slow
def test_overlay_existing_directory(mocker, temp_dir):
    name = 'lorin.openstack-ansible-modules'
    repo = 'https://github.com/lorin/openstack-ansible-modules.git'
    branch = 'master'
    clone_dir = os.path.join(temp_dir.strpath, name)
    dst_dir = os.path.join(temp_dir.strpath, 'dst', '')

    os.mkdir(clone_dir)
    os.mkdir(dst_dir)
    os.mkdir(os.path.join(dst_dir, 'tests'))

    files = [
        mocker.Mock(
            src=os.path.join(clone_dir, 'tests'),
            dst=os.path.join(dst_dir, 'tests'))
    ]
    git.clone(name, repo, clone_dir)
    git.overlay(clone_dir, files, branch)

    assert 2 == len(glob.glob('{}/*'.format(os.path.join(dst_dir, 'tests'))))


@pytest.fixture()
def patched_run_command(mocker):
    return mocker.patch('gilt.util.run_command')


def test_get_version_has_branch(mocker, patched_run_command):
    mocker.patch('gilt.git._has_branch').return_value = True
    mocker.patch('gilt.git._has_tag').return_value = False
    mocker.patch('gilt.git._has_commit').return_value = False
    git._get_version('branch')
    # yapf: disable
    expected = [
        mocker.call(sh.git.bake('checkout', 'branch'), debug=False),
        mocker.call(sh.git.bake('clean', '-d', '-x', '-f'), debug=False),
        mocker.call(sh.git.bake('pull', rebase=True, ff_only=True),
                    debug=False)
    ]
    # yapf: enable

    assert expected == patched_run_command.mock_calls


def test_get_version_has_tag(mocker, patched_run_command):
    mocker.patch('gilt.git._has_branch').return_value = False
    mocker.patch('gilt.git._has_tag').return_value = True
    mocker.patch('gilt.git._has_commit').return_value = False
    git._get_version('tag_name')
    # yapf: disable
    expected = [
        mocker.call(sh.git.bake('checkout', 'tag_name'), debug=False),
        mocker.call(sh.git.bake('clean', '-d', '-x', '-f'), debug=False)
    ]
    # yapf: enable

    assert expected == patched_run_command.mock_calls


def test_get_version_has_commit(mocker, patched_run_command):
    mocker.patch('gilt.git._has_branch').return_value = False
    mocker.patch('gilt.git._has_tag').return_value = False
    mocker.patch('gilt.git._has_commit').return_value = True
    git._get_version('commit_sha')
    # yapf: disable
    expected = [
        mocker.call(sh.git.bake('checkout', 'commit_sha'), debug=False),
        mocker.call(sh.git.bake('clean', '-d', '-x', '-f'), debug=False)
    ]
    # yapf: enable

    assert expected == patched_run_command.mock_calls


def test_get_version_needs_fetch(mocker, patched_run_command):
    mocker.patch('gilt.git._has_branch').return_value = False
    mocker.patch('gilt.git._has_tag').return_value = False
    mocker.patch('gilt.git._has_commit').return_value = False
    git._get_version('remote_tag')
    # yapf: disable
    expected = [
        mocker.call(sh.git.bake('fetch'), debug=False),
        mocker.call(sh.git.bake('checkout', 'remote_tag'), debug=False),
        mocker.call(sh.git.bake('clean', '-d', '-x', '-f'), debug=False)
    ]
    # yapf: enable

    assert expected == patched_run_command.mock_calls


def test_get_version_needs_pull(mocker, patched_run_command):
    mocker.patch('gilt.git._has_branch').side_effect = [False, True]
    mocker.patch('gilt.git._has_tag').return_value = False
    mocker.patch('gilt.git._has_commit').return_value = False
    git._get_version('remote_branch')
    # yapf: disable
    expected = [
        mocker.call(sh.git.bake('fetch'), debug=False),
        mocker.call(sh.git.bake('checkout', 'remote_branch'), debug=False),
        mocker.call(sh.git.bake('clean', '-d', '-x', '-f'), debug=False),
        mocker.call(sh.git.bake('pull', rebase=True, ff_only=True),
                    debug=False)
    ]
    # yapf: enable

    assert expected == patched_run_command.mock_calls


@slow
def test_has_version(temp_dir):
    name = 'retr0h.ansible-etcd'
    repo = 'https://github.com/retr0h/ansible-etcd.git'
    destination = os.path.join(temp_dir.strpath, name)
    git.clone(name, repo, destination)
    os.chdir(destination)
    # _has_branch tests
    assert git._has_branch('master')
    assert not git._has_branch('1.1')
    assert not git._has_branch('888ef7b')

    # _has_tag tests
    assert not git._has_tag('master')
    assert git._has_tag('1.1')
    assert not git._has_tag('888ef7b')

    # _has_commit tests
    assert not git._has_commit('master')
    assert not git._has_commit('1.1')
    assert git._has_commit('888ef7b')
