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

func TestErrors(t *testing.T) {
	Convey("Reason works", t, func() {
		e := rsn("because")
		So(e.Error(), ShouldContainSubstring,
			"errors_test.go:25: github.com/stockparfait/errors.rsn() because")
	})

	Convey("Annotate works", t, func() {
		Convey("annotates non-nil error", func() {
			e := ann(rsn("because"), "failed %s", "me")
			So(e.Error(), ShouldContainSubstring,
				"errors_test.go:30: github.com/stockparfait/errors.ann() failed me")
			So(e.Error(), ShouldContainSubstring,
				"errors_test.go:25: github.com/stockparfait/errors.rsn() because")
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
}
