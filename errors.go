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
	"strings"
)

// annotatedError annotates the original error with the current message.
type annotatedError struct {
	orig error
	curr string
}

// Error implements error.
func (e annotatedError) Error() string {
	if e.orig == nil {
		return e.curr
	}
	return fmt.Sprintf("%s\n%s", e.curr, e.orig.Error())
}

// Unwrap returns the original error being annotated. See also As and Is methods.
func (e annotatedError) Unwrap() error {
	return e.orig
}

// annotate must be called from ReasonStack or AnnotateStack only.
func annotate(stack int, s string, args ...any) string {
	// Frame 2 is the caller of Reason / Annotate.
	pc, filename, line, ok := runtime.Caller(stack)
	a := "ERROR: ???: "
	if ok {
		a = fmt.Sprintf("ERROR: %s:%d: %s() ", filename, line, runtime.FuncForPC(pc).Name())
	}
	return a + fmt.Sprintf(s, args...)
}

// ReasonStack returns an error annotated with location `stack` levels up, and
// message. Its arguments are the same as for fmt.Printf.
func ReasonStack(stack int, s string, args ...any) error {
	return &annotatedError{curr: annotate(stack, s, args...)}
}

// AnnotateStack annotates the existing error with location `stack` levels up,
// and message, formatted as fmt.Printf(s, args...). If the original error is
// nil, returns nil.
func AnnotateStack(e error, stack int, s string, args ...any) error {
	if e == nil {
		return nil
	}
	return &annotatedError{orig: e, curr: annotate(stack, s, args...)}
}

// Reason returns an error annotated with location and message. Its arguments
// are the same as for fmt.Printf.
func Reason(s string, args ...any) error {
	return ReasonStack(3, s, args...)
}

// Annotate the existing error with location and message, formatted as
// fmt.Printf(s, args...). If the original error is nil, returns nil.
func Annotate(e error, s string, args ...any) error {
	return AnnotateStack(e, 3, s, args...)
}

// ReasonPanic is equivalent to panic(Reason(s, args...)).  This allows using
// panic as an exception for error handling.  See also FromPanic for converting
// such panic back into error.
func ReasonPanic(s string, args ...any) {
	panic(ReasonStack(3, s, args...))
}

// trimFrames to keep only the portion from panic to the top user main(). If in
// doubt, keep the frames.
func trimFrames(frames []runtime.Frame) []runtime.Frame {
	for i, f := range frames {
		if f.Function == "runtime.gopanic" {
			frames = frames[i+1:]
			break
		}
	}
	for i, f := range frames {
		if f.Function == "runtime.main" {
			frames = frames[:i]
			break
		}
	}
	return frames
}

// FromPanic converts an intentional panic back to error and annotates it with
// the panic call stack. Other panics are re-raised. It is intended to be used
// in defer:
//
//	func Foo() (err error) {
//	  defer func() { err = FromPanic(recover()) }()
//	  // Foo body, may panic on error
//	}
func FromPanic(p any) error {
	if p == nil {
		return nil
	}
	if err, ok := p.(*annotatedError); ok {
		pc := make([]uintptr, 20)
		n := runtime.Callers(3, pc)
		if n == 0 { // shouldn't happen, defensive code
			return err
		}
		pc = pc[:n] // use only valid pcs
		framesIter := runtime.CallersFrames(pc)

		frames := []runtime.Frame{}
		for {
			frame, more := framesIter.Next()
			frames = append(frames, frame)
			if !more {
				break
			}
		}
		frames = trimFrames(frames)
		// Invert frames in place.
		for l, h := 0, len(frames)-1; l < h; l, h = l+1, h-1 {
			frames[l], frames[h] = frames[h], frames[l]
		}
		traces := []string{}
		for _, frame := range frames {
			traces = append(traces, fmt.Sprintf("PANIC: %s:%d %s()",
				frame.File, frame.Line, frame.Function))
		}
		if len(traces) == 0 { // no panic stack found, defensive code
			return err
		}
		return &annotatedError{orig: err, curr: strings.Join(traces, "\n")}
	}
	// Re-raise all other panics.
	panic(p)
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
func As(err error, target any) bool {
	return errors.As(err, target)
}
