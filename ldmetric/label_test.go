/*
 * Copyright (C) distroy
 */

package ldmetric

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/distroy/ldgo-base/3rd/convey"
	"github.com/distroy/ldgo-base/ldconv"
	"github.com/distroy/ldgo-base/lderr"
	"github.com/distroy/ldgo-base/ldptr"
)

type testStringer []byte

func (s testStringer) String() string { return ldconv.BytesToStrUnsafe(s) }

func TestToLabels(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		type SubLabelStruct struct {
			Bool   *bool
			String string
		}

		c.Convey("anonymous struct", func(c convey.C) {
			type LabelStruct struct {
				SubLabelStruct

				Int        int          `ldmetric:"name:int"`
				UintPtr    *uint        `ldmetric:""`
				Float      float32      `ldmetric:""`
				Time       time.Time    `ldmetric:""`
				Error      error        ``
				Stringer   fmt.Stringer ``
				unexported int          ``
				Unexported int          `ldmetric:"-"`
			}

			in := &LabelStruct{
				SubLabelStruct: SubLabelStruct{
					Bool:   nil,
					String: "aaa",
				},
				Int:        100,
				UintPtr:    ldptr.New[uint](200),
				Float:      300,
				Time:       time.Unix(1723707841, 0).In(time.FixedZone("", 0)),
				unexported: 400,
				Unexported: 500,
			}
			out := ToLabels(in)
			c.So(out, convey.ShouldResemble, map[string]string{
				`bool`:          nilStr,
				`string`:        "aaa",
				`int`:           "100",
				`uint_ptr`:      "200",
				`float`:         "300",
				`time`:          "2024-08-15T07:44:01Z",
				`error_code`:    "0",
				`error_message`: lderr.ErrSuccess.Error(),
				`stringer`:      nilStr,
			})
		})

		c.Convey("struct field", func(c convey.C) {
			type LabelStruct struct {
				SubLabelStruct0 *SubLabelStruct ``
				SubLabelStruct1 *SubLabelStruct ``
				Int             int             `ldmetric:"name:int"`
				UintPtr         *uint           `ldmetric:""`
				Float           float32         `ldmetric:""`
				Error           error           `ldmetric:""`
				Stringer        fmt.Stringer    ``
			}
			in := &LabelStruct{
				SubLabelStruct0: &SubLabelStruct{
					Bool:   ldptr.New(true),
					String: "aaa",
				},
				Int:      100,
				UintPtr:  ldptr.New[uint](200),
				Float:    300,
				Error:    lderr.ErrUnkown,
				Stringer: testStringer("bbb"),
			}
			out := ToLabels(in)
			c.So(out, convey.ShouldResemble, map[string]string{
				`sub_label_struct0.bool`:   "true",
				`sub_label_struct0.string`: "aaa",
				`sub_label_struct1.bool`:   nilStr,
				`sub_label_struct1.string`: "",
				`int`:                      "100",
				`uint_ptr`:                 "200",
				`float`:                    "300",
				`error_code`:               strconv.Itoa(lderr.ErrUnkown.Code()),
				`error_message`:            lderr.ErrUnkown.Error(),
				`stringer`:                 "bbb",
			})
		})
	})
}
