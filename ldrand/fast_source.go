/*
 * Copyright (C) distroy
 */

package ldrand

import (
	"math/rand"
	"sync/atomic"
	"unsafe"

	"github.com/distroy/ldgo-base/ldconv"
)

var (
	_ Rand          = (*fastSource)(nil)
	_ rand.Source   = (*fastSource)(nil)
	_ rand.Source64 = (*fastSource)(nil)
)

const (
	//                 7654321076543210
	// fastSourceStep = 0x1753715715313157
	// fastSourceStep = 0x5371795313b93157
	fastSourceStep = 0x5371795313b93d57
)

var fastSourceXor [16]uint64 = [...]uint64{
	0x3b78bb846a443726,
	0x772c3f9b061acdac,
	0x62dfa820aede2c50,
	0xa51729651df84541,
	0x93c695f2456306c5,
	0x2aa1107ecc9b60f2,
	0x19fd8406df31d838,
	0xf0ea0719370c8f9e,
	0x418e4a4878a79b6b,
	0xe49476ea842fb1e4,
	0x0c09febc20b9e3df,
	0xdf63eca7538d728d,
	0x56b2d3cd92e65e77,
	0xbe5552d1e952f903,
	0xc83bc15ff1c0a4ba,
	0x8d406d33bb751a19,
}

var fastSourcePosition [16]uint64 = [...]uint64{
	0x3158b9640fa2edc7,
	0xb364c795f1ea082d,
	0x59483721bfadec60,
	0x2405613bd98a7fec,
	0xd5f0c3b16ae97842,
	0xc7de804b29f153a6,
	0x79610fac3d28be45,
	0x326efd917405bac8,
	0x0ec92d145f3a76b8,
	0x0e98f2bd3c46a175,
	0xc1e9a670d32f548b,
	0x24db8f93ce57a061,
	0xd46e9a85013b2c7f,
	0xb568e9a20471cf3d,
	0xe3f14d70a29b65c8,
	0xacfb169238de5704,
}

// noescape hides a pointer from escape analysis.  noescape is
// the identity function but escape analysis doesn't think the
// output depends on the input. noescape is inlined and currently
// compiles down to zero instructions.
// USE CAREFULLY!
// This was copied from the runtime; see issues 23382 and 7921.
//
//go:nosplit
//go:nocheckptr
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0) // nolint
}

func NewFastSource(seed int64) rand.Source64 {
	r := &fastSource{}
	r.Seed(seed)
	// r.rand = rand.New((*fastSource)(noescape(unsafe.Pointer(r))))
	return r
}

func newFastSource(seed int64, xor [16]uint64, pos [16]uint64) *fastSource {
	src := &fastSource{
		seed: uint64(seed),
		xor:  xor,
		pos:  pos,
	}
	// src.rand = rand.New(src)
	return src
}

type fastSource struct {
	seed uint64
	xor  [16]uint64
	pos  [16]uint64
}

func (r *fastSource) Seed(seed int64) {
	n := initFastSourceXor(seed, r.xor[:], r.pos[:])
	atomic.StoreUint64(&r.seed, n)
}

func (r *fastSource) Uint64() uint64 {
	n := fastSourceNext(&r.seed, r.xor[:])
	n = fastSourceReposition(n, r.pos[:])
	return n
}

func (r *fastSource) Int63() int64   { return int64(r.Uint64() << 1 >> 1) }
func (r *fastSource) Uint32() uint32 { return uint32(r.Uint64()) }
func (r *fastSource) Int31() int32   { return int32(r.Uint32() << 1 >> 1) }
func (r *fastSource) Uint() uint     { return uint(r.Uint64()) }
func (r *fastSource) Int() int       { return int(r.Uint() << 1 >> 1) }

func (r *fastSource) Int63n(n int64) int64 { return r.Int63() % n }
func (r *fastSource) Int31n(n int32) int32 { return r.Int31() % n }
func (r *fastSource) Intn(n int) int       { return r.Int() % n }

func (r *fastSource) Float32() float32 {
again:
	f := float32(r.Float64())
	if f == 1 {
		goto again // resample; this branch is taken O(very rarely)
	}
	return f
}

func (r *fastSource) Float64() float64 {
	// A clearer, simpler implementation would be:
	//	return float64(r.Int63n(1<<53)) / (1<<53)
	// However, Go 1 shipped with
	//	return float64(r.Int63()) / (1 << 63)
	// and we want to preserve that value stream.
	//
	// There is one bug in the value stream: r.Int63() may be so close
	// to 1<<63 that the division rounds up to 1.0, and we've guaranteed
	// that the result is always less than 1.0.
	//
	// We tried to fix this by mapping 1.0 back to 0.0, but since float64
	// values near 0 are much denser than near 1, mapping 1 to 0 caused
	// a theoretically significant overshoot in the probability of returning 0.
	// Instead of that, if we round up to 1, just try again.
	// Getting 1 only happens 1/2⁵³ of the time, so most clients
	// will not observe it anyway.
again:
	f := float64(r.Int63()) / (1 << 63)
	if f == 1 {
		goto again // resample; this branch is taken O(never)
	}
	return f
}

func (r *fastSource) Read(p []byte) (int, error) {
	pos := 0
	val := uint64(0)
	mask := r.Uint64()
	for i := range len(p) {
		if pos == 0 {
			val = r.Uint64() ^ mask
			pos = 8
		}
		p[i] = byte(val)
		val >>= 8
		pos--
	}
	return len(p), nil
}

func (r *fastSource) String(n int, letters ...[]byte) string {
	if len(letters) > 0 && len(letters[0]) > 0 {
		return randString(r, n, letters[0])
	}
	return randString(r, n, randStringLetters)
}

func (r *fastSource) Bytes(n int) []byte {
	if n <= 0 {
		return nil
	}
	b := make([]byte, n)
	r.Read(b)
	return b
}

func (r *fastSource) Perm(n int) []int {
	m := make([]int, n)
	// In the following loop, the iteration when i=0 always swaps m[0] with m[0].
	// A change to remove this useless iteration is to assign 1 to i in the init
	// statement. But Perm also effects r. Making this change will affect
	// the final state of r. So this change can't be made for compatibility
	// reasons for Go 1.
	for i := range n {
		j := r.Intn(i + 1)
		m[i] = m[j]
		m[j] = i
	}
	return m
}

func (r *fastSource) Shuffle(n int, swap func(i, j int)) {
	i := n - 1
	for ; i > 1<<31-1-1; i-- {
		j := int(r.Int63n(int64(i + 1)))
		swap(i, j)
	}
	for ; i > 0; i-- {
		j := int(r.Int31n(int32(i + 1)))
		swap(i, j)
	}
}

func fastSourceNext(seed *uint64, xor []uint64) uint64 {
	n := atomic.AddUint64(seed, fastSourceStep)
	x := n
	b := n & 0xf
	// b := (n >> 60) & 0xf

	b = b ^ ((n >> 4) & 0xf)
	b = b ^ ((n >> 8) & 0xf)
	b = b ^ ((n >> 12) & 0xf)

	x = x ^ xor[b]
	x = x ^ (n << 4)
	x = x ^ ((n >> 28) & 0xffffffff)
	x = x - (x & 0xf)

	b = b ^ ((n >> 16) & 0xf)
	b = b ^ ((n >> 20) & 0xf)
	b = b ^ ((n >> 24) & 0xf)
	b = b ^ ((n >> 28) & 0xf)

	b = b ^ ((n >> 60) & 0xf)

	x = x | uint64(b)
	// // x = (x << 4 >> 4) | (b << 60)
	// x = ((x & 0xffffffff) << 32) | ((x >> 32) & 0xffffffff)
	return x
}

func fastSourceReposition(n uint64, pos []uint64) uint64 {
	b := n & 0xf
	p := pos[b]

	x := uint64(0)
	for range 16 {
		b := n & 0xf
		n = n >> 4

		i := p & 0xf
		p = p >> 4
		x |= b << (i * 4)
	}

	return (x >> 4) | (b << 60)
	// return x
}

func initFastSourceXor(seed int64, xor []uint64, pos []uint64) uint64 {
	buf := [16]byte{}
	for i := range buf {
		buf[i] = byte(i)
	}

	r := New(newFastSource(seed, fastSourceXor, fastSourcePosition))

	// generate xor
	for range 16 {
		r.Shuffle(len(buf), func(i, j int) { buf[i], buf[j] = buf[j], buf[i] })
		for j := range xor {
			xor[j] = (xor[j] << 4) | uint64(buf[j])
		}
	}

	// generate pos
	for i := range 16 {
		r.Shuffle(len(buf), func(i, j int) { buf[i], buf[j] = buf[j], buf[i] })
		// fmt.Fprintf(os.Stdout, "%#v\n", buf)

		var n uint64
		for _, v := range buf {
			if v == 0 {
				continue
			}
			n = (n << 4) | uint64(v)
		}
		pos[i] = n << 4
	}

	return r.Uint64()
}

func randString(r rand.Source64, n int, letters []byte) string {
	size := len(letters)
	if size == 0 {
		letters = randStringLetters
	}
	if size == 1 {
		b := make([]byte, n)
		for i := range b {
			b[i] = letters[0]
		}
		return ldconv.BytesToStrUnsafe(b)
	}

	maxPos := 1
	shift := 0
	switch {
	case size < (1 << 8):
		maxPos = 8
		shift = 8
	case size < (1 << 16):
		maxPos = 4
		shift = 16
	case size < (1 << 32):
		maxPos = 2
		shift = 32
	default:
		maxPos = 1
		shift = 0
	}

	mask := r.Uint64()
	b := make([]byte, n)

	pos := 0
	val := uint64(0)
	for i := range b {
		if pos == 0 {
			val = r.Uint64() ^ mask
			pos = maxPos
		}
		b[i] = letters[val%uint64(size)]
		val >>= shift
		pos--
	}

	return ldconv.BytesToStrUnsafe(b)
}
