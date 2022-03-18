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
//   func MyFunc(x int) error {
//     if x < 0 {
//       return errors.Reason("x = %d is negative", x)
//     }
//     return nil
//
//   if err := MyFunc(val); err != nil {
//     return errors.Annotate(err, "cannot use %d", val)
//   }
package errors

import (
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
	return fmt.Sprintf("ERROR: %s:\n%s", e.curr, e.orig.Error())
}

// annotate must be called from Reason or Annotate only.
func annotate(s string, args ...interface{}) string {
	// Frame 2 is the caller of Reason / Annotate.
	pc, filename, line, ok := runtime.Caller(2)
	a := "???: "
	if ok {
		a = fmt.Sprintf("%s:%d: %s() ", filename, line, runtime.FuncForPC(pc).Name())
	}
	return a + fmt.Sprintf(s, args...)
}

// Reason returns an annotated error. Its arguments are the same as for
// fmt.Printf.
func Reason(s string, args ...interface{}) error {
	return fmt.Errorf(annotate(s, args...))
}

func Annotate(e error, s string, args ...interface{}) error {
	return &annotatedError{orig: e, curr: annotate(s, args...)}
}
