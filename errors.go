// Copyright 2022 Stock Parfait

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package errors implements custom error reporting.
//
// Example usage:
//
//	func MyFunc(x int) error {
//	  if x < 0 {
//	    return errors.Reason("x = %d is negative", x)
//	  }
//	  return nil
//
//	if err := MyFunc(val); err != nil {
//	  return errors.Annotate(err, "cannot use %d", val)
//	}
package errors

import (
	"errors"
	"fmt"
	"runtime"
)

// annotatedError annotates the original error with the current message.
type annotatedError struct {
	orig error
	curr string
}

// Error implements error.
func (e annotatedError) Error() string {
	return fmt.Sprintf("%s\n%s", e.curr, e.orig.Error())
}

// Unwrap returns the original error being annotated. See also As and Is methods.
func (e annotatedError) Unwrap() error {
	return e.orig
}

// annotate must be called from Reason or Annotate only.
func annotate(stack int, s string, args ...interface{}) string {
	// Frame 2 is the caller of Reason / Annotate.
	pc, filename, line, ok := runtime.Caller(stack)
	a := "ERROR: ???: "
	if ok {
		a = fmt.Sprintf("ERROR: %s:%d: %s() ", filename, line, runtime.FuncForPC(pc).Name())
	}
	return a + fmt.Sprintf(s, args...)
}

// Reason returns an error annotated with location and message. Its arguments
// are the same as for fmt.Printf.
func Reason(s string, args ...interface{}) error {
	return ReasonStack(3, s, args...)
}

// Annotate the existing error with location and message, formatted as
// fmt.Printf(s, args...). If the original error is nil, returns nil.
func Annotate(e error, s string, args ...interface{}) error {
	return AnnotateStack(e, 3, s, args...)
}

// ReasonStack returns an error annotated with location `stack` levels up, and
// message. Its arguments are the same as for fmt.Printf.
func ReasonStack(stack int, s string, args ...interface{}) error {
	return fmt.Errorf(annotate(stack, s, args...))
}

// AnnotateStack annotates the existing error with location `stack` levels up,
// and message, formatted as fmt.Printf(s, args...). If the original error is
// nil, returns nil.
func AnnotateStack(e error, stack int, s string, args ...interface{}) error {
	if e == nil {
		return nil
	}
	return &annotatedError{orig: e, curr: annotate(stack, s, args...)}
}

// Is reports whether any error in err's "Unwrap" chain matches target.
//
// It is exactly as Go's errors.Is method, and is provided to match the
// functionality.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As sets the target to the first applicable value in err's "Unwrap" chain.
//
// It is exactly as Go's errors.As method, and is provided to match the
// functionality.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
