/*
 * Copyright (C) distroy
 */

package lderr

import (
	"errors"
	"maps"
)

type (
	isIface = interface{ Is(error) bool }

	unwrapIface  = interface{ Unwrap() error }
	unwrapsIface = interface{ Unwrap() []error }

	codeIface   = interface{ Code() int }
	statusIface = interface{ Status() int }

	detailsIface = interface{ Details() []string }
	detailIface  = interface{ Detail() string }

	extraIface = interface{ Extra() map[string]string }
)

var (
	_ Error       = (*commError)(nil)
	_ isIface     = (*commError)(nil)
	_ codeIface   = (*commError)(nil)
	_ statusIface = (*commError)(nil)
	_ unwrapIface = (*commError)(nil)

	_ error        = (*detailsError)(nil)
	_ detailsIface = (*detailsError)(nil)
	_ unwrapIface  = (*detailsError)(nil)

	_ error       = (*extraError)(nil)
	_ extraIface  = (*extraError)(nil)
	_ unwrapIface = (*extraError)(nil)
)

type Error interface {
	error

	Status() int
	Code() int
}

// Is reports whether any error in err's tree matches target.
//
// The tree consists of err itself, followed by the errors obtained by repeatedly
// calling Unwrap. When err wraps multiple errors, Is examines err followed by a
// depth-first traversal of its children.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
//
// An error type might provide an Is method so it can be treated as equivalent
// to an existing error. For example, if MyError defines
//
//	func (m MyError) Is(target error) bool { return target == fs.ErrExist }
//
// then Is(MyError{}, fs.ErrExist) returns true. See syscall.Errno.Is for
// an example in the standard library. An Is method should only shallowly
// compare err and the target and not call Unwrap on either.
func Is(err, target error) bool {
	if err == nil {
		return IsSuccess(target)
	}
	if target == nil {
		return IsSuccess(err)
	}
	return errors.Is(err, target)
}

func In(err error, targets ...error) bool {
	for _, target := range targets {
		if Is(err, target) {
			return true
		}
	}
	return false
}

func IsSuccess(err error) bool {
	if err == nil {
		return true
	}
	if GetCode(err) == 0 {
		return true
	}
	return false
}

func New(status, code int, message string) error {
	return commError{
		error:  strError(message),
		status: status,
		code:   code,
	}
}

func Wrap(err error, def ...Error) error {
	if err == nil {
		return nil
	}

	if v, ok := err.(Error); ok {
		return v
	}

	if e := getMatchError(err); e != nil {
		return e
	}

	code, codeOk := getCode(err)
	status, statusOk := getStatus(err)

	if codeOk && statusOk {
		return err
	}

	d := ErrUnkown
	if len(def) != 0 {
		d = def[0]
	}

	if !codeOk {
		code = d.Code()
	}
	if !statusOk {
		status = d.Status()
	}

	return commError{
		error:  err,
		status: status,
		code:   code,
	}
}

func Override(err error, message string) error {
	e := commError{
		error:  strError(message),
		status: GetStatus(err),
		code:   GetCode(err),
	}
	return newWithDetails(e, GetDetails(err))
}

func newWithDetails(err commError, details []string) error {
	if len(details) == 0 {
		return err
	}
	return &detailsError{
		error:   err,
		details: details,
	}
}

type commError struct {
	error

	status int
	code   int
}

func (e commError) Status() int   { return e.status }
func (e commError) Code() int     { return e.code }
func (e commError) Unwrap() error { return e.error }
func (e commError) Is(target error) bool {
	if err, _ := target.(codeIface); err != nil && e.Code() == err.Code() {
		return true
	}
	return false
}

type strError string

func (e strError) Error() string { return string(e) }

func WithDetail(err error, details ...string) error {
	return WithDetails(err, details)
}

func WithDetails(err error, details []string) error {
	if len(details) == 0 {
		return err
	}
	t := GetDetails(err)

	if len(details)+len(t) == 0 {
		return Wrap(err)
	}
	if e := getMatchError(err); e != nil {
		err = e
	}

	if len(t) == 0 {
		return &detailsError{
			error:   err,
			details: details,
		}
	}

	d := make([]string, 0, len(details)+len(t))
	d = append(d, t...)
	d = append(d, details...)

	return &detailsError{
		error:   err,
		details: d,
	}
}

type detailsError struct {
	error

	details []string
}

func (e *detailsError) Details() []string { return e.details }
func (e *detailsError) Unwrap() error     { return e.error }

func WithExtra(err error, extra map[string]string) error {
	if len(extra) == 0 {
		return err
	}
	if err == nil {
		err = ErrSuccess
	}
	m := GetExtra(err)
	if len(m) != 0 {
		x := make(map[string]string, len(extra)+len(m))
		maps.Copy(x, m)
		maps.Copy(x, extra)
		extra = m
	}

	return &extraError{
		error: err,
		extra: extra,
	}
}

type extraError struct {
	error

	extra map[string]string
}

func (e *extraError) Extra() map[string]string { return e.extra }
func (e *extraError) Unwrap() error            { return e.error }
