[![Build Status](https://img.shields.io/travis/walle/lll.svg?style=flat)](https://travis-ci.org/walle/lll)
[![Coverage](https://img.shields.io/codecov/c/github/walle/lll.svg?style=flat)](https://codecov.io/github/walle/lll)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/walle/lll)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/walle/lll/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/walle/lll)](http:/goreportcard.com/report/walle/lll)

# lll

Line length linter, used to enforce line length in files.
Support for only checking go files.

## Installation

```shell
$ go get github.com/walle/lll/...
```

## Usage

```shell
usage: lll [--maxlength MAXLENGTH] [--tabwidth TABWIDTH] [--goonly] [--skiplist SKIPLIST] [--vendor] [--files] [--exclude EXCLUDE] [INPUT [INPUT ...]]

positional arguments:
  input

options:
  --maxlength MAXLENGTH, -l MAXLENGTH
                         max line length to check for [default: 80]
  --tabwidth TABWIDTH, -w TABWIDTH
                         tab width in spaces [default: 1]
  --goonly, -g           only check .go files
  --skiplist SKIPLIST, -s SKIPLIST
                         list of dirs to skip [default: [.git vendor]]
  --vendor               check files in vendor directory
  --files                read file names from stdin one at each line
  --exclude EXCLUDE, -e EXCLUDE
                         exclude lines that matches this regex
  --help, -h             display this help and exit
```

Example usage to check only go files for lines more than 100 characters.
Excluding lines that contain the words TODO or FIXME.
`lll -l 100 -g -e "TODO|FIXME" path/to/myproject`.

You can also define the flags using environment variables, eg. 
`MAXLENGTH=100 GOONLY=true lll path/to/my/project`.

## Testing

Use the `go test` tool.

```shell
$ go test -cover
```

## Contributing

All contributions are welcome! See [CONTRIBUTING](CONTRIBUTING.md) for more
info.

## License

The code is under the MIT license. See [LICENSE](LICENSE) for more
information.
