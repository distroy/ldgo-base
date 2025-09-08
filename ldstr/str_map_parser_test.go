/*
 * Copyright (C) distroy
 */

package ldstr

import (
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
)

func TestStrMapParser(t *testing.T) {
	type subtest struct {
		text    string
		wantErr bool
		want    map[string]string
	}

	type test struct {
		tmpl     string
		splits   []string
		wantErr  bool
		subtests []subtest
	}

	tests := []test{
		{
			tmpl:    "{}: {user}",
			splits:  nil,
			wantErr: false,
			subtests: []subtest{{
				text:    "user: xxx",
				wantErr: false,
				want:    map[string]string{`user`: "xxx"},
			}},
		},
		{
			tmpl:    "#key#: #value#",
			splits:  []string{"#"},
			wantErr: false,
			subtests: []subtest{{
				text:    "user: xxx",
				wantErr: false,
				want: map[string]string{
					`key`:   "user",
					`value`: "xxx",
				},
			}},
		},
		{
			tmpl:    "{key}: {value}",
			splits:  []string{"{", "}"},
			wantErr: false,
			subtests: []subtest{{
				text:    "user: xxx",
				wantErr: false,
				want: map[string]string{
					`key`:   "user",
					`value`: "xxx",
				},
			}},
		},
		{
			tmpl:    "{key}{}: { {value}{ignore} }",
			splits:  nil,
			wantErr: false,
			subtests: []subtest{{
				text:    "user: { xxx }",
				wantErr: false,
				want: map[string]string{
					`key`:    "user",
					`value`:  "xxx",
					`ignore`: "",
				},
			}},
		},
		{
			tmpl:    "user: {user}, welcome {user}{action}",
			splits:  nil,
			wantErr: false,
			subtests: []subtest{
				{
					text:    "user: xxx, welcome xxx join",
					wantErr: false,
					want: map[string]string{
						`user`:   "xxx",
						`action`: " join",
					},
				},
				{
					text:    "user: xxx, welcome xyz join",
					wantErr: true,
					want:    nil,
				},
			},
		},
		{
			tmpl:    "user: {user}{text}{}{user}{action}",
			splits:  nil,
			wantErr: false,
			subtests: []subtest{
				{
					text:    "user: xxx, welcome xxx join",
					wantErr: false,
					want: map[string]string{
						`user`:   "xxx",
						`text`:   ", welcome ",
						`action`: " join",
					},
				},
				{
					text:    "user: xxx, welcome xyz join",
					wantErr: false,
					want: map[string]string{
						`user`:   "x",
						`text`:   "xx, welcome ",
						`action`: "yz join",
					},
				},
			},
		},
		{
			tmpl:    "{prefix}{user}{text}{}welcome {user}{action}",
			splits:  nil,
			wantErr: false,
			subtests: []subtest{{
				text:    "user: xxx, welcome xxx join",
				wantErr: false,
				want: map[string]string{
					`prefix`: "user: xxx, ",
					`user`:   "",
					`text`:   "",
					`action`: "xxx join",
				},
			}},
		},
		{
			tmpl:    "{prefix}{user}{text}{}welcome {user} join",
			splits:  nil,
			wantErr: false,
			subtests: []subtest{{
				text:    "user: xxx, welcome xxx join",
				wantErr: false,
				want: map[string]string{
					`prefix`: "user: ",
					`user`:   "xxx",
					`text`:   ", ",
				},
			}},
		},
	}

	subtestRun := func(c convey.C, parser *StrMapParser, subtests []subtest) {
		for _, stt := range subtests {
			c.Convey(stt.text, func(c convey.C) {
				got, err := parser.Parse(stt.text)
				if stt.wantErr {
					c.So(err, convey.ShouldNotBeNil)
				} else {
					c.So(err, convey.ShouldBeNil)
				}
				c.So(got, convey.ShouldResemble, stt.want)
			})
		}
	}

	convey.Convey(t.Name(), t, func(c convey.C) {
		for _, tt := range tests {
			c.Convey(tt.tmpl, func(c convey.C) {
				// log.Printf("")
				parser := &StrMapParser{}
				err := parser.Init(tt.tmpl, tt.splits...)
				// log.Printf(" === %s", mustMarshalJson(parser.fields))
				if tt.wantErr {
					c.So(err, convey.ShouldNotBeNil)
					return
				}

				c.So(err, convey.ShouldBeNil)
				c.So(parser.Template(), convey.ShouldEqual, tt.tmpl)

				subtestRun(c, parser, tt.subtests)
			})
		}
	})
}
