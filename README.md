[![Go Report Card](https://goreportcard.com/badge/github.com/chendrix/gerrors)](https://goreportcard.com/report/github.com/chendrix/gerrors)

Gomega matchers for other error packages
==================================

This package provides [Gomega](https://github.com/onsi/gomega) matchers to write assertions against errors:

- wrapped using Dave Cheney's [github.com/pkg/errors](https://github.com/pkg/errors) package
- included in Hashicorp's [github.com/hashicorp/go-multierror](https://github.com/hashicorp/go-multierror) package

This package is needed to address deficiencies in gomega's default `MatchError()`

MatchWrappedError()
-------------------
Verifies that an error matches the expected one, including errors that have been wrapped by a context.


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

ContainMatchedError()
---------------------
Verifies that an error is included in a `go-multierror` error.

```go
import (
  "errors"
   "github.com/hashicorp/go-multierror"
  
  . "github.com/chendrix/gerrors"
)

var result error 
err := errors.New("some error")

result = multierror.Append(result, err)

Expect(err).To(MatchError("some error")) // Pass
Expect(result).To(MatchError("some error")) // Fail!

Expect(err).To(ContainWrappedError("some error")) // Pass
Expect(result).To(ContainWrappedError("some error")) // Pass
```
