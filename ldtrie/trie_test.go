/*
 * Copyright (C) distroy
 */

package ldtrie

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/distroy/ldgo-base/ldsort"
)

type (
	testWantFunc[R any]   = func(text string, blacklist []string) (R, bool)
	testNewGotFunc[R any] = func(blacklist []string) func(text string) (R, bool)
)

func testSearchIgnoreCase(text string, blacklist []string) (string, bool) {
	text = strToLower(text)
	return testSearch(text, blacklist)
}
func testSearch(text string, blacklist []string) (string, bool) {
	res := ""
	idx := -1
	for _, sub := range blacklist {
		if pos := strings.Index(text, sub); pos >= 0 && (res == "" || pos < idx) {
			res = sub
			idx = pos
		}
	}
	return res, idx >= 0
}

func testSearchAllIgnoreCase(text string, blacklist []string) ([]string, bool) {
	text = strToLower(text)
	return testSearchAll(text, blacklist)
}
func testSearchAll(text string, blacklist []string) ([]string, bool) {
	res := make([]string, 0, 16)
	isMapping := false
	for _, sub := range blacklist {
		if strings.Contains(text, sub) {
			isMapping = true
			res = append(res, sub)
		}
	}
	return res, isMapping
}

func testSearchBest(text string, blacklist []string) (string, bool) {
	res := ""
	idx := -1
	for _, sub := range blacklist {
		if len(sub) < len(res) {
			continue
		}

		pos := strings.Index(text, sub)
		if pos < 0 {
			continue

		} else if idx < 0 || pos <= idx || len(sub) > len(res) {
			res = sub
			idx = pos
		}
	}
	return res, idx >= 0
}

func testSearchExact(text string, blacklist []string) (int, bool) {
	ok := slices.Contains(blacklist, text)
	return 0, ok
}

func testSearchPrefixIgnoreCase(text string, blacklist []string) (string, bool) {
	text = strToLower(text)
	return testSearchPrefix(text, blacklist)
}
func testSearchPrefix(text string, blacklist []string) (string, bool) {
	res := ""
	ok := false
	for _, sub := range blacklist {
		if !strings.HasPrefix(text, sub) {
			continue
		}
		if !ok || len(sub) < len(res) {
			res = sub
			ok = true
		}
	}
	return res, ok
}
func testSearchPrefixBest(text string, blacklist []string) (string, bool) {
	res := ""
	ok := false
	for _, sub := range blacklist {
		if !strings.HasPrefix(text, sub) {
			continue
		}
		if !ok || len(sub) > len(res) {
			res = sub
			ok = true
		}
	}
	return res, ok
}
func testSearchPrefixAll(text string, blacklist []string) ([]string, bool) {
	res := make([]string, 0, 16)
	for _, sub := range blacklist {
		if !strings.HasPrefix(text, sub) {
			continue
		}
		res = append(res, sub)
	}
	return res, len(res) > 0
}

func testCommon[R any](t testing.TB, ignoreCase bool, fnWant testWantFunc[R], fnNewGot testNewGotFunc[R]) {
	blacklist := testGetBlacklist(t, ignoreCase)
	testcases := getTestcases(t)

	fnGot := fnNewGot(blacklist)

	trueCount := 0
	falseCount := 0
	for _, text := range testcases {
		res0, ok0 := fnWant(text, blacklist)
		res1, ok1 := fnGot(text)
		if ok0 != ok1 || !reflect.DeepEqual(res0, res1) {
			t.Errorf("%s fail. text:%s, want:[%v:%v], got:[%v:%v]", t.Name(), text, ok0, res0, ok1, res1)
			continue
		}

		if ok0 {
			trueCount++
		} else {
			falseCount++
		}
	}

	t.Logf("%s result. true:%d, false:%d", t.Name(), trueCount, falseCount)
}

func testTrieSearchIgnoreCase(t testing.TB, fnNew func(...Option) Trie) {
	testCommon(t, true, testSearchIgnoreCase, func(blacklist []string) func(text string) (string, bool) {
		tt := fnNew(IgnoreCase(true))
		tt.Insert(blacklist...)
		return tt.Search
	})
}

func testTrieSearchAllIgnoreCase(t testing.TB, fnNew func(...Option) Trie) {
	fnWant := func(text string, blacklist []string) ([]string, bool) {
		res, ok := testSearchAllIgnoreCase(text, blacklist)
		ldsort.SortStrings(res)
		return ldsort.UniqStrings(res), ok
	}
	testCommon(t, true, fnWant, func(blacklist []string) func(text string) ([]string, bool) {
		tt := fnNew(IgnoreCase(true))
		tt.Insert(blacklist...)
		return func(text string) ([]string, bool) {
			res, ok := tt.SearchAll(text)
			ldsort.SortStrings(res)
			return ldsort.UniqStrings(res), ok
		}
	})
}

func testTrieSearchBest(t testing.TB, fnNew func(...Option) Trie) {
	testCommon(t, false, testSearchBest, func(blacklist []string) func(text string) (string, bool) {
		tt := fnNew()
		tt.Insert(blacklist...)
		tt.Insert("")
		return tt.SearchBest
	})
}

func testTrieSearchPrefixDisableBest(t testing.TB, fnNew func(...Option) Trie) {
	testCommon(t, true, testSearchPrefixIgnoreCase, func(blacklist []string) func(text string) (string, bool) {
		tt := fnNew(DisableBest(true))
		tt.Insert(blacklist...)
		tt.Insert("")
		return tt.SearchPrefix
	})
}

func testTrieSearchPrefixBest(t testing.TB, fnNew func(...Option) Trie) {
	testCommon(t, false, testSearchPrefixBest, func(blacklist []string) func(text string) (string, bool) {
		tt := fnNew(DisableBest(false))
		tt.Insert(blacklist...)
		tt.Insert("")
		return tt.SearchPrefixBest
	})
}

func testTrieSearchPrefixAll(t testing.TB, fnNew func(...Option) Trie) {
	fnWant := func(text string, blacklist []string) ([]string, bool) {
		res, ok := testSearchPrefixAll(text, blacklist)
		ldsort.SortStrings(res)
		return ldsort.UniqStrings(res), ok
	}
	testCommon(t, false, fnWant, func(blacklist []string) func(text string) ([]string, bool) {
		tt := fnNew(DisableBest(false))
		tt.Insert(blacklist...)
		tt.Insert("")
		return func(text string) ([]string, bool) {
			res, ok := tt.SearchPrefixAll(text)
			ldsort.SortStrings(res)
			return ldsort.UniqStrings(res), ok
		}
	})
}

func testTrieSearchExactAllowEmpty(t testing.TB, fnNew func(...Option) Trie) {
	testCommon(t, false, testSearchExact, func(blacklist []string) func(text string) (int, bool) {
		tt := fnNew(AllowEmpty(true))
		tt.Insert(blacklist...)
		tt.Insert("")
		return func(text string) (int, bool) {
			ok := tt.SearchExact(text)
			return 0, ok
		}
	})
}

func TestByteTrieSearchIgnoreCase(t *testing.T) {
	fnNew := func(opts ...Option) Trie {
		return newByteTrie(getOpt(opts...))
	}
	testTrieSearchIgnoreCase(t, fnNew)
}

func TestRuneTrieSearchIgnoreCase(t *testing.T) {
	fnNew := func(opts ...Option) Trie {
		return newRuneTrie(getOpt(opts...))
	}
	testTrieSearchIgnoreCase(t, fnNew)
}

func TestByteTrieSearchAllIgnoreCase(t *testing.T) {
	fnNew := func(opts ...Option) Trie {
		return newByteTrie(getOpt(opts...))
	}
	testTrieSearchAllIgnoreCase(t, fnNew)
}
func TestRuneTrieSearchAllIgnoreCase(t *testing.T) {
	fnNew := func(opts ...Option) Trie {
		return newRuneTrie(getOpt(opts...))
	}
	testTrieSearchAllIgnoreCase(t, fnNew)
}

func TestByteTrieSearchBest(t *testing.T) {
	fnNew := func(opts ...Option) Trie {
		return newByteTrie(getOpt(opts...))
	}
	testTrieSearchBest(t, fnNew)
}
func TestRuneTrieSearchBest(t *testing.T) {
	fnNew := func(opts ...Option) Trie {
		return newRuneTrie(getOpt(opts...))
	}
	testTrieSearchBest(t, fnNew)
}

func TestByteTrieSearchPrefixDisableBest(t *testing.T) {
	fnNew := func(opts ...Option) Trie {
		return newByteTrie(getOpt(opts...))
	}
	testTrieSearchPrefixDisableBest(t, fnNew)
}
func TestRuneTrieSearchPrefixDisableBest(t *testing.T) {
	fnNew := func(opts ...Option) Trie {
		return newRuneTrie(getOpt(opts...))
	}
	testTrieSearchPrefixDisableBest(t, fnNew)
}

func TestByteTrieSearchPrefixBest(t *testing.T) {
	fnNew := func(opts ...Option) Trie {
		return newByteTrie(getOpt(opts...))
	}
	testTrieSearchPrefixBest(t, fnNew)
}
func TestRuneTrieSearchPrefixBest(t *testing.T) {
	fnNew := func(opts ...Option) Trie {
		return newRuneTrie(getOpt(opts...))
	}
	testTrieSearchPrefixBest(t, fnNew)
}

func TestByteTrieSearchPrefixAll(t *testing.T) {
	fnNew := func(opts ...Option) Trie {
		return newByteTrie(getOpt(opts...))
	}
	testTrieSearchPrefixAll(t, fnNew)
}
func TestRuneTrieSearchPrefixAll(t *testing.T) {
	fnNew := func(opts ...Option) Trie {
		return newRuneTrie(getOpt(opts...))
	}
	testTrieSearchPrefixAll(t, fnNew)
}

func TestByteTrieSearchExactAllowEmpty(t *testing.T) {
	fnNew := func(opts ...Option) Trie {
		return newByteTrie(getOpt(opts...))
	}
	testTrieSearchExactAllowEmpty(t, fnNew)
}
func TestRuneTrieSearchExactAllowEmpty(t *testing.T) {
	fnNew := func(opts ...Option) Trie {
		return newRuneTrie(getOpt(opts...))
	}
	testTrieSearchExactAllowEmpty(t, fnNew)
}
