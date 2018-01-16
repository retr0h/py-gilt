[![Build Status](http://img.shields.io/travis/retr0h/go-gilt.svg?style=flat-square)](https://travis-ci.org/retr0h/go-gilt)
[![Coveralls github](https://img.shields.io/coveralls/github/retr0h/go-gilt.svg?style=flat-square)](https://coveralls.io/github/retr0h/go-gilt)
[![Go Report Card](https://goreportcard.com/badge/github.com/retr0h/go-gilt?style=flat-square)](https://goreportcard.com/report/github.com/retr0h/go-gilt)

# go-gilt

## Installation

    $  go get github.com/retr0h/gofile

## Usage

Overlay a remote repository into the destination provided.

```yaml
---
- url: https://github.com/retr0h/ansible-etcd.git
  version: 77a95b7
  dst: roles/retr0h.ansible-etcd
```

```bash
$ gilt overlay
```

Optionally, override gilt's cache location (defaults to `~/.gilt/cache`):

```bash
$ gilt --giltdir ~/alternate/directory overlay
```

Display the git commands being executed.

```bash
$ gilt --debug overlay
```

Use an alternate config file (default `gilt.yml`).

```bash
$ gilt overlay --filename /path/to/gilt.yml
```

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
