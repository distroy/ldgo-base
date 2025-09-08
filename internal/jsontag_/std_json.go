/*
 * Copyright (C) distroy
 */

package jsontag_

// import (
// 	"bytes"
// 	_ "encoding/json"
// 	"reflect"
// 	_ "unsafe"
// )
//
// func Get(typ reflect.Type) Struct {
// 	return cachedTypeFields(typ)
// }
//
// type (
// 	Struct = structFields
// 	Field  = field
// )
//
// type encodeState struct {
// 	bytes.Buffer // accumulated output
// 	scratch      [64]byte
//
// 	// Keep track of what pointers we've seen in the current recursive call
// 	// path, to avoid cycles that could lead to a stack overflow. Only do
// 	// the relatively expensive map operations if ptrLevel is larger than
// 	// startDetectingCyclesAfter, so that we skip the work if we're within a
// 	// reasonable amount of nested pointers deep.
// 	ptrLevel uint
// 	ptrSeen  map[any]struct{}
// }
// type encOpts struct {
// 	// quoted causes primitive fields to be encoded inside JSON strings.
// 	quoted bool
// 	// escapeHTML causes '<', '>', and '&' to be escaped in JSON strings.
// 	escapeHTML bool
// }
//
// type encoderFunc func(e *encodeState, v reflect.Value, opts encOpts)
//
// type field struct {
// 	Name      string
// 	NameBytes []byte                 // []byte(name)
// 	EqualFold func(s, t []byte) bool // bytes.EqualFold or equivalent
//
// 	NameNonEsc  string // `"` + name + `":`
// 	NameEscHTML string // `"` + HTMLEscape(name) + `":`
//
// 	Tag       bool
// 	Index     []int
// 	Typ       reflect.Type
// 	OmitEmpty bool
// 	Quoted    bool
//
// 	Encoder encoderFunc
// }
// type structFields struct {
// 	List      []field
// 	NameIndex map[string]int
// }
//
// //go:linkname replace encoding/json.cachedTypeFields
// func cachedTypeFields(t reflect.Type) structFields
