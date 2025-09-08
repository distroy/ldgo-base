/*
 * Copyright (C) distroy
 */

package ldctx

import (
	"bytes"
	"testing"
	"time"
	_ "unsafe"

	"github.com/distroy/ldgo-base/3rd/convey"
	"github.com/distroy/ldgo-base/ldhook"
	"github.com/distroy/ldgo-base/ldlog"
)

func TestLog(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		patches := ldhook.NewPatches()
		defer patches.Reset()
		patches.Applys([]ldhook.Hook{
			ldhook.FuncHook{
				Target: time.Now,
				Double: ldhook.Values{time.Unix(1629610258, 0)},
			},
		})

		b := bytes.NewBuffer(make([]byte, 0, 1024))
		l := ldlog.New(ldlog.NewHandler(b, nil), ldlog.SetLevel(ldlog.LevelTrace))
		l = l.With(ldlog.String("abc", "xxx"))

		ctx := WithLogger(nil, l)

		c.Convey("LogT", func(c convey.C) {
			LogT(ctx, "test")
			c.So(b.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|TRACE|-|ldctx/ctx_log_test.go:36|test,abc=xxx\n")
		})
		c.Convey("LogD", func(c convey.C) {
			LogD(ctx, "test")
			c.So(b.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|DEBUG|-|ldctx/ctx_log_test.go:41|test,abc=xxx\n")
		})
		c.Convey("LogI", func(c convey.C) {
			LogI(ctx, "test")
			c.So(b.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|INFO|-|ldctx/ctx_log_test.go:46|test,abc=xxx\n")
		})
		c.Convey("LogW", func(c convey.C) {
			LogW(ctx, "test")
			c.So(b.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|WARN|-|ldctx/ctx_log_test.go:51|test,abc=xxx\n")
		})
		c.Convey("LogE", func(c convey.C) {
			LogE(ctx, "test")
			c.So(b.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|ERROR|-|ldctx/ctx_log_test.go:56|test,abc=xxx\n")
		})

		c.Convey("LogTf", func(c convey.C) {
			LogTf(ctx, "test")
			c.So(b.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|TRACE|-|ldctx/ctx_log_test.go:62|test,abc=xxx\n")
		})
		c.Convey("LogDf", func(c convey.C) {
			LogDf(ctx, "test")
			c.So(b.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|DEBUG|-|ldctx/ctx_log_test.go:67|test,abc=xxx\n")
		})
		c.Convey("LogIf", func(c convey.C) {
			LogIf(ctx, "test")
			c.So(b.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|INFO|-|ldctx/ctx_log_test.go:72|test,abc=xxx\n")
		})
		c.Convey("LogWf", func(c convey.C) {
			LogWf(ctx, "test")
			c.So(b.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|WARN|-|ldctx/ctx_log_test.go:77|test,abc=xxx\n")
		})
		c.Convey("LogEf", func(c convey.C) {
			LogEf(ctx, "test")
			c.So(b.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|ERROR|-|ldctx/ctx_log_test.go:82|test,abc=xxx\n")
		})
	})
}
