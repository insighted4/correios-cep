// Copyright 2023 The Correios CEP Admin Authors
//
// Licensed under the AGPL, Version 3.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.gnu.org/licenses/agpl-3.0.en.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"github.com/sirupsen/logrus"
)

// Error is typed error, it carries information and behavior as to what caused
// this error so that callers can implement logic around it.
type Error struct {
	// Kind categorizes errors into a smaller subset. This way we can generalize
	// what an error really is: such as "not found", "bad request", etc.
	// The official categories are HTTP status code but the ones we use are
	// imported into this package.
	Kind     int
	Op       Op
	Object   Obj
	Version  V
	Err      error
	Severity logrus.Level
}

// Error returns the underlying error's
// string message. The logger takes care
// of filling out the stack levels and
// extra information.
func (e Error) Error() string {
	return e.Err.Error()
}

// Is checks an error against a kind (shorthand).
func Is(err error, kind int) bool {
	if err == nil {
		return false
	}
	return Kind(err) == kind
}

// IsErr is a convenience wrapper around the std library errors.Is.
func IsErr(err, target error) bool {
	return errors.Is(err, target)
}

// AsErr is a convenience wrapper around the std library errors.As.
func AsErr(err error, target any) bool {
	return errors.As(err, target)
}

// Op describes any independent function or method. A series of operations
// forms a more readable stack trace.
type Op string

func (o Op) String() string {
	return string(o)
}

// Obj represents a module in an error
// this is so that we can distinguish
// an object from a regular error string or version.
type Obj string

// V represents a module version in an error
type V string

// E is a helper function to construct an Error type
// Operation always comes first, module path and version
// come second, they are optional. Args must have at least
// an error or a string to describe what exactly went wrong.
// You can optionally pass a Logrus severity to indicate
// the log level of an error based on the context it was constructed in.
func E(op Op, args ...interface{}) error {
	e := Error{Op: op}
	if len(args) == 0 {
		msg := "errors.E called with 0 args"
		_, file, line, ok := runtime.Caller(1)
		if ok {
			msg = fmt.Sprintf("%v - %v:%v", msg, file, line)
		}
		e.Err = errors.New(msg)
	}
	for _, a := range args {
		switch a := a.(type) {
		case error:
			e.Err = a
		case string:
			e.Err = errors.New(a)
		case Obj:
			e.Object = a
		case V:
			e.Version = a
		case logrus.Level:
			e.Severity = a
		case int:
			e.Kind = a
		}
	}
	if e.Err == nil {
		e.Err = errors.New(KindText(e))
	}
	return e
}

// Severity returns the log level of an error
// if none exists, then the level is Error because
// it is an unexpected.
func Severity(err error) logrus.Level {
	e, ok := err.(Error)
	if !ok {
		return logrus.ErrorLevel
	}

	// if there's no severity (0 is Panic level in logrus
	// which we should not use since cloud providers only have
	// debug, info, warn, and error) then look for the
	// child's severity.
	if e.Severity < logrus.ErrorLevel {
		return Severity(e.Err)
	}

	return e.Severity
}

// Expect is a helper that returns an Info level
// if the error has the expected kind, otherwise
// it returns an Error level.
func Expect(err error, kinds ...int) logrus.Level {
	for _, kind := range kinds {
		if Kind(err) == kind {
			return logrus.InfoLevel
		}
	}
	return logrus.ErrorLevel
}

// Kind recursively searches for the
// first error kind it finds.
func Kind(err error) int {
	e, ok := err.(Error)
	if !ok {
		return KindUnexpected
	}

	if e.Kind != 0 {
		return e.Kind
	}

	return Kind(e.Err)
}

// KindText returns a friendly string
// of the Kind type. Since we use http
// status codes to represent error kinds,
// this method just deferrs to the net/http
// text representations of statuses.
func KindText(err error) string {
	return http.StatusText(Kind(err))
}

// Ops aggregates the error's operation
// with all the embedded errors' operations.
// This way you can construct a queryable
// stack trace.
func Ops(err Error) []Op {
	ops := []Op{err.Op}
	for {
		embeddedErr, ok := err.Err.(Error)
		if !ok {
			break
		}

		ops = append(ops, embeddedErr.Op)
		err = embeddedErr
	}

	return ops
}
