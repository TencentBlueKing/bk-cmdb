# godate
Golang package for intelligently parsing date strings like javascript's Date.parse() and getting the layout of date strings.
Fully compatible with the native time package.

This package is still under development.

# Usage
### Installation
```
go get github.com/joyt/godate
```

In your program:
```go
var timestamp time.Time
var err error
timestamp, err = date.Parse("Mar 14 2003")
```

## Example
```go
import (
  "github.com/joyt/godate"
  "time"
  "fmt"
)

func main() {
  var t time.Time
  
  t = date.MustParse("Aug 31 1999")
  fmt.Println(t)
  // Prints 1999-08-31 00:00:00 +0000 UTC
  
  t = date.MustParse("Tuesday, August 31, 1999")
  fmt.Println(t)
  // Prints 1999-08-31 00:00:00 +0000 UTC
  
  t = date.MustParse("Tue 31 Aug '99")
  fmt.Println(t)
  // Prints 1999-08-31 00:00:00 +0000 UTC
  
  t = date.MustParse("08/31/1999")
  fmt.Println(t)
  // Prints 1999-08-31 00:00:00 +0000 UTC
  
  t = date.MustParse("8/31/1999 20:05")
  fmt.Println(t)
  // Prints 1999-08-31 21:05:00 +0000 UTC
  
  t = date.MustParse("31/08/1999 8:05pm")
  fmt.Println(t)
  // Prints 1999-08-31 21:05:00 +0000 UTC
  
  t = date.MustParse("8/31/1999 8:05PM EST")
  fmt.Println(t)
  // Prints 1999-08-31 21:05:00 -0400 EDT
  
  t = date.MustParse("Aug-1999")
  fmt.Println(t)
  // Prints 1999-08-01 00:00:00 +0000 UTC
}
```

# Notes

The parser is extremely lenient, and will try to interpret whatever it is given as a date as much as possible.

In cases where the meaning of the date is ambiguous (such as 6/09, which could mean Jun 9th or Jun 2009), the parser generally defaults to the higher resolution date (Jun 9th). An exception is made for dates without separators such as "200609", where the parser will always try to assume the year is first (200609 -> Sep 2006, NOT Jun 20th 2009).
