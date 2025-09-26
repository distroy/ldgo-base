/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"io"
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestAddLevelTo(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		l := testNewLogger(io.Discard)

		l = l.WithOptions(SetLevel(LevelInfo))
		c.So(l.Level(), convey.ShouldEqual, LevelInfo)

		l = l.WithOptions(AddLevelTo(LevelTrace))
		c.So(l.Level(), convey.ShouldEqual, LevelInfo)

		l = l.WithOptions(AddLevelTo(LevelWarn))
		c.So(l.Level(), convey.ShouldEqual, LevelWarn)
	})
}
