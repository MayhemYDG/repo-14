# (Yet another) error handling library for Go

When an unhandled error occurs, it is usually caught and logged at some top layer within an application. For example, in a CLI this might be in `main()`. For a request-response service, this is usually in a middleware layer.

A good error message should include enough information to allow a developer to quickly understand and debug the error. The two most important pieces of information that aid in this are:

1. The location in the code (file, line number, and, optionally function name) where the error originated, and
2. What action was in progress when the error happened. This should have enough context to  identify the specific entities which were involved.

However, by default in Go, errors do not include any context information whatsoever. Therefore it is up to the developers to choose the mechanism we want to use for collecting and propagating that context along with the error.

`cerrors` is one way that this can be done. It works like this:

- When an error is created or wrapped, attaches information about the source location, including the file, line number, and function name. This achieves (1) above.

- Encourages developers to wrap unhandled errors with a description of what was happening at the time the error was encountered. This achieves (2) above.

This is similar in spirit to the original ["Go 2" proposal for the error handling improvements](https://go.googlesource.com/proposal/+/master/design/29934-error-values.md). It has been only partially implemented by Go so far (the stack frame information is not currently tracked). It is also similar to how [pkg/errors](https://github.com/pkg/errors) worked in its earlier iterations.

## Usage

```go
import "github.com/1debit/cerrors"
// ...

// create a new error
err := cerrors.New("total fail")

// wrap the error a few times, including a non-cerrors wrap
err = cerrors.Wrapf(err, "while doing xyz with (foo=%s)", foo)
err = fmt.Errorf("standard Go wrap: %w", err)
err = cerrors.Wrap(err, "during abc")

// get the short form error message
fmt.Println(err.Error())
// =>
// during abc: standard Go wrap: while doing xyz with (foo=blah): total fail

// print the long form error
fmt.Printf("%#v\n", err)
// =>
// during abc
//     main.main
//         main.go:19
//   - standard Go wrap
//   - while doing xyz with (foo=blah)
//     main.main
//         main.go:18
//   - total fail
//     main.main
//         main.go:15
```

The long-form error message shown above includes the stack, with each layer showing not only the function name, the filename and the line number, but also the description of what was happening at the time. **This helps the reader of the message to determine the cause of the issue without having to delve into the source**. However, it *only* includes the frames where the error is wrapped using `cerrors.Wrap`. This gives the developer the option to skip stack frames, such as helper functions, which would not be informative when debugging.
