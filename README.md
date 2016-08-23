[![Go Report Card](https://goreportcard.com/badge/github.com/chendrix/gerrors)](https://goreportcard.com/report/github.com/chendrix/gerrors)

Gomega matchers for wrapped errors
==================================

This package provides [Gomega](https://github.com/onsi/gomega) matchers to write assertions against errors wrapped using Dave Cheney's [github.com/pkg/errors](https://github.com/pkg/errors) package

MatchWrappedError()
-------------------
Verifies that an error matches the expected one, including errors that have been wrapped by a context.

This was needed because gomega's default `MatchError()` fails on wrapped errors

```go
import (
  "errors"
  perrors "github.com/pkg/errors"
  
  . "github.com/chendrix/gerrors"
)

err := errors.New("some error")
err2 := perrors.Wrap(err, "read failed")

Expect(err).To(MatchError("some error")) // Pass
Expect(err2).To(MatchError("some error")) // Fail!

Expect(err).To(MatchWrappedError("some error")) // Pass
Expect(err2).To(MatchWrappedError("some error")) // Pass
```
