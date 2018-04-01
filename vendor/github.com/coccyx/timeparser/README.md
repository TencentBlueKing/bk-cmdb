Timeparser
==========

[![Build Status](https://travis-ci.org/coccyx/timeparser.svg)](https://travis-ci.org/coccyx/timeparser) [![Go Report Card](https://goreportcard.com/badge/github.com/coccyx/timeparser)](https://goreportcard.com/report/github.com/coccyx/timeparser) [![GoDoc](https://godoc.org/github.com/coccyx/timeparser?status.svg)](https://godoc.org/github.com/coccyx/timeparser)

Go (golang) package for parsing time in Splunk's relative time syntax.  [See docs for format documentation](https://docs.splunk.com/Documentation/Splunk/6.4.3/SearchReference/SearchTimeModifiers), but here's some examples:

| Time             | Description                                                                        |
|------------------|------------------------------------------------------------------------------------|
| -1m              | One minute ago                                                                     |
| +30m             | 30 minutes from now                                                                |
| -4h@h            | Four hours ago, snapped to the hour                                                |
| -1week@week+1day | One week ago, snapped to Monday (1 day after Sunday, which is the default snap to) |


For details, see the [API documentation](http://godoc.org/github.com/coccyx/timeparser).