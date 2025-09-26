/*
 * Copyright (C) distroy
 */

package ldtrie

type Option func(opt *trieOpt)

type trieOpt struct {
	IgnoreCase  bool
	AllowEmpty  bool
	DisableBest bool
}

// default: false
func IgnoreCase(b bool) Option {
	return func(option *trieOpt) {
		option.IgnoreCase = b
	}
}

// default: false
func AllowEmpty(b bool) Option {
	return func(option *trieOpt) {
		option.AllowEmpty = b
	}
}

// default: false
func DisableBest(b bool) Option {
	return func(option *trieOpt) {
		option.DisableBest = b
	}
}

func getOpt(opts ...Option) trieOpt {
	option := trieOpt{}
	for _, opt := range opts {
		opt(&option)
	}
	return option
}
