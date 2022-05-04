// Package stacktrace extracts merged/simplified stacktrace from wrapped errors.
//
// Usage
//
//     func funcA() error {
//     	return errors.Wrap(funcB(), "error A")
//     }
//
//     func funcB() error {
//     	return errors.Wrap(funcC(), "error B")
//     }
//
//     func funcC() error {
//     	return errors.New("error C")
//     }
//
//
// Case 1: Simple output
//
//     fmt.Printf("%+v\n", a())
//
// vvv
//
//     error C
//             path/to/mypkg/stacktrace_test.go:30
//     error B
//             path/to/mypkg/stacktrace_test.go:26
//     error A
//             path/to/mypkg/stacktrace_test.go:22
//             path/to/mypkg/stacktrace_test.go:12
//             go/current/src/testing/run_example.go:63
//             go/current/src/testing/example.go:44
//             go/current/src/testing/testing.go:1721
//             _testmain.go:49
//             go/current/src/runtime/proc.go:250
//             go/current/src/runtime/asm_amd64.s:1571
//
//
// Case 2: With function names
//
//     fmt.Printf("%+v\n", stacktrace.New(err))
//
// vvv
//
//     error C
//     mypkg.funcC
//             path/to/mypkg/mysource.go:30
//     error B
//     mypkg.funcB
//             path/to/mypkg/mysource.go:26
//     error A
//     mypkg.funcA
//             path/to/mypkg/mysource.go:22
//     mypkg.Example
//             path/to/mypkg/mysource.go:12
//     testing.runExample
//             go/current/src/testing/run_example.go:63
//     testing.runExamples
//             go/current/src/testing/example.go:44
//     testing.(*M).Run
//             go/current/src/testing/testing.go:1721
//     main.main
//             _testmain.go:49
//     runtime.main
//             go/current/src/runtime/proc.go:250
//     runtime.goexit
//             go/current/src/runtime/asm_amd64.s:1571
//
//
// Case 3: (without this package) fmt.Printf("%+v\n", a())
//
// Redundant output.
//
//     error C
//     mypkg.funcC
//             path/to/mypkg/mysource.go:30
//     mypkg.funcB
//             path/to/mypkg/mysource.go:26
//     mypkg.funcA
//             path/to/mypkg/mysource.go:22
//     mypkg.Example
//             path/to/mypkg/mysource.go:12
//     testing.runExample
//             go/current/src/testing/run_example.go:63
//     testing.runExamples
//             go/current/src/testing/example.go:44
//     testing.(*M).Run
//             go/current/src/testing/testing.go:1721
//     main.main
//             _testmain.go:49
//     runtime.main
//             go/current/src/runtime/proc.go:250
//     runtime.goexit
//             go/current/src/runtime/asm_amd64.s:1571
//     error B
//     mypkg.funcB
//             path/to/mypkg/mysource.go:26
//     mypkg.funcA
//             path/to/mypkg/mysource.go:22
//     mypkg.Example
//             path/to/mypkg/mysource.go:12
//     testing.runExample
//             go/current/src/testing/run_example.go:63
//     testing.runExamples
//             go/current/src/testing/example.go:44
//     testing.(*M).Run
//             go/current/src/testing/testing.go:1721
//     main.main
//             _testmain.go:49
//     runtime.main
//             go/current/src/runtime/proc.go:250
//     runtime.goexit
//             go/current/src/runtime/asm_amd64.s:1571
//     error A
//     mypkg.funcA
//             path/to/mypkg/mysource.go:22
//     mypkg.Example
//             path/to/mypkg/mysource.go:12
//     testing.runExample
//             go/current/src/testing/run_example.go:63
//     testing.runExamples
//             go/current/src/testing/example.go:44
//     testing.(*M).Run
//             go/current/src/testing/testing.go:1721
//     main.main
//             _testmain.go:49
//     runtime.main
//             go/current/src/runtime/proc.go:250
//     runtime.goexit
//             go/current/src/runtime/asm_amd64.s:1571
package stacktrace

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
)

// StackTrace is stack of Frames from innermost (newest) to outermost (oldest).
type StackTrace []Frame

// Frame contains parsial string represendations of errors.Frame.
type Frame struct {
	Message  string // error message
	FuncName string // function name
	Source   string // path of source file and line
}

func DumpErrors(err error, n int) {
	indent := strings.Repeat(" ", n)
	_, ok := err.(stackTracer)
	fmt.Printf("%s(%T; trace=%v)%v\n", indent, err, ok, err)

	if err := errors.Unwrap(err); err != nil {
		DumpErrors(err, n+1)
	}
}

// New extracts StackTrace from error err.
func New(err error) StackTrace {
	if err == nil {
		return nil
	}

	inner := errors.Unwrap(err)
	innerST := New(inner)

	st := []Frame{}
	if tracer, ok := err.(stackTracer); ok {
		msg := err.Error()
		for _, f := range tracer.StackTrace() {
			s := fmt.Sprintf("%+v", f)
			ss := strings.Split(s, "\n\t")
			if len(ss) > 1 {
				st = append(st, Frame{Message: msg, FuncName: ss[0], Source: ss[1]})
			}
		}
	} else {
		return innerST
	}
	if innerST == nil {
		return st
	}

	// merge
	idx := -1
loop:
	for i := 0; i < len(st); i++ {
		for j := 0; j < len(innerST); j++ {
			if st[i].FuncName == innerST[j].FuncName && st[i].Source == innerST[j].Source {
				idx = j
				break loop
			}
		}
	}
	if idx != -1 {
		innerST = innerST[:idx]
	}
	innerST = append(innerST, st...)

	return innerST
}

// String returns string representation of the StackTrace.
func (st StackTrace) String() string {
	return strings.Join(st.strings(true), "\n")
}

func (st StackTrace) strings(funcName bool) []string {
	if len(st) == 0 {
		return nil
	}

	ss := make([]string, 0, len(st))

	prevMsg := ""
	prevFunc := ""
	for _, f := range st {
		if f.Message != prevMsg {
			curr := f.Message
			prev := prevMsg
			if strings.HasSuffix(curr, ": "+prev) {
				curr = curr[:len(curr)-len(prev)-2]
			}
			ss = append(ss, curr)
		}

		if funcName && (f.Message != prevMsg || f.FuncName != prevFunc) {
			ss = append(ss, f.FuncName)
		}

		if f.Source != "" {
			ss = append(ss, "\t"+f.Source)
		}

		prevMsg = f.Message
		prevFunc = f.FuncName
	}

	return ss
}

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
// %s, %v  list frames
func (st StackTrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		ss := st.strings(s.Flag('+'))
		io.WriteString(s, strings.Join(ss, "\n"))

	case 's':
		io.WriteString(s, st.String())
	default:
		// nop
	}
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}
