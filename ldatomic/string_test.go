/*
 * Copyright (C) distroy
 */

package ldatomic

import (
	"log"
	"os"
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestMain(m *testing.M) {
	log.SetFlags(log.Flags() | log.Lshortfile)
	os.Exit(m.Run())
}

func TestString(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		convey.So(NewString("abc").Load(), convey.ShouldEqual, "abc")

		p := String{}
		convey.So(p.Load(), convey.ShouldEqual, "")

		p.Store("123")
		convey.So(p.Load(), convey.ShouldEqual, "123")
		p.Store("abc")
		convey.So(p.Load(), convey.ShouldEqual, "abc")
		p.Store("")
		convey.So(p.String(), convey.ShouldEqual, "")

		convey.So(p.Swap("123"), convey.ShouldEqual, "")
		convey.So(p.Swap("abc"), convey.ShouldEqual, "123")
		convey.So(p.Swap(""), convey.ShouldEqual, "abc")
	})
}
