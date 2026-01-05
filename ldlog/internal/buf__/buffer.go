/*
 * Copyright (C) distroy
 */

package buf__

import (
	"encoding/json"
	"slices"
	"strconv"
	"time"

	"github.com/distroy/ldgo-base/ldsync"
)

const (
	TimeLayout = "2006-01-02T15:04:05.000Z0700"
)

type Buffer []byte

// Having an initial size gives a dramatic speedup.
var bufPool = ldsync.GetPool(func() []byte { return make([]byte, 0, 1024) })

func NewBuffer() *Buffer {
	buf := bufPool.Get()
	buf = buf[:0]
	return (*Buffer)(&buf)
}

func (b *Buffer) Free() {
	// To reduce peak allocation, return only smaller buffers to the pool.
	const maxBufferSize = 16 << 10
	b.Reset()
	if b.Cap() <= maxBufferSize {
		var bb []byte = *b
		bufPool.Put(bb)
	}
}

func (b *Buffer) Grow(n int) {
	if n <= 0 {
		return
	}
	*b = slices.Grow(*b, n)
}

// TrimNewline trims any final "\n" byte from the end of the buffer.
func (b *Buffer) TrimNewline() {
	// *b = bytes.TrimSuffix(*b, []byte{'\n'})
	if i := len(*b) - 1; i >= 0 {
		if (*b)[i] == '\n' {
			*b = (*b)[:i]
		}
	}
}

func (b *Buffer) Reset() { b.SetLen(0) }

func (b *Buffer) Write(p []byte) (int, error) {
	*b = append(*b, p...)
	return len(p), nil
}

func (b *Buffer) WriteString(s string) (int, error) {
	*b = append(*b, s...)
	return len(s), nil
}

func (b *Buffer) WriteByte(c byte) error {
	*b = append(*b, c)
	return nil
}

func (b *Buffer) WriteBool(v bool) (int, error)    { return b.write(func() { b.AppendBool(v) }) }
func (b *Buffer) WriteQuote(s string) (int, error) { return b.write(func() { b.AppendQuote(s) }) }

func (b *Buffer) WriteInt(v int64) (int, error)   { return b.write(func() { b.AppendInt(v) }) }
func (b *Buffer) WriteUint(v uint64) (int, error) { return b.write(func() { b.AppendUint(v) }) }
func (b *Buffer) WriteFloat(v float64, bitSize int) (int, error) {
	return b.write(func() { b.AppendFloat(v, bitSize) })
}
func (b *Buffer) WriteComplex(v complex128, bitSize int) (int, error) {
	return b.write(func() { b.AppendComplex(v, bitSize) })
}

func (b *Buffer) WriteDuration(d time.Duration) (int, error) {
	return b.write(func() { b.AppendDuration(d) })
}
func (b *Buffer) WriteTime(t time.Time, layout string) (int, error) {
	return b.write(func() { b.AppendTime(t, layout) })
}
func (b *Buffer) WriteJson(v any) (int, error) {
	l0 := b.Len()
	err := b.AppendJson(v)
	return b.Len() - l0, err
}

func (b *Buffer) write(f func()) (int, error) {
	l0 := b.Len()
	f()
	return b.Len() - l0, nil
}

func (b *Buffer) Bytes() []byte  { return *b }
func (b *Buffer) String() string { return string(*b) }

func (b *Buffer) Len() int     { return len(*b) }
func (b *Buffer) Cap() int     { return cap(*b) }
func (b *Buffer) SetLen(n int) { *b = (*b)[:n] }

func (b *Buffer) AppendBool(v bool)     { *b = strconv.AppendBool(*b, v) }
func (b *Buffer) AppendByte(v byte)     { *b = append(*b, v) }
func (b *Buffer) AppendBytes(v []byte)  { *b = append(*b, v...) }
func (b *Buffer) AppendString(v string) { *b = append(*b, v...) }
func (b *Buffer) AppendQuote(v string)  { *b = strconv.AppendQuote(*b, v) }

func (b *Buffer) AppendInt(v int64)   { *b = strconv.AppendInt(*b, v, 10) }
func (b *Buffer) AppendUint(v uint64) { *b = strconv.AppendUint(*b, v, 10) }
func (b *Buffer) AppendFloat(v float64, bitSize int) {
	*b = strconv.AppendFloat(*b, v, 'f', -1, bitSize)
}
func (b *Buffer) AppendComplex(v complex128, bitSize int) {
	s := strconv.FormatComplex(v, 'f', -1, bitSize)
	*b = append(*b, s...)
}

func (b *Buffer) AppendDuration(d time.Duration) { *b = append(*b, d.String()...) }
func (b *Buffer) AppendTime(t time.Time, layout string) {
	if layout == "" {
		layout = TimeLayout
	}
	*b = t.AppendFormat(*b, layout)
}

func (b *Buffer) AppendJson(v any) error {
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return err
	}
	b.TrimNewline()
	return nil
}
