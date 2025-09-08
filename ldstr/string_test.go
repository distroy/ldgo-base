/*
 * Copyright (C) distroy
 */

package ldstr

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestStrMapReplace(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		m := map[string]string{
			"a":   "A",
			"b":   "B",
			"env": "test",
		}

		convey.So(StrMapReplace("a}{{a}{abc}{b}{}.{{a}}", m, "{", "}"), convey.ShouldEqual, "a}{A{abc}B{}.{A}")
		convey.So(StrMapReplace("./config_{env}.ini", m, "{", "}"), convey.ShouldEqual, "./config_test.ini")
	})
}

func TestStrIMapReplace(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		s0 := "a}{{a}{abc}{b}{}.{{a}}"
		m0 := map[string]string{
			"a": "A",
			"b": "B",
		}
		convey.So(StrIMapReplace(s0, m0, "{", "}"), convey.ShouldEqual, "a}{A{abc}B{}.{A}")

		m1 := map[string]int{
			"a": 1,
			"b": 2,
		}
		convey.So(StrIMapReplace(s0, m1, "{", "}"), convey.ShouldEqual, "a}{1{abc}2{}.{1}")
	})
}

func TestStrMapParse(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		var (
			tmpl string
			text string
			res  map[string]string
			err  error
		)
		tmpl = "{}: {user}"
		text = "user: xxx"
		res, err = StrMapParse(tmpl, text)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldResemble, map[string]string{`user`: "xxx"})
	})
}
