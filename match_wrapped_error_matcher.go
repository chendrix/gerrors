package gerrors

import (
	"github.com/onsi/gomega/types"

	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/matchers"
)

func MatchWrappedError(expected interface{}) types.GomegaMatcher {
	return &MatchWrappedErrorMatcher{
		Expected: expected,
	}

}

type MatchWrappedErrorMatcher struct {
	Expected interface{}
}

type Causer interface {
	Cause() error
}

func (matcher *MatchWrappedErrorMatcher) Match(actual interface{}) (bool, error) {
	success, err := matchError(matcher.Expected, actual)

	// Matched or it's an incompatible use of the matcher
	if success || err != nil {
		return success, err
	}

	// Did not match, but it's a correct usage of the matcher
	// Time to check to see if it's actually a valid wrapped error
	actualErr, aok := actual.(error)
	expectedErr, eok := matcher.Expected.(error)

	// One of them is not actually an error, default to what MatchError says
	if !aok || !eok {
		return success, err
	}

	_, expectedCauserOk := matcher.Expected.(Causer)
	_, actualCauserOK := actual.(Causer)

	// XOR, one of them is a wrapped error and the other isn't
	if expectedCauserOk != actualCauserOK {
		underlyingExpected := unwindError(expectedErr)
		underlyingActual := unwindError(actualErr)
		return matchError(underlyingExpected, underlyingActual)
	}

	// Both wrapped errors
	return reflect.DeepEqual(matcher.Expected, actual), nil
}

func (matcher *MatchWrappedErrorMatcher) FailureMessage(actual interface{}) string {
	return format.Message(actual, "to match wrapped error", matcher.Expected)
}

func (matcher *MatchWrappedErrorMatcher) NegatedFailureMessage(actual interface{}) string {
	return format.Message(actual, "not to match wrapped error", matcher.Expected)
}

func matchError(expected, actual interface{}) (bool, error) {
	return (&matchers.MatchErrorMatcher{
		Expected: expected,
	}).Match(actual)
}

func unwindError(e error) error {
	causer, ok := e.(Causer)
	if !ok {
		return e
	}

	return unwindError(causer.Cause())
}
