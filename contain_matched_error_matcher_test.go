package gerrors_test

import (
	. "github.com/chendrix/gerrors"

	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/go-multierror"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ContainMatchedError", func() {
	var (
		result error

		c CustomError
	)

	BeforeEach(func() {
		result = nil
	})

	Context("when asserting an error vs a multi error", func() {
		It("succeeds the underlying error is in the multi-error, but not the other way around", func() {
			Expect(multierror.Append(result, sql.ErrNoRows)).To(ContainMatchedError(sql.ErrNoRows))

			Expect(sql.ErrNoRows).NotTo(ContainMatchedError(multierror.Append(result, sql.ErrNoRows)))
		})

		It("succeeds when the underlying errors is in the multi-error, even with other errors", func() {
			result = multierror.Append(result, sql.ErrNoRows)
			result = multierror.Append(result, os.ErrNotExist)
			Expect(result).To(ContainMatchedError(sql.ErrNoRows))
			Expect(result).To(ContainMatchedError(os.ErrNotExist))
		})

		It("fails when passed nil", func() {
			// Expect(maybeError).To(ContainMatchedError(multiError.Append(nil, nil))
			// Bad Code, but this is what happens
			_, err := (&ContainMatchedErrorMatcher{
				Expected: multierror.Append(nil, nil),
			}).Match(nil)
			Expect(err).To(HaveOccurred())

			// Expect(maybeError).To(ContainMatchedError(nil))
			// Bad Code, but this is what happens
			_, err = (&ContainMatchedErrorMatcher{
				Expected: nil,
			}).Match(multierror.Append(nil, nil))
			Expect(err).To(HaveOccurred())
		})

		It("fails when the underlying error is not in the multi-error", func() {
			Expect(multierror.Append(nil)).NotTo(ContainMatchedError(sql.ErrNoRows))

			result = multierror.Append(result, sql.ErrNoRows)
			Expect(result).NotTo(ContainMatchedError(os.ErrNotExist))
		})
	})

	Context("when asserting an error vs an error (behaves like a MatchErrorMatcher)", func() {
		Context("When asserting against an error", func() {
			It("succeeds when matching with an error", func() {
				err := errors.New("an error")
				fmtErr := fmt.Errorf("an error")
				customErr := CustomError{}

				Expect(err).To(ContainMatchedError(errors.New("an error")))
				Expect(err).ToNot(ContainMatchedError(errors.New("another error")))

				Expect(fmtErr).To(ContainMatchedError(errors.New("an error")))
				Expect(customErr).To(ContainMatchedError(CustomError{}))
			})

			It("succeeds when matching with a string", func() {
				err := errors.New("an error")
				fmtErr := fmt.Errorf("an error")
				customErr := CustomError{}

				Expect(err).To(ContainMatchedError("an error"))
				Expect(err).ToNot(ContainMatchedError("another error"))

				Expect(fmtErr).To(ContainMatchedError("an error"))
				Expect(customErr).To(ContainMatchedError("an error"))
			})

			Context("when passed a matcher", func() {
				It("passes if the matcher passes against the error string", func() {
					err := errors.New("error 123 abc")

					Expect(err).To(ContainMatchedError(MatchRegexp(`\d{3}`)))
				})

				It("fails if the matcher fails against the error string", func() {
					err := errors.New("no digits")
					Expect(err).ToNot(ContainMatchedError(MatchRegexp(`\d`)))
				})
			})

			It("fails when passed anything else", func() {
				actualErr := errors.New("an error")
				_, err := (&ContainMatchedErrorMatcher{
					Expected: []byte("an error"),
				}).Match(actualErr)
				Expect(err).To(HaveOccurred())

				_, err = (&ContainMatchedErrorMatcher{
					Expected: 3,
				}).Match(actualErr)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when passed nil", func() {
			It("fails", func() {
				// Expect(maybeError).To(ContainMatchedError("an error"))
				_, err := (&ContainMatchedErrorMatcher{
					Expected: "an error",
				}).Match(nil)
				Expect(err).To(HaveOccurred())

				// Expect(someValidError).To(ContainMatchedError(nil))
				// Bad Code, but this is what happens
				_, err = (&ContainMatchedErrorMatcher{
					Expected: nil,
				}).Match(c)
				Expect(err).To(HaveOccurred())

				// Expect(maybeError).To(ContainMatchedError(nil))
				// Bad code, but this is what happens
				_, err = (&ContainMatchedErrorMatcher{
					Expected: nil,
				}).Match(nil)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when passed a non-error", func() {
			It("fails", func() {
				_, err := (&ContainMatchedErrorMatcher{
					Expected: "an error",
				}).Match("an error")
				Expect(err).To(HaveOccurred())

				_, err = (&ContainMatchedErrorMatcher{
					Expected: "an error",
				}).Match(3)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
