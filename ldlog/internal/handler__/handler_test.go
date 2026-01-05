/*
 * Copyright (C) distroy
 */

package handler__

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"
	"unsafe"

	"github.com/distroy/ldgo-base/3rd/convey"
	"github.com/distroy/ldgo-base/ldhook"
	"github.com/distroy/ldgo-base/ldptr"
)

func testNewLogger(w io.Writer) *slog.Logger {
	return slog.New(NewHandler(w, &Options{
		Caller: true,
		Level:  LevelDebug,
	}))
}

func TestHandler(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		patches := ldhook.NewPatches()
		defer patches.Reset()
		patches.Applys([]ldhook.Hook{
			ldhook.FuncHook{
				Target: time.Now,
				Double: ldhook.Values{time.Unix(1629610258, 0)},
			},
		})

		type LoggerValue struct {
			Name string
		}
		writer := bytes.NewBuffer(nil)
		l := testNewLogger(writer)
		l = l.With(slog.String("abc", "xxx"))

		c.Convey("error", func(c convey.C) {
			l.Error("error message")
			c.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|ERROR|-|handler__/handler_test.go:47|error message,abc=xxx\n")
		})

		c.Convey("warn", func(c convey.C) {
			l.Warn("warn message", slog.Int("int", 123))
			c.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|WARN|-|handler__/handler_test.go:53|warn message,abc=xxx|int=123\n")
		})

		c.Convey("info", func(c convey.C) {
			l.Info("info message", "dur", (10 * time.Millisecond))
			c.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|INFO|-|handler__/handler_test.go:59|info message,abc=xxx|dur=10ms\n")
		})

		c.Convey("warnln", func(c convey.C) {
			l.Warn("warnln message", "*int", ldptr.New(1234), "map", map[string]string{"a": "b"})
			c.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|WARN|-|handler__/handler_test.go:65|warnln message,abc=xxx|*int=1234,map={\"a\":\"b\"}\n")
		})

		c.Convey("infoln", func(c convey.C) {
			l.Info("infoln message", "obj", &LoggerValue{Name: "abc"}, "slice", []any{ldptr.New(1234), "*int", (*int)(nil)})
			c.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|INFO|-|handler__/handler_test.go:71|infoln message,abc=xxx|obj={\"Name\":\"abc\"},slice=[1234,\"*int\",null]\n")
		})

		c.Convey("errorln", func(c convey.C) {
			l.Error("errorln message", "ptr", (*LoggerValue)(nil), "ptr2", unsafe.Pointer(uintptr(0x2345)))
			c.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|ERROR|-|handler__/handler_test.go:77|errorln message,abc=xxx|ptr=null,ptr2=0x2345\n")
		})

		c.Convey("map", func(c convey.C) {
			l.Warn("warnln message", "*int", ldptr.New(1234), "map", map[string]any{
				"a":   "b",
				"100": 124,
				"10":  234,
			})
			c.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|WARN|-|handler__/handler_test.go:83|warnln message,abc=xxx|*int=1234,map={\"10\":234,\"100\":124,\"a\":\"b\"}\n")
		})

		c.Convey("errorf", func(c convey.C) {
			l.Error("errorf message", "int", 1234)
			c.So(writer.String(), convey.ShouldEqual,
				"2021-08-22T13:30:58.000+0800|ERROR|-|handler__/handler_test.go:93|errorf message,abc=xxx|int=1234\n")
		})

		c.Convey("debug", func(c convey.C) {
			type Object struct {
				Int int    `json:"int"`
				Str string `json:"str"`
			}
			obj := &Object{
				Int: 123,
				Str: "abc",
			}
			l.Debug("error message", slog.Any("obj", obj))
			c.So(writer.String(), convey.ShouldEqual,
				`2021-08-22T13:30:58.000+0800|DEBUG|-|handler__/handler_test.go:107|error message,abc=xxx|obj={"int":123,"str":"abc"}`+"\n")
		})
	})
}

func testSLogPrint(log *slog.Logger) {
	log = log.With(slog.Int("int", 123))
	log = log.WithGroup("g")
	log.Debug("test", slog.String("str", "abc"), slog.Group("g1", slog.Int("int1", 234)),
		slog.Group("", slog.String("str1", "xyz")))
}

func TestHandlerPrint(t *testing.T) {
	opts := &Options{Caller: true, Level: slog.LevelDebug}
	buf := bytes.NewBuffer(nil)
	log := slog.New(NewHandler(buf, opts))
	testSLogPrint(log)
	fmt.Printf("%s\n", buf.Bytes())
}

func TestSLogTextHandler(t *testing.T) {
	opts := &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}
	buf := bytes.NewBuffer(nil)
	log := slog.New(slog.NewTextHandler(buf, opts))
	testSLogPrint(log)
	fmt.Printf("%s\n", buf.Bytes())
}
