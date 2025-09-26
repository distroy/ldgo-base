/*
 * Copyright (C) distroy
 */

package ldtrie

import (
	"strings"

	"github.com/distroy/ldgo-base/ldconv"
	"github.com/distroy/ldgo-base/lditer"
)

/*
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkContainsIgnoreCase-16             23204             51972 ns/op
BenchmarkRuneTrieIgnoreCase-16          17215753                61.98 ns/op
BenchmarkByteTrieIgnoreCase-16          18361742                63.27 ns/op
*/

type Trie interface {
	Insert(words ...string) int

	Search(text string) (string, bool)
	SearchAll(text string) ([]string, bool)
	SearchBest(text string) (string, bool)

	SearchPrefix(text string) (string, bool)
	SearchPrefixAll(text string) ([]string, bool)
	SearchPrefixBest(text string) (string, bool)

	SearchExact(text string) bool
}

func NewTrie(opts ...Option) Trie {
	option := &trieOpt{}
	for _, opt := range opts {
		opt(option)
	}
	return newRuneTrie(option)
	// return newByteTrie(option)
}

func newRuneTrie(option *trieOpt) *trie[runeIface, rune] { return newTrie[runeIface](option) }
func newByteTrie(option *trieOpt) *trie[byteIface, byte] { return newTrie[byteIface](option) }

func newTrie[Iface itemIface[T], T byte | rune](option *trieOpt) *trie[Iface, T] {
	return &trie[Iface, T]{
		Option: *option,
		Root:   trieNode[T]{
			// Children: make(map[T]*trieNode[T]),
		},
	}
}

type itemIface[T byte | rune] interface {
	Conv(s string) []T
	Iter(s string) lditer.Seq2[int, T]
}

type byteIface struct{}

func (i byteIface) Conv(s string) []byte                 { return ldconv.StrToBytesUnsafe(s) }
func (i byteIface) Iter(s string) lditer.Seq2[int, byte] { return lditer.Slice(i.Conv(s)) }

type runeIface struct{}

func (i runeIface) Conv(s string) []rune { return []rune(s) }
func (i runeIface) Iter(s string) lditer.Seq2[int, rune] {
	return func(yield func(int, rune) bool) {
		for i, v := range s {
			if !yield(i, v) {
				break
			}
		}
	}
}

type trieNode[T byte | rune] struct {
	Children map[T]*trieNode[T]
	Origin   string
	Item     T
	IsEnd    bool
}

type trie[Iface itemIface[T], T byte | rune] struct {
	Root   trieNode[T]
	Option trieOpt
}

func (t *trie[Iface, T]) Insert(words ...string) int {
	count := 0
	for _, word := range words {
		ok := t.insert(word)
		if ok {
			count++
		}
	}
	return count
}

func (t *trie[Iface, T]) iface() (x Iface) { return }

func (t *trie[Iface, T]) insert(word string) bool {
	if word == "" && !t.Option.AllowEmpty {
		return false
	}

	origin := word
	if t.Option.IgnoreCase {
		word = strings.ToLower(word)
	}

	curr := &t.Root
	for _, c := range t.iface().Iter(word) {
		next := curr.Children[c]
		if next == nil {
			next = &trieNode[T]{
				Item: c,
			}
			if curr.Children == nil {
				curr.Children = make(map[T]*trieNode[T])
			}
			curr.Children[c] = next
		}
		curr = next
		if curr.IsEnd && t.Option.DisableBest {
			return false
		}
	}

	curr.IsEnd = true
	curr.Origin = origin
	if t.Option.DisableBest {
		curr.Children = nil
	}
	return true
}

func (t *trie[Iface, T]) SearchExact(text string) bool {
	if t.Option.IgnoreCase {
		text = strings.ToLower(text)
	}
	curr := &t.Root
	for _, c := range t.iface().Iter(text) {
		next := curr.Children[c]
		if next == nil {
			return false
		}

		curr = next
	}
	return curr.IsEnd
}

func (t *trie[Iface, T]) Search(text string) (string, bool) {
	if t.Option.IgnoreCase {
		text = strings.ToLower(text)
	}
	node := t.search(text)
	if node != nil {
		return node.Origin, true
	}
	return "", false
}

func (t *trie[Iface, T]) SearchAll(text string) ([]string, bool) {
	if t.Option.IgnoreCase {
		text = strings.ToLower(text)
	}
	res, cnt := t.searchAll(text)
	return res, cnt > 0
}

func (t *trie[Iface, T]) SearchBest(text string) (string, bool) {
	if t.Option.IgnoreCase {
		text = strings.ToLower(text)
	}
	node := t.searchBest(text)
	if node != nil {
		return node.Origin, true
	}
	return "", false
}

func (t *trie[Iface, T]) SearchPrefix(text string) (string, bool) {
	if t.Option.IgnoreCase {
		text = strings.ToLower(text)
	}
	items := t.iface().Conv(text)
	length := len(items)
	node := t.searchItems(items, 0, length)
	if node != nil {
		return node.Origin, true
	}
	return "", false
}

func (t *trie[Iface, T]) SearchPrefixAll(text string) ([]string, bool) {
	if t.Option.IgnoreCase {
		text = strings.ToLower(text)
	}
	items := t.iface().Conv(text)
	length := len(items)
	res, cnt := t.searchAllItems(nil, items, 0, length)
	return res, cnt > 0
}

func (t *trie[Iface, T]) SearchPrefixBest(text string) (string, bool) {
	if t.Option.IgnoreCase {
		text = strings.ToLower(text)
	}
	items := t.iface().Conv(text)
	length := len(items)
	node := t.searchBestItems(items, 0, length)
	if node != nil {
		return node.Origin, true
	}
	return "", false
}

func (t *trie[Iface, T]) search(text string) *trieNode[T] {
	items := t.iface().Conv(text)
	length := len(items)
	for i := range length {
		node := t.searchItems(items, i, length)
		if node != nil {
			return node
		}
	}
	return nil
}

func (t *trie[Iface, T]) searchItems(items []T, start, length int) *trieNode[T] {
	curr := &t.Root
	for i := start; i < length; i++ {
		c := items[i]
		next := curr.Children[c]
		if next == nil {
			return nil
		}

		curr = next
		if curr.IsEnd {
			return curr
		}
	}
	return nil
}

func (t *trie[Iface, T]) searchAll(text string) ([]string, int) {
	res := make([]string, 0, 16)
	items := t.iface().Conv(text)
	length := len(items)
	count := 0
	for i := range length {
		tmpRes, tmpCnt := t.searchAllItems(res, items, i, length)
		res = tmpRes
		count += tmpCnt
	}
	return res, count
}

func (t *trie[Iface, T]) searchAllItems(res []string, items []T, start, length int) ([]string, int) {
	if res == nil {
		res = make([]string, 0, 16)
	}
	curr := &t.Root
	count := 0
	for i := start; i < length; i++ {
		c := items[i]
		next := curr.Children[c]
		if next == nil {
			return res, count
		}

		curr = next
		if curr.IsEnd {
			res = append(res, curr.Origin)
			count++
		}
	}
	return res, count
}

func (t *trie[Iface, T]) searchBest(text string) *trieNode[T] {
	items := t.iface().Conv(text)
	length := len(items)
	res := (*trieNode[T])(nil)
	for i := range length {
		node := t.searchBestItems(items, i, length)
		if res == nil || node != nil && len(node.Origin) > len(res.Origin) {
			res = node
			// return node
		}
	}
	return res
}

func (t *trie[Iface, T]) searchBestItems(items []T, start, length int) *trieNode[T] {
	curr := &t.Root
	res := (*trieNode[T])(nil)
	for i := start; i < length; i++ {
		c := items[i]
		next := curr.Children[c]
		if next == nil {
			break
		}
		curr = next
		if curr.IsEnd {
			res = curr
			// return curr
		}
	}
	return res
}
