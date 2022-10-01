package testutils

import (
	"fmt"
	"testing"
)

type Testcase[In any, Out any] struct {
	Name    string
	Input   In
	Output  Out
	Process func(In) Out
	Check   func(In, Out) error
}

func (tc *Testcase[In, Out]) RunFunc() func(*testing.T) {
	return func(t *testing.T) {
		tc.Output = tc.Process(tc.Input)

		err := tc.Check(tc.Input, tc.Output)
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}

func Assert[T any](
	name string,
	a T,
	b T,
	predicate func(T, T) bool,
) error {
	if !predicate(a, b) {
		return fmt.Errorf(
			"%s: expected %v got %v",
			name,
			a,
			b,
		)
	}

	return nil
}

func AssertNotNil[T any](
	name string,
	a *T,
) error {
	if a == nil {
		return fmt.Errorf(
			"%s: should have not been nil",
			name,
		)
	}

	return nil
}

func AssertNil[T any](
	name string,
	a *T,
) error {
	if a != nil {
		return fmt.Errorf(
			"%s: should have been nil",
			name,
		)
	}

	return nil
}

func StringsEqual(a string, b string) bool {
	return a == b
}
