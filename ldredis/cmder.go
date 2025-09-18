/*
 * Copyright (C) distroy
 */

package ldredis

type Cmder interface {
	// command name.
	// e.g. "set k v ex 10" -> "set", "cluster info" -> "cluster".
	Name() string

	// full command name.
	// e.g. "set k v ex 10" -> "set", "cluster info" -> "cluster info".
	FullName() string

	// all args of the command.
	// e.g. "set k v ex 10" -> "[set k v ex 10]".
	Args() []any

	// format request and response string.
	// e.g. "set k v ex 10" -> "set k v ex 10: OK", "get k" -> "get k: v".
	String() string

	SetErr(error)
	Err() error

	// SetFirstKeyPos(int8)
}

type valCmder[V any] interface {
	Cmder

	Result() (V, error)
	SetVal(val V)
	Val() V
}

type BoolCmder interface {
	valCmder[bool]

	// Result() (bool, error)
	// SetVal(val bool)
	// Val() bool
}

type IntCmder interface {
	valCmder[int64]

	// Result() (int64, error)
	// SetVal(val int64)
	// Val() int64

	// Uint64() (uint64, error)
}

type StringCmder interface {
	valCmder[string]

	// Result() (string, error)
	// SetVal(val string)
	// Val() string

	// Bool() (bool, error)
	// Bytes() ([]byte, error)
	// Float32() (float32, error)
	// Float64() (float64, error)
	// Int() (int, error)
	// Int64() (int64, error)
	// Scan(val any) error
	// Time() (time.Time, error)
	// Uint64() (uint64, error)
}
