/*
 * Copyright (C) distroy
 */

package ldsheet

import (
	"fmt"
	"io"
	"reflect"
)

const (
	_TAG_NAME = "ldsheet"
)

type Reader struct {
	Reader LineReader
	model  model
}

func NewReader(r LineReader) *Reader {
	return &Reader{
		Reader: r,
	}
}

func (r *Reader) Read(d interface{}) error {
	if err := r.init(d); err != nil {
		return err
	}

	if err := r.readHeader(); err != nil {
		return err
	}

	if err := r.readData(d); err != nil {
		return err
	}

	return nil
}

func (r *Reader) ReadAll(data interface{}) error {
	pVal := reflect.ValueOf(data)
	pTyp := pVal.Type()
	if pTyp.Kind() != reflect.Ptr && pTyp.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("[ldsheet] the param of read all must be pointer to slice. type:%s", pTyp.String())
	}

	sVal := pVal.Elem()
	eTyp := pTyp.Elem().Elem()
	if (eTyp.Kind() != reflect.Ptr || eTyp.Elem().Kind() != reflect.Struct) && eTyp.Kind() != reflect.Struct {
		return fmt.Errorf("[ldsheet] the element of slice must be the struct or the pointer to struct. type:%s", eTyp.String())
	}

	fnAppend := func(v reflect.Value) { sVal.Set(reflect.Append(sVal, v.Elem())) }
	typ := eTyp
	if eTyp.Kind() == reflect.Ptr {
		typ = eTyp.Elem()
		fnAppend = func(v reflect.Value) { sVal.Set(reflect.Append(sVal, v)) }
	}

	if err := r.init(reflect.New(typ).Interface()); err != nil {
		return err
	}

	if err := r.readHeader(); err != nil {
		return err
	}

	for {
		val := reflect.New(typ)
		err := r.readDataValue(val.Elem())
		switch err {
		case nil:
			fnAppend(val)

		case io.EOF:
			return nil

		default:
			return fmt.Errorf("[ldsheet] read all data fail. err:%s", err)
		}
	}
}

func (r *Reader) init(d interface{}) error {
	m := &r.model
	if m.Inited {
		return nil
	}

	typ := reflect.TypeOf(d)
	if typ.Kind() != reflect.Ptr && typ.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("[ldsheet] the data must be the pointer to struct. type:%s", typ.String())
	}

	pTyp := typ         // pointer type
	oTyp := pTyp.Elem() // object type

	fields := make([]*modelField, 0, oTyp.NumField())
	for i, l := 0, oTyp.NumField(); i < l; i++ {
		v := parseModelField(oTyp.Field(i))
		if v == nil {
			continue
		}

		v.FieldIndex = i
		fields = append(fields, v)
	}

	m.Type = pTyp
	m.Fields = fields
	m.Inited = true
	return nil
}

func (r *Reader) readHeader() error {
	m := &r.model
	if m.Headed {
		return nil
	}

	header, err := r.Reader.Read()
	if err != nil {
		return fmt.Errorf("[ldsheet] read header fail. err:%s", err)
	}

	headerMap := make(map[string]int, len(header))
	for i, v := range header {
		headerMap[v] = i
	}

	for _, field := range m.Fields {
		field.HeaderIndex = -1

		if idx := field.GetIndex(headerMap); idx >= 0 {
			field.HeaderIndex = idx

		} else if err := field.Validate(); err != nil {
			return err
		}
	}

	m.Headers = headerMap
	m.Headed = true
	return nil
}

func (r *Reader) readData(d interface{}) error {
	m := &r.model

	val := reflect.ValueOf(d)
	typ := val.Type()
	if m.Type != typ {
		return fmt.Errorf("[ldsheet] read the inconsistently type. expected:%s", m.Type.String())
	}

	obj := val.Elem()
	obj.Set(reflect.Zero(obj.Type()))

	return r.readDataValue(obj)
}

func (r *Reader) readDataValue(obj reflect.Value) error {
	m := &r.model

	line, err := r.Reader.Read()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("[ldsheet] read data fail. err:%s", err)
	}

	for _, field := range m.Fields {
		err := field.ParseValue(obj, line)
		if err != nil {
			return err
		}
	}

	return nil
}
