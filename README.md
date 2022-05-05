Package stacktrace extracts merged stacktrace from wrapped errors.

[![Go Report Card](https://goreportcard.com/badge/github.com/shu-go/stacktrace)](https://goreportcard.com/report/github.com/shu-go/stacktrace)
![MIT License](https://img.shields.io/badge/License-MIT-blue)

# Usage

    func funcA() error {
    	return errors.Wrap(funcB(), "error A")
    }

    func funcB() error {
    	return errors.Wrap(funcC(), "error B")
    }

    func funcC() error {
    	return errors.New("error C")
    }

## Case 1: Simple output

    fmt.Printf("%v\n", stacktrace.New(err))

vvv

    error C
            path/to/mypkg/stacktrace_test.go:30
    error B
            path/to/mypkg/stacktrace_test.go:26
    error A
            path/to/mypkg/stacktrace_test.go:22
            path/to/mypkg/stacktrace_test.go:12
            go/current/src/testing/run_example.go:63
            go/current/src/testing/example.go:44
            go/current/src/testing/testing.go:1721
            _testmain.go:49
            go/current/src/runtime/proc.go:250
            go/current/src/runtime/asm_amd64.s:1571

## Case 2: With function names

    fmt.Printf("%+v\n", stacktrace.New(err))

vvv

    error C
    mypkg.funcC
            path/to/mypkg/mysource.go:30
    error B
    mypkg.funcB
            path/to/mypkg/mysource.go:26
    error A
    mypkg.funcA
            path/to/mypkg/mysource.go:22
    mypkg.Example
            path/to/mypkg/mysource.go:12
    testing.runExample
            go/current/src/testing/run_example.go:63
    testing.runExamples
            go/current/src/testing/example.go:44
    testing.(*M).Run
            go/current/src/testing/testing.go:1721
    main.main
            _testmain.go:49
    runtime.main
            go/current/src/runtime/proc.go:250
    runtime.goexit
            go/current/src/runtime/asm_amd64.s:1571

## Case 3: (without this package) fmt.Printf("%+v\n", a())

Redundant output.

    error C
    mypkg.funcC
            path/to/mypkg/mysource.go:30
    mypkg.funcB
            path/to/mypkg/mysource.go:26
    mypkg.funcA
            path/to/mypkg/mysource.go:22
    mypkg.Example
            path/to/mypkg/mysource.go:12
    testing.runExample
            go/current/src/testing/run_example.go:63
    testing.runExamples
            go/current/src/testing/example.go:44
    testing.(*M).Run
            go/current/src/testing/testing.go:1721
    main.main
            _testmain.go:49
    runtime.main
            go/current/src/runtime/proc.go:250
    runtime.goexit
            go/current/src/runtime/asm_amd64.s:1571
    error B
    mypkg.funcB
            path/to/mypkg/mysource.go:26
    mypkg.funcA
            path/to/mypkg/mysource.go:22
    mypkg.Example
            path/to/mypkg/mysource.go:12
    testing.runExample
            go/current/src/testing/run_example.go:63
    testing.runExamples
            go/current/src/testing/example.go:44
    testing.(*M).Run
            go/current/src/testing/testing.go:1721
    main.main
            _testmain.go:49
    runtime.main
            go/current/src/runtime/proc.go:250
    runtime.goexit
            go/current/src/runtime/asm_amd64.s:1571
    error A
    mypkg.funcA
            path/to/mypkg/mysource.go:22
    mypkg.Example
            path/to/mypkg/mysource.go:12
    testing.runExample
            go/current/src/testing/run_example.go:63
    testing.runExamples
            go/current/src/testing/example.go:44
    testing.(*M).Run
            go/current/src/testing/testing.go:1721
    main.main
            _testmain.go:49
    runtime.main
            go/current/src/runtime/proc.go:250
    runtime.goexit
            go/current/src/runtime/asm_amd64.s:1571
