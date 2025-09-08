/*
 * Copyright (C) distroy
 */

package ldsheet

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func Test_splitStringWord(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(splitStringWord("HttpUrl"), convey.ShouldEqual, "Http Url")
		convey.So(splitStringWord("HTTPAddress"), convey.ShouldEqual, "HTTP Address")
		convey.So(splitStringWord("UserAccountID"), convey.ShouldEqual, "User Account ID")
	})
}
