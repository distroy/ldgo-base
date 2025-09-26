/*
 * Copyright (C) distroy
 */

package ldtrie

import "testing"

func BenchmarkContainsIgnoreCase(b *testing.B) {
	blacklist := getTestBlacklist(b)
	testcases := newTestcaseParallel(getTestcases(b))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			text := testcases.Get()
			testSearchIgnoreCase(text, blacklist)
		}
	})
}

func BenchmarkRuneTrieIgnoreCase(b *testing.B) {
	blacklist := getTestBlacklist(b)
	testcases := newTestcaseParallel(getTestcases(b))

	tt := newRuneTrie(getOpt(IgnoreCase(true)))
	tt.Insert(blacklist...)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tt.search(testcases.Get())
		}
	})
}

func BenchmarkByteTrieIgnoreCase(b *testing.B) {
	blacklist := getTestBlacklist(b)
	testcases := newTestcaseParallel(getTestcases(b))

	tt := newByteTrie(getOpt(IgnoreCase(true)))
	tt.Insert(blacklist...)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tt.search(testcases.Get())
		}
	})
}
