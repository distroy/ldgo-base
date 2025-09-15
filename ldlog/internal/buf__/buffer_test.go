/*
 * Copyright (C) distroy
 */

package buf__

import (
	"testing"
	"time"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestBuffer(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		b := NewBuffer()
		defer b.Free()

		c.Convey("grow", func(c convey.C) {
			b.Grow(4 << 10)
			c.So(b.Len(), convey.ShouldEqual, 0)
			c.So(b.Cap(), convey.ShouldEqual, 4<<10)
		})

		c.Convey("write", func(c convey.C) {
			c.Convey("write", func(c convey.C) {
				n, err := b.Write(([]byte)("0123456789"))
				c.So(err, convey.ShouldBeNil)
				c.So(n, convey.ShouldEqual, 10)
				c.So(b.String(), convey.ShouldEqual, "0123456789")
			})
			c.Convey("string", func(c convey.C) {
				n, err := b.WriteString("0123456789")
				c.So(err, convey.ShouldBeNil)
				c.So(n, convey.ShouldEqual, 10)
				c.So(b.String(), convey.ShouldEqual, "0123456789")
			})
			c.Convey("byte", func(c convey.C) {
				err := b.WriteByte('a')
				c.So(err, convey.ShouldBeNil)
				c.So(b.String(), convey.ShouldEqual, "a")
			})
			c.Convey("bool", func(c convey.C) {
				n, err := b.WriteBool(true)
				c.So(err, convey.ShouldBeNil)
				c.So(n, convey.ShouldEqual, 4)
				c.So(b.String(), convey.ShouldEqual, "true")
			})
			c.Convey("quote", func(c convey.C) {
				n, err := b.WriteQuote("abc")
				c.So(err, convey.ShouldBeNil)
				c.So(n, convey.ShouldEqual, 5)
				c.So(b.String(), convey.ShouldEqual, `"abc"`)
			})
			c.Convey("int", func(c convey.C) {
				n, err := b.WriteInt(-1234)
				c.So(err, convey.ShouldBeNil)
				c.So(n, convey.ShouldEqual, 5)
				c.So(b.String(), convey.ShouldEqual, "-1234")
			})
			c.Convey("uint", func(c convey.C) {
				n, err := b.WriteUint(1234)
				c.So(err, convey.ShouldBeNil)
				c.So(n, convey.ShouldEqual, 4)
				c.So(b.String(), convey.ShouldEqual, "1234")
			})
			c.Convey("float", func(c convey.C) {
				n, err := b.WriteFloat(-1.234, 64)
				c.So(err, convey.ShouldBeNil)
				c.So(n, convey.ShouldEqual, 6)
				c.So(b.String(), convey.ShouldEqual, "-1.234")
			})
			c.Convey("complex", func(c convey.C) {
				n, err := b.WriteComplex(complex(1, -234), 64)
				c.So(err, convey.ShouldBeNil)
				c.So(n, convey.ShouldEqual, 8)
				c.So(b.String(), convey.ShouldEqual, "(1-234i)")
			})
			c.Convey("duration", func(c convey.C) {
				n, err := b.WriteDuration(time.Minute + time.Second*10)
				c.So(err, convey.ShouldBeNil)
				c.So(n, convey.ShouldEqual, 5)
				c.So(b.String(), convey.ShouldEqual, "1m10s")
			})
			c.Convey("time", func(c convey.C) {
				n, err := b.WriteTime(time.Unix(1629610258, 0), "")
				c.So(err, convey.ShouldBeNil)
				c.So(n, convey.ShouldEqual, 28)
				c.So(b.String(), convey.ShouldEqual, "2021-08-22T13:30:58.000+0800")
			})
			c.Convey("json", func(c convey.C) {
				v := &struct {
					Int  int    `json:"int"`
					Uint uint   `json:"uint"`
					Str  string `json:"str"`
				}{
					Int:  -1234,
					Uint: 1234,
					Str:  "abc",
				}
				n, err := b.WriteJson(v)
				c.So(err, convey.ShouldBeNil)
				c.So(n, convey.ShouldEqual, 37)
				c.So(b.String(), convey.ShouldEqual, `{"int":-1234,"uint":1234,"str":"abc"}`)
			})
		})

		c.Convey("append", func(c convey.C) {
			c.Convey("byte", func(c convey.C) {
				b.AppendByte('a')
				c.So(b.String(), convey.ShouldEqual, "a")
			})
			c.Convey("bytes", func(c convey.C) {
				b.AppendBytes(([]byte)("0123456789"))
				c.So(b.String(), convey.ShouldEqual, "0123456789")
				c.So(b.Bytes(), convey.ShouldResemble, ([]byte)("0123456789"))
			})
			c.Convey("string", func(c convey.C) {
				b.AppendString("0123456789")
				c.So(b.String(), convey.ShouldEqual, "0123456789")
			})
		})
	})
}
