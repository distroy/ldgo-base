/*
 * Copyright (C) distroy
 */

package ldbyte

import (
	"io"

	"github.com/distroy/ldgo-base/internal/_buf"
)

var (
	_ io.ReadWriteCloser = (*RingBuffer)(nil)
)

func NewRingBuffer(n int) *RingBuffer {
	if n <= 0 {
		n = 1024
	}
	buf := make([]byte, n)
	return &RingBuffer{
		buf: _buf.NewRing(buf),
	}
}

type RingBuffer struct {
	buf *_buf.Ring[byte]
}

func (b RingBuffer) Close() error { return b.buf.Close() }
func (b RingBuffer) Closed() bool { return b.buf.Closed() }

func (b RingBuffer) Cap() int  { return b.buf.Cap() }
func (b RingBuffer) Size() int { return b.buf.Size() }

func (b RingBuffer) Write(d []byte) (int, error) { return b.buf.Write(d) }
func (b RingBuffer) Read(d []byte) (int, error)  { return b.buf.Read(d) }

func NewBlockingRingBuffer(n int) BlockingRingBuffer {
	if n <= 0 {
		n = 1024
	}
	b := BlockingRingBuffer{
		buf: _buf.NewBlockingRing(make([]byte, n)),
	}
	return b
}

type BlockingRingBuffer struct {
	buf *_buf.BlockingRing[byte]
}

func (b BlockingRingBuffer) Close() error { return b.buf.Close() }
func (b BlockingRingBuffer) Closed() bool { return b.buf.Closed() }

func (b BlockingRingBuffer) Cap() int  { return b.buf.Cap() }
func (b BlockingRingBuffer) Size() int { return b.buf.Size() }

func (b BlockingRingBuffer) Write(d []byte) (int, error) { return b.buf.Write(d) }
func (b BlockingRingBuffer) Read(d []byte) (int, error)  { return b.buf.Read(d) }
