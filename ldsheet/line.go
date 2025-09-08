/*
 * Copyright (C) distroy
 */

package ldsheet

import "io"

type LineReader interface {
	Read() ([]string, error)
}

func Lines(content [][]string) LineReader {
	return &lines{lines: content}
}

type lines struct {
	lines [][]string
	index int
}

func (p *lines) Read() ([]string, error) {
	if p.index >= len(p.lines) {
		return nil, io.EOF
	}

	res := p.lines[p.index]
	p.index++
	return res, nil
}
