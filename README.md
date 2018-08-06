[![Build Status](http://img.shields.io/travis/retr0h/go-gilt.svg?style=flat-square)](https://travis-ci.org/retr0h/go-gilt)
[![Coveralls github](https://img.shields.io/coveralls/github/retr0h/go-gilt.svg?style=flat-square)](https://coveralls.io/github/retr0h/go-gilt)
[![Go Report Card](https://goreportcard.com/badge/github.com/retr0h/go-gilt?style=flat-square)](https://goreportcard.com/report/github.com/retr0h/go-gilt)

# go-gilt

Gilt is a tool which aims to make Ansible repo management, managable.  Gilt
clones repositories at a particular version, then overlays the repository to
the provided destination.

What makes Gilt interesting, is the ability to overlay particular files and/or
directories from the specified repository to given destinations.  This is quite
hepful when working with Ansible, since libraries, plugins, and playbooks are
often shared, but [Galaxy][1] has no mechanism to cleanly handle this.

[1]: https://docs.ansible.com/ansible/latest/reference_appendices/galaxy.html

## Port

This project is a port of [Gilt](http://gilt.readthedocs.io/en/latest/), it is
not 100% compatible with the python version, and not yet complete.

## Installation

    $  go get github.com/retr0h/gofile

## Usage

### Overlay Repository

Create the giltfile (`gilt.yml`).

Clone the specified `url`@`version` to the configurable path `--giltdir`, and
extract the repository to the provided `dst`.

```yaml
---
- url: https://github.com/retr0h/ansible-etcd.git
  version: 77a95b7
  dst: roles/retr0h.ansible-etcd
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
