# Annotated errors

This package implements error values annotated by the file name, line number and
the function name where the error happened. It supports a "stack trace" style
errors by annotating existing errors as they pass through the calling functions.

This package is inspired by the
[`luci/common/errors`](https://pkg.go.dev/go.chromium.org/luci)
package. However, LUCI is a huge mono-repo, which may be an overkill for smaller
projects, or projects unrelated to Cloud apps.

This repo is dedicated to a very light-weght implementation of a simple API that
gets most of the job done.

## Installation

```
go get github.com/stockparfait/errors
```

## Example usage

```go
// example.go

package main

import (
	"fmt"

	"github.com/stockparfait/errors"
)

func Use(x int) error {
	if x < 0 {
		return errors.Reason("x = %d is negative", x)
	}
	return nil
}

func Top(val int) error {
	if err := Use(val); err != nil {
		return errors.Annotate(err, "cannot use %d", val)
	}
	return nil
}

func main() {
	fmt.Printf("%s\n", Top(-42).Error())
}
```

This should print something like this:
```
ERROR: /path/to/example.go:20: main.Top() cannot use -42:
ERROR: /path/to/example.go:13: main.Use() x = -42 is negative
```

## Development

Clone and initialize the repository, run tests:

```sh
git clone git@github.com:stockparfait/errors.git
cd errors
make init
make test
```
