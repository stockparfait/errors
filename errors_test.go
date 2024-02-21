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

package errors

import (
	"runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// rsn provides a fixed function name and line number in error annotations.
func rsn(r string) error {
	return Reason(r)
}

// ann provides a fixed function name and line number in error annotations.
func ann(e error, r string, args ...interface{}) error {
	return Annotate(e, r, args...)
}

type myError string

func (e myError) Error() string { return string(e) }

func fnA(mode string) (err error) {
	defer func() { err = FromPanic(recover()) }()
	fnB(mode)
	return nil
}

func fnB(mode string) {
	fnC(mode)
}

func fnC(mode string) {
	switch mode {
	case "error":
		ReasonPanic("error in %s", "fnC")
	case "panic":
		panic("panic in fnC")
	default:
		// no error or panic
	}
}

func TestErrors(t *testing.T) {
	Convey("Reason works", t, func() {
		e := rsn("because")
		So(e.Error(), ShouldContainSubstring,
			"errors_test.go:26: github.com/stockparfait/errors.rsn() because")
	})

	Convey("Annotate works", t, func() {
		Convey("annotates non-nil error", func() {
			e := ann(rsn("because"), "failed %s", "me")
			So(e.Error(), ShouldContainSubstring,
				"errors_test.go:31: github.com/stockparfait/errors.ann() failed me")
			So(e.Error(), ShouldContainSubstring,
				"errors_test.go:26: github.com/stockparfait/errors.rsn() because")
		})

		Convey("passes through nil error", func() {
			So(ann(nil, "you won't see this"), ShouldBeNil)
		})

		Convey("Is and As work", func() {
			err := myError("mine")
			annotated := ann(err, "annotated")
			So(Is(annotated, err), ShouldBeTrue)
			var err2 myError
			So(As(annotated, &err2), ShouldBeTrue)
			So(err2, ShouldEqual, err)
		})
	})

	Convey("Panic methods work", t, func() {

		Convey("trimFrames", func() {
			frames := []runtime.Frame{
				{Function: "aa"},
				{Function: "runtime.gopanic"},
				{Function: "cc"},
				{Function: "bb"},
				{Function: "runtime.main"},
				{Function: "runtime.top"},
			}
			Convey("trims the frames when possible", func() {
				So(trimFrames(frames), ShouldResemble, frames[2:4])
			})

			Convey("keep the frames when no matches found", func() {
				So(trimFrames(frames[2:4]), ShouldResemble, frames[2:4])
			})
		})

		Convey("recover an error panic", func() {
			err := fnA("error")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring,
				"errors_test.go:51: github.com/stockparfait/errors.fnC() error in fnC")
			So(err.Error(), ShouldContainSubstring,
				"errors_test.go:40 github.com/stockparfait/errors.fnA()")
			So(err.Error(), ShouldContainSubstring,
				"errors_test.go:45 github.com/stockparfait/errors.fnB()")
			So(err.Error(), ShouldContainSubstring,
				"errors_test.go:51 github.com/stockparfait/errors.fnC()")
		})

		Convey("re-raise non-error panic", func() {
			So(func() { fnA("panic") }, ShouldPanic)
		})

		Convey("no-op without panic", func() {
			So(fnA("none"), ShouldBeNil)
		})
	})
}
