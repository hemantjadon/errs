package errs_test

import (
	"errors"
	"fmt"

	"github.com/hemantjadon/errs"
)

func ExampleNew() {
	fn := func() error {
		return errs.New("error occurred")
	}
	fmt.Println(fn())
	// Output: error occurred
}

func ExampleNew_withFields() {
	fn := func() error {
		return errs.New("error occurred", errs.F("temperature", 10), errs.F("state", "heating"))
	}
	fmt.Println(fn())
	// Output: error occurred (temperature=10 state=heating)
}

func ExampleBox() {
	fn := func() error {
		err := doSomething()
		if err != nil {
			return errs.Box(err, "error occurred")
		}
		return nil
	}
	fmt.Println(fn())
	fmt.Println(errors.Is(fn(), ErrSomethingFailed))
	// Output:
	// error occurred: something failed
	// false
}

func ExampleBox_withFields() {
	fn := func() error {
		err := doSomething()
		if err != nil {
			return errs.Box(err, "error occurred", errs.F("temperature", 10), errs.F("state", "heating"))
		}
		return nil
	}
	fmt.Println(fn())
	fmt.Println(errors.Is(fn(), ErrSomethingFailed))
	// Output:
	// error occurred (temperature=10 state=heating): something failed
	// false
}

func ExampleWrap() {
	fn := func() error {
		err := doSomething()
		if err != nil {
			return errs.Wrap(err, "error occurred")
		}
		return nil
	}
	fmt.Println(fn())
	fmt.Println(errors.Is(fn(), ErrSomethingFailed))
	// Output:
	// error occurred: something failed
	// true
}

func ExampleWrap_withFields() {
	fn := func() error {
		err := doSomething()
		if err != nil {
			return errs.Wrap(err, "error occurred", errs.F("temperature", 10), errs.F("state", "heating"))
		}
		return nil
	}
	fmt.Println(fn())
	fmt.Println(errors.Is(fn(), ErrSomethingFailed))
	// Output:
	// error occurred (temperature=10 state=heating): something failed
	// true
}

func ExampleWrap_chain() {
	fn := func() error {
		err := doSomething()
		if err != nil {
			return errs.Wrap(err, "error occurred", errs.F("temperature", 10), errs.F("state", "heating"))
		}
		return nil
	}
	fmt.Println(fn().(errs.ChainError).Chain())
	// Output:
	// [error occurred (temperature=10 state=heating) something failed]
}

const ErrSomethingFailed testErr = "something failed"

func doSomething(_ ...string) error {
	return ErrSomethingFailed
}

type testErr string

func (e testErr) Error() string {
	return string(e)
}
