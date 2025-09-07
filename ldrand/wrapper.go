/*
 * Copyright (C) distroy
 */

package ldrand

import "math/rand"

// A Rand is a source of random numbers.
type Rand interface {
	Int() int
	Int31() int32
	Int63() int64

	Uint() uint
	Uint32() uint32
	Uint64() uint64

	Intn(n int) int
	Int31n(n int32) int32
	Int63n(n int64) int64

	// Float32 returns, as a float32, a pseudo-random number in [0.0,1.0).
	Float32() float32
	// Float64 returns, as a float64, a pseudo-random number in [0.0,1.0).
	Float64() float64

	// Perm returns, as a slice of n ints, a pseudo-random permutation of the integers [0,n).
	Perm(n int) []int

	// Shuffle pseudo-randomizes the order of elements.
	// n is the number of elements. Shuffle panics if n < 0.
	// swap swaps the elements with indexes i and j.
	Shuffle(n int, swap func(i, j int))

	// Read generates len(p) random bytes and writes them into p. It
	// always returns len(p) and a nil error.
	Read(p []byte) (int, error)

	// Bytes generates `n` random bytes and writes them into p. It
	// always returns byte slice of length `n`.
	Bytes(n int) []byte

	String(n int, letters ...[]byte) string

	// NormFloat64 returns a normally distributed float64 in the range
	// [-math.MaxFloat64, +math.MaxFloat64] with
	// standard normal distribution (mean = 0, stddev = 1)
	// from the default Source.
	// To produce a different normal distribution, callers can
	// adjust the output using:
	//
	//  sample = NormFloat64() * desiredStdDev + desiredMean
	//
	NormFloat64() float64

	// ExpFloat64 returns an exponentially distributed float64 in the range
	// (0, +math.MaxFloat64] with an exponential distribution whose rate parameter
	// (lambda) is 1 and whose mean is 1/lambda (1) from the default Source.
	// To produce a distribution with a different rate parameter,
	// callers can adjust the output using:
	//
	//  sample = ExpFloat64() / desiredRateParameter
	//
	ExpFloat64() float64
}

func New(src rand.Source) Rand {
	switch r := src.(type) {
	case Rand:
		return r
	case *rand.Rand:
		return randWrapper{Rand: r}
	}
	return randWrapper{Rand: rand.New(src)}
}

type randWrapper struct {
	*rand.Rand
}

func (r randWrapper) Uint() uint { return uint(r.Rand.Uint64()) }
func (r randWrapper) Bytes(n int) []byte {
	if n <= 0 {
		return nil
	}
	b := make([]byte, n)
	r.Read(b)
	return b
}
func (r randWrapper) String(n int, letters ...[]byte) string {
	if len(letters) > 0 && len(letters[0]) > 0 {
		return randString(r, n, letters[0])
	}
	return randString(r, n, randStringLetters)
}
