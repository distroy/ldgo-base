/*
 * Copyright (C) distroy
 */

package lderr

import (
	"context"
	"io"
	"testing"
)

type testUnwrapError struct {
	error
}

func (e testUnwrapError) Unwrap() error { return e.error }

func TestIs(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{
			name:   "nil & nil",
			err:    nil,
			target: nil,
			want:   true,
		},
		{
			name:   "nil & succ",
			err:    nil,
			target: ErrSuccess,
			want:   true,
		},
		{
			name:   "succ & nil",
			err:    ErrSuccess,
			target: nil,
			want:   true,
		},
		{
			name:   "succ & succ",
			err:    ErrSuccess,
			target: ErrSuccess,
			want:   true,
		},
		{
			name:   "nil & unknown",
			err:    nil,
			target: ErrUnkown,
			want:   false,
		},
		{
			name:   "succ & unknown",
			err:    ErrSuccess,
			target: ErrUnkown,
			want:   false,
		},
		{
			name:   "context canceled",
			err:    ErrCtxCanceled,
			target: context.Canceled,
			want:   true,
		},
		{
			name:   "context deadline exceeded",
			err:    ErrCtxDeadlineExceeded,
			target: context.DeadlineExceeded,
			want:   true,
		},
		{
			name:   "iof & wrap iof",
			err:    io.EOF,
			target: Wrap(io.EOF),
			want:   false,
		},
		{
			name:   "wrap iof & iof",
			err:    Wrap(io.EOF),
			target: io.EOF,
			want:   true,
		},
		{
			name:   "wrap iof & unknown",
			err:    Wrap(io.EOF),
			target: io.EOF,
			want:   true,
		},
		{
			name:   "wrap iof & succ",
			err:    Wrap(io.EOF),
			target: ErrSuccess,
			want:   false,
		},
		{
			name:   "unknown & overwrite unknown",
			err:    ErrUnkown,
			target: Override(ErrUnkown, "abc"),
			want:   true,
		},
		{
			name:   "unknown & string error",
			err:    ErrUnkown,
			target: strError(errMessageUnkonw),
			want:   true,
		},
		{
			name:   "unknown & panic",
			err:    testUnwrapError{ErrUnkown},
			target: ErrServicePanic,
			want:   false,
		},
		{
			name:   "unknown & nil",
			err:    testUnwrapError{ErrUnkown},
			target: nil,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.err, tt.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testDetailsErrorEqual(a, b error) bool {
	if a == b {
		return true
	}
	if aa, bb := GetCode(a), GetCode(b); aa != bb {
		return false
	}
	if aa, bb := GetStatus(a), GetStatus(b); aa != bb {
		return false
	}
	if aa, bb := GetMessage(a), GetMessage(b); aa != bb {
		return false
	}
	aa, bb := GetDetails(a), GetDetails(b)
	if len(aa) != len(bb) {
		return false
	}
	for i := range aa {
		if aa[i] != bb[i] {
			return false
		}
	}
	return true
}

func TestWithDetail(t *testing.T) {
	equal := testDetailsErrorEqual

	type args struct {
		err     error
		details []string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "nil",
			args: args{
				err:     nil,
				details: []string{},
			},
			want: nil,
		},
		{
			name: "nil with details",
			args: args{
				err: nil,
				details: []string{
					"detail-0",
					"detail-1",
				},
			},
			want: &detailsError{
				error: commError{
					error:  ErrSuccess,
					status: ErrSuccess.Status(),
					code:   ErrSuccess.Code(),
				},
				details: []string{
					"detail-0",
					"detail-1",
				},
			},
		},
		{
			name: "unknown",
			args: args{
				err:     ErrUnkown,
				details: []string{},
			},
			want: ErrUnkown,
		},
		{
			name: "unknown with details",
			args: args{
				err: ErrUnkown,
				details: []string{
					"detail-0",
					"detail-1",
				},
			},
			want: &detailsError{
				error: commError{
					error:  ErrUnkown,
					status: ErrUnkown.Status(),
					code:   ErrUnkown.Code(),
				},
				details: []string{
					"detail-0",
					"detail-1",
				},
			},
		},
		{
			name: "str error with details",
			args: args{
				err: strError("str error"),
				details: []string{
					"detail-0",
					"detail-1",
				},
			},
			want: &detailsError{
				error: commError{
					error:  strError("str error"),
					status: ErrUnkown.Status(),
					code:   ErrUnkown.Code(),
				},
				details: []string{
					"detail-0",
					"detail-1",
				},
			},
		},
		{
			name: "details error with details",
			args: args{
				err: &detailsError{
					error: commError{
						error:  strError("details error"),
						status: ErrUnkown.Status(),
						code:   ErrUnkown.Code(),
					},
					details: []string{
						"detail-0",
					},
				},
				details: []string{
					"detail-1",
				},
			},
			want: &detailsError{
				error: commError{
					error:  strError("details error"),
					status: ErrUnkown.Status(),
					code:   ErrUnkown.Code(),
				},
				details: []string{
					"detail-0",
					"detail-1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDetail(tt.args.err, tt.args.details...); !equal(got, tt.want) {
				t.Errorf("WithDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}
