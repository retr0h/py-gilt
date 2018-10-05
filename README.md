[![Build Status](http://img.shields.io/travis/retr0h/go-gilt.svg?style=flat-square)](https://travis-ci.org/retr0h/go-gilt)
[![Coveralls github](https://img.shields.io/coveralls/github/retr0h/go-gilt.svg?style=flat-square)](https://coveralls.io/github/retr0h/go-gilt)
[![Go Report Card](https://goreportcard.com/badge/github.com/retr0h/go-gilt?style=flat-square)](https://goreportcard.com/report/github.com/retr0h/go-gilt)

# go-gilt

Gilt is a tool which aims to make repo management, manageable.  Gilt
clones repositories at a particular version, then overlays the repository to
the provided destination.  An alternate approach to "vendoring".

What makes Gilt interesting, is the ability to overlay particular files and/or
directories from the specified repository to given destinations.  This is quite
helpful for those using Ansible, since libraries, plugins, and playbooks are
often shared, but Ansible's [Galaxy][1] has no mechanism to handle this.

[1]: https://docs.ansible.com/ansible/latest/reference_appendices/galaxy.html

## Port

This project is a port of [Gilt](http://gilt.readthedocs.io/en/latest/), it is
not 100% compatible with the python version, and aims to correct some poor decisions
made in the python version of Gilt.

This version of Gilt does not provide built in locking, unlike our python friend. If one
wishes to use locking, [flock(1)](https://linux.die.net/man/1/flock) should be used.

## Installation

    $  go get github.com/retr0h/go-gilt

## Usage

### Overlay Repository

Create the giltfile (`gilt.yml`).

Clone the specified `url`@`version` to the configurable path `--giltdir`.
Extract the repo the `dstDir` when `dstDir` is provided.  Otherwise, copy files
and/or directories to the desired destinations.

```yaml
---
- git: https://github.com/retr0h/ansible-etcd.git
  version: 77a95b7
  dstDir: roles/retr0h.ansible-etcd

- git: https://github.com/lorin/openstack-ansible-modules.git
  version: 2677cc3
  sources:
    - src: "*_manage"
      dstDir: library
    - src: nova_quota
      dstDir: library
    - src: neutron_router
      dstFile: library/neutron_router.py
    - src: tests
      dstDir: tests
```

Overlay a remote repository into the destination provided.

```bash
$ gilt overlay
```

Use an alternate config file (default `gilt.yml`).

```bash
$ gilt overlay --filename /path/to/gilt.yml
```

Optionally, override gilt's cache location (defaults to `~/.gilt/clone`):

```bash
$ gilt --giltdir ~/alternate/directory overlay
```

### Debug

Display the git commands being executed.

```bash
$ gilt --debug overlay
```

[![asciicast](https://asciinema.org/a/195036.png)](https://asciinema.org/a/195036?speed=2&autoplay=1&loop=1)

## Dependencies

```bash
$ go get github.com/golang/dep/cmd/dep
```

## Building

```bash
$ make build
$ tree .build/
```

## Testing

```bash
$ make test
```

## License

MIT
