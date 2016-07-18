package gerrors_test

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	. "github.com/chendrix/gerrors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	perrors "github.com/pkg/errors"
)

type CustomError struct {
}

func (c CustomError) Error() string {
	return "an error"
}

var _ = Describe("MatchWrappedErrorMatcher", func() {
	var (
		c CustomError
	)

	Context("when asserting an error vs a wrapped error", func() {
		It("succeeds when the underlying errors are the same", func() {
			Expect(sql.ErrNoRows).To(MatchWrappedError(perrors.Wrap(sql.ErrNoRows, "wrapped this error")))
			Expect(perrors.Wrap(sql.ErrNoRows, "wrapped this error")).To(MatchWrappedError(sql.ErrNoRows))

		})

		It("succeeds when the underlying errors are the same, even in multiple levels of wrapping", func() {
			Expect(sql.ErrNoRows).To(MatchWrappedError(
				perrors.Wrap(perrors.Wrap(sql.ErrNoRows, "wrapped this error"), "another wrapping"),
			))

			Expect(
				perrors.Wrap(perrors.Wrap(sql.ErrNoRows, "wrapped this error"), "another wrapping"),
			).To(MatchWrappedError(sql.ErrNoRows))
		})

		It("fails when passed nil", func() {
			// Expect(maybeError).To(MatchWrappedError(perrors.Wrap(nil, "wrapped this error"))
			// Bad Code, but this is what happens
			_, err := (&MatchWrappedErrorMatcher{
				Expected: perrors.Wrap(nil, "wrapped this error"),
			}).Match(nil)
			Expect(err).To(HaveOccurred())

			// Expect(maybeError).To(MatchWrappedError(nil))
			// Bad Code, but this is what happens
			_, err = (&MatchWrappedErrorMatcher{
				Expected: nil,
			}).Match(perrors.Wrap(nil, "wrapped this error"))
			Expect(err).To(HaveOccurred())
		})

		It("fails when the underlying error is different", func() {
			Expect(os.ErrNotExist).NotTo(MatchWrappedError(perrors.Wrap(sql.ErrNoRows, "wrapped this error")))
			Expect(perrors.Wrap(sql.ErrNoRows, "wrapped this error")).NotTo(MatchWrappedError(os.ErrNotExist))
		})
	})

	Context("when asserting a wrapped error vs a wrapped error", func() {
		It("succeeds when both sides are the same", func() {
			Expect(perrors.Wrap(c, "foo")).NotTo(MatchWrappedError(perrors.Wrap(c, "foo")))
			Expect(perrors.Wrap(perrors.Wrap(c, "foo"), "bar")).NotTo(MatchWrappedError(perrors.Wrap(perrors.Wrap(c, "foo"), "bar")))
		})

		It("fails when one of the levels of wrapping is different", func() {
			Expect(perrors.Wrap(c, "foo")).NotTo(MatchWrappedError(perrors.Wrap(c, "bar")))
			Expect(perrors.Wrap(perrors.Wrap(c, "foo"), "bar")).NotTo(MatchWrappedError(perrors.Wrap(perrors.Wrap(c, "foo"), "baz")))
			Expect(perrors.Wrap(perrors.Wrap(c, "foo"), "bar")).NotTo(MatchWrappedError(perrors.Wrap(perrors.Wrap(c, "baz"), "bar")))
		})

		It("fails when the underlying error is different", func() {
			Expect(perrors.Wrap(c, "wrapped this error")).NotTo(MatchWrappedError(perrors.Wrap(sql.ErrNoRows, "wrapped this error")))
		})
	})

	Context("when asserting an error vs an error (behaves like a MatchErrorMatcher)", func() {
		Context("When asserting against an error", func() {
			It("succeeds when matching with an error", func() {
				err := errors.New("an error")
				fmtErr := fmt.Errorf("an error")
				customErr := CustomError{}

				Expect(err).To(MatchWrappedError(errors.New("an error")))
				Expect(err).ToNot(MatchWrappedError(errors.New("another error")))

				Expect(fmtErr).To(MatchWrappedError(errors.New("an error")))
				Expect(customErr).To(MatchWrappedError(CustomError{}))
			})

			It("succeeds when matching with a string", func() {
				err := errors.New("an error")
				fmtErr := fmt.Errorf("an error")
				customErr := CustomError{}

				Expect(err).To(MatchWrappedError("an error"))
				Expect(err).ToNot(MatchWrappedError("another error"))

				Expect(fmtErr).To(MatchWrappedError("an error"))
				Expect(customErr).To(MatchWrappedError("an error"))
			})

			Context("when passed a matcher", func() {
				It("passes if the matcher passes against the error string", func() {
					err := errors.New("error 123 abc")

					Expect(err).To(MatchWrappedError(MatchRegexp(`\d{3}`)))
				})

				It("fails if the matcher fails against the error string", func() {
					err := errors.New("no digits")
					Expect(err).ToNot(MatchWrappedError(MatchRegexp(`\d`)))
				})
			})

			It("fails when passed anything else", func() {
				actualErr := errors.New("an error")
				_, err := (&MatchWrappedErrorMatcher{
					Expected: []byte("an error"),
				}).Match(actualErr)
				Expect(err).To(HaveOccurred())

				_, err = (&MatchWrappedErrorMatcher{
					Expected: 3,
				}).Match(actualErr)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when passed nil", func() {
			It("fails", func() {
				// Expect(maybeError).To(MatchWrappedError("an error"))
				_, err := (&MatchWrappedErrorMatcher{
					Expected: "an error",
				}).Match(nil)
				Expect(err).To(HaveOccurred())

				// Expect(someValidError).To(MatchWrappedError(nil))
				// Bad Code, but this is what happens
				_, err = (&MatchWrappedErrorMatcher{
					Expected: nil,
				}).Match(c)
				Expect(err).To(HaveOccurred())

				// Expect(maybeError).To(MatchWrappedError(nil))
				// Bad code, but this is what happens
				_, err = (&MatchWrappedErrorMatcher{
					Expected: nil,
				}).Match(nil)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when passed a non-error", func() {
			It("fails", func() {
				_, err := (&MatchWrappedErrorMatcher{
					Expected: "an error",
				}).Match("an error")
				Expect(err).To(HaveOccurred())

				_, err = (&MatchWrappedErrorMatcher{
					Expected: "an error",
				}).Match(3)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
