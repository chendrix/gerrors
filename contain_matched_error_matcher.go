package gerrors

import (
	multierror "github.com/hashicorp/go-multierror"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

func ContainMatchedError(expected interface{}) types.GomegaMatcher {
	return &ContainMatchedErrorMatcher{
		Expected: expected,
	}
}

type ContainMatchedErrorMatcher struct {
	Expected interface{}
}

func (matcher *ContainMatchedErrorMatcher) Match(actual interface{}) (bool, error) {
	success, err := matchError(matcher.Expected, actual)

	// Matched or it's an incompatible use of the matcher
	if success || err != nil {
		return success, err
	}

	merr, ok := actual.(*multierror.Error)
	if !ok {
		return success, err
	}

	innerErrors := merr.WrappedErrors()
	for _, e := range innerErrors {
		success, err = matchError(e, matcher.Expected)
		// Early exit if it matches
		if err == nil && success {
			return success, err
		}
	}

	// Otherwise just return most recent result
	return success, err
}

func (matcher *ContainMatchedErrorMatcher) FailureMessage(actual interface{}) string {
	return format.Message(actual, "to contain matched error", matcher.Expected)
}

func (matcher *ContainMatchedErrorMatcher) NegatedFailureMessage(actual interface{}) string {
	return format.Message(actual, "not to contain matched error", matcher.Expected)
}
