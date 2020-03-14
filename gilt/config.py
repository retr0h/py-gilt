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
import errno
import os
import urllib.parse

import yaml

from gilt import interpolation


class ParseError(Exception):
    """Error raised when a config can't be loaded properly. """

    pass


BASE_WORKING_DIR = os.environ.get("GILT_CACHE_DIRECTORY", "~/.gilt")


def config(filename):
    """Construct `Config` object and return a list.

    :parse filename: A string containing the path to YAML file.
    :return: list
    """
    Config = collections.namedtuple(
        "Config",
        [
            "git",
            "lock_file",
            "version",
            "name",
            "src",
            "dst",
            "files",
            "post_commands",
        ],
    )

    return [Config(**d) for d in _get_config_generator(filename)]


def _parse_repo_uri(uri, scm="git"):
    """Construct and return a `ParsedRepo` object.

    :param uri: A SCM repository URI.
    :return: ParsedRepo
    """
    ParsedRepo = collections.namedtuple('ParsedRepo', [
        'hostname',
        'owner',
        'name'
    ])

    scm_ext = "." + scm
    o = urllib.parse.urlparse(uri)
    if not o.hostname:  # scp-style "URI", so fake it
        uri = "ssh://" + uri.replace(":/", "/").replace(":", "/")
        o = urllib.parse.urlparse(uri)
    path = o.path
    if path.endswith(scm_ext):
        path = path[0:-len(scm_ext)]
    try:
        name, owner = path.rsplit("/", 2)[-1:-3:-1]
        owner = owner.lstrip("~")
    except ValueError:
        # if repo is the root of the path, then there is no owner
        owner = ''

    return ParsedRepo(o.hostname, owner, name)


def _get_files_config(src_dir, files_list):
    """Construct `FileConfig` object and return a list.

    :param src_dir: A string containing the source directory.
    :param files_list: A list of dicts containing the src/dst mapping of files
     to overlay.
    :return: list
    """
    FilesConfig = collections.namedtuple(
        "FilesConfig", ["src", "dst", "post_commands"]
    )

    return [
        FilesConfig(**d) for d in _get_files_generator(src_dir, files_list)
    ]


def _get_config_generator(filename):
    """A generator which populates and return a dict.

    :parse filename: A string containing the path to YAML file.
    :return: dict
    """
    for d in _get_config(filename):
        repo = d['git']
        parsedrepo = _parse_repo_uri(repo)
        if parsedrepo.owner:
            name = '{}.{}'.format(parsedrepo.owner, parsedrepo.name)
        else:
            name = parsedrepo.name
        src_dir = os.path.join(_get_clone_dir(), parsedrepo.hostname, name)
        files = d.get('files')
        post_commands = d.get('post_commands', [])
        dst_dir = None
        if not files:
            dst_dir = _get_dst_dir(d["dst"])
        yield {
            'git': repo,
            'lock_file': os.path.join(
                _get_lock_dir(), parsedrepo.hostname, name
            ),
            'version': d['version'],
            'name': name,
            'src': src_dir,
            'dst': dst_dir,
            'files': _get_files_config(src_dir, files),
            'post_commands': post_commands,
        }


def _get_files_generator(src_dir, files_list):
    """A generator which populates and return a dict.

    :param src_dir: A string containing the source directory.
    :param files_list: A list of dicts containing the src/dst mapping of files
     to overlay.
    :return: dict
    """
    if files_list:
        for d in files_list:
            yield {
                "src": os.path.join(src_dir, d["src"]),
                "dst": _get_dst_dir(d["dst"]),
                "post_commands": d.get("post_commands", []),
            }


def _get_config(filename):
    """Parse the provided YAML file and return a dict.

    :parse filename: A string containing the path to YAML file.
    :return: dict
    """
    i = interpolation.Interpolator(
        interpolation.TemplateWithDefaults, os.environ
    )

    with open(filename, "r") as stream:
        try:
            interpolated_config = i.interpolate(stream.read())
            return yaml.safe_load(interpolated_config)
        except yaml.parser.ParserError as e:
            msg = "Error parsing gilt config: {0}".format(e)
            raise ParseError(msg)


def _get_dst_dir(dst_dir):
    """Prefix the provided string with working directory and return a str.

    :param dst_dir: A string to be prefixed with the working dir.
    :return: str
    """
    wd = os.getcwd()
    _makedirs(dst_dir)

    return os.path.join(wd, dst_dir)


def _get_base_dir():
    """Return gilt's base working directory. """
    return os.path.expanduser(BASE_WORKING_DIR)


def _get_lock_dir():
    """Construct gilt's lock directory and return a str.

    :return: str
    """
    return os.path.join(_get_base_dir(), "lock",)


def _get_clone_dir():
    """Construct gilt's clone directory and return a str.

    :return: str
    """
    return os.path.join(_get_base_dir(), "clone",)


def _makedirs(path):
    """Create a base directory of the provided path and return None.

    :param path: A string containing a path to be deconstructed and basedir
     created.
    :return: None
    """
    dirname, _ = os.path.split(path)
    try:
        os.makedirs(dirname)
    except OSError as exc:
        if exc.errno == errno.EEXIST:
            pass
        else:
            raise
