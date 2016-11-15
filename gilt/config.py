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

import collections
import os

import giturlparse
import yaml


class ParseError(Exception):
    """ Error raised when a config can't be loaded properly. """
    pass


def config(filename):
    """
    Construct `Config` object and return a list.

    :parse filename: A string containing the path to YAML file.
    :return: list
    """
    Config = collections.namedtuple(
        'Config', ['git', 'version', 'name', 'src', 'dst', 'files'])

    return [Config(**d) for d in _get_config_generator(filename)]


def _get_files_config(src_dir, files_list):
    """
    Construct `FileConfig` object and return a list.

    :param src_dir: A string containing the source directory.
    :param files_list: A list of dicts containing the src/dst mapping of files
     to overlay.
    :return: list
    """
    FilesConfig = collections.namedtuple('FilesConfig', ['src', 'dst'])

    return [
        FilesConfig(**d) for d in _get_files_generator(src_dir, files_list)
    ]


def _get_config_generator(filename):
    """
    A generator which populates and return a dict.

    :parse filename: A string containing the path to YAML file.
    :return: dict
    """
    for d in _get_config(filename):
        parsedrepo = giturlparse.parse(d.get('git'), False)
        name = '{}.{}'.format(parsedrepo.owner, parsedrepo.repo)
        src_dir = os.path.join(_get_clone_dir(), name)
        yield {
            'git': d.get('git'),
            'version': d.get('version'),
            'name': name,
            'src': src_dir,
            'dst': _get_dst_dir(d),
            'files': _get_files_config(src_dir, d.get('files'))
        }


def _get_files_generator(src_dir, files_list):
    """
    A generator which populates and return a dict.

    :param src_dir: A string containing the source directory.
    :param files_list: A list of dicts containing the src/dst mapping of files
     to overlay.
    :return: dict
    """
    if files_list:
        for d in files_list:
            yield {
                'src': os.path.join(src_dir, d.get('src')),
                'dst': _get_dst_dir(d)
            }


def _get_config(filename):
    """
    Parse the provided YAML file and return a dict.

    :parse filename: A string containing the path to YAML file.
    :return: dict
    """
    with open(filename, 'r') as stream:
        try:
            return yaml.safe_load(stream)
        except yaml.parser.ParserError as e:
            msg = 'Error parsing gilt config: {0}'.format(e)
            raise ParseError(msg)


def _get_dst_dir(d):
    """
    Prefix the provided dict's `dst` key with working directory and return a
    dict.

    :param d: A config dict which contains a `dst` key.
    :return: str
    """
    wd = os.getcwd()
    if d.get('dst'):
        return os.path.join(wd, d.get('dst'))


def _get_clone_dir():
    """
    Construct gilt's clone directory and return a str.

    :return: str
    """
    path = '~/.gilt/clone'

    return os.path.expanduser(path)
