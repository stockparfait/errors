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

func rsn(r string) error {
	return Reason(r)
}

func ann(e error, r string, args ...interface{}) error {
	return Annotate(e, r, args...)
}

func TestErrors(t *testing.T) {
	Convey("Reason works", t, func() {
		e := rsn("because")
		So(e.Error(), ShouldContainSubstring,
			"errors_test.go:24: github.com/stockparfait/errors.rsn() because")
	})

	Convey("Annotate works", t, func() {
		e := ann(rsn("because"), "failed %s", "me")
		So(e.Error(), ShouldContainSubstring,
			"errors_test.go:28: github.com/stockparfait/errors.ann() failed me")
		So(e.Error(), ShouldContainSubstring,
			"errors_test.go:24: github.com/stockparfait/errors.rsn() because")
	})
}
