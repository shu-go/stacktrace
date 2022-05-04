package stacktrace_test

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/shu-go/stacktrace"
)

func Example() {
	err := funcA()

	fmt.Printf("%v\n", stacktrace.New(err))
	fmt.Printf("%+v\n", stacktrace.New(err))
	fmt.Printf("%+v\n", err)

	// Output:
}

func funcA() error {
	return errors.Wrap(funcB(), "error A")
}

func funcB() error {
	return errors.Wrap(funcC(), "error B")
}

func funcC() error {
	return errors.New("error C")
}
