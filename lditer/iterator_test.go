/*
 * Copyright (C) distroy
 */

package lditer

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"
)

type testPair[K, V any] struct {
	key   K
	value V
}

func testReadSeq[T any](iter Seq[T], yield func(v T) bool) []T {
	if yield == nil {
		yield = func(v T) bool { return true }
	}
	res := make([]T, 0, 16)
	for v := range iter {
		if !yield(v) {
			break
		}
		res = append(res, v)
	}
	return res
}

func testReadSeq2[K, V any](iter Seq2[K, V], yield func(k K, v V) bool) []testPair[K, V] {
	if yield == nil {
		yield = func(k K, v V) bool { return true }
	}
	res := make([]testPair[K, V], 0, 16)
	for k, v := range iter {
		if !yield(k, v) {
			break
		}
		res = append(res, testPair[K, V]{k, v})
	}
	return res
}

func TestChan(t *testing.T) {
	tests := []struct {
		name  string
		yield func(v int) bool
		slice []int
		want  []int
	}{
		{
			name:  "v > 0",
			yield: func(v int) bool { return v > 0 },
			slice: []int{1, 2, 3, 4, 0, 5},
			want:  []int{1, 2, 3, 4},
		},
	}

	fnMakeChan := func(s []int) chan int {
		ch := make(chan int, len(s))
		for _, n := range s {
			ch <- n
		}
		close(ch)
		return ch
	}

	for i, tt := range tests {
		name := fmt.Sprintf("%d: %s", i, tt.name)
		t.Run(name, func(t *testing.T) {
			ch := fnMakeChan(tt.slice)
			got := testReadSeq(Chan(ch), tt.yield)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestChan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToSeq2(t *testing.T) {
	tests := []struct {
		name  string
		yield func(k, v int) bool
		slice []int
		want  []testPair[int, int]
	}{
		{
			name:  "v > 0",
			yield: func(i, v int) bool { return v > 0 },
			slice: []int{1, 2, 3, 4, 0, 5},
			want:  []testPair[int, int]{{0, 1}, {1, 2}, {2, 3}, {3, 4}},
		},
	}

	fnMakeChan := func(s []int) chan int {
		ch := make(chan int, len(s))
		for _, n := range s {
			ch <- n
		}
		close(ch)
		return ch
	}

	for i, tt := range tests {
		name := fmt.Sprintf("%d: %s", i, tt.name)
		t.Run(name, func(t *testing.T) {
			ch := fnMakeChan(tt.slice)
			// got := testReadSeq2(Chan2(ch), tt.yield)
			got := testReadSeq2(ToSeq2(Chan(ch)), tt.yield)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestToSeq2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		name  string
		n     int
		ns    []int
		yield func(v int) bool
		want  []int
	}{
		{
			name:  "5",
			n:     5,
			yield: func(v int) bool { return true },
			want:  []int{0, 1, 2, 3, 4},
		},
		{
			name:  "8, 13",
			n:     8,
			ns:    []int{13},
			yield: func(v int) bool { return v < 10 },
			want:  []int{8, 9},
		},
		{
			name:  "3, 13, 2",
			n:     3,
			ns:    []int{13, 2},
			yield: func(v int) bool { return v < 10 },
			want:  []int{3, 5, 7, 9},
		},
	}

	for i, tt := range tests {
		name := fmt.Sprintf("%d: %s", i, tt.name)
		t.Run(name, func(t *testing.T) {
			got := testReadSeq(Int(tt.n, tt.ns...), tt.yield)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice(t *testing.T) {
	tests := []struct {
		name  string
		yield func(i, v int) bool
		slice []int
		want  []testPair[int, int]
	}{
		{
			name:  "v > 0",
			yield: func(i, v int) bool { return v > 0 },
			slice: []int{1, 2, 3, 4, 0, 5},
			want:  []testPair[int, int]{{0, 1}, {1, 2}, {2, 3}, {3, 4}},
		},
	}

	for i, tt := range tests {
		name := fmt.Sprintf("%d: %s", i, tt.name)
		t.Run(name, func(t *testing.T) {
			got := testReadSeq2(Slice(tt.slice), tt.yield)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceBackward(t *testing.T) {
	tests := []struct {
		name  string
		yield func(i, v int) bool
		slice []int
		want  []testPair[int, int]
	}{
		{
			name:  "v > 0",
			yield: func(i, v int) bool { return v > 0 },
			slice: []int{1, 0, 2, 3, 4, 5},
			want:  []testPair[int, int]{{5, 5}, {4, 4}, {3, 3}, {2, 2}},
		},
	}

	for i, tt := range tests {
		name := fmt.Sprintf("%d: %s", i, tt.name)
		t.Run(name, func(t *testing.T) {
			got := testReadSeq2(SliceBackward(tt.slice), tt.yield)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestSlice2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	type testcase struct {
		name      string
		yield     func(k int, v string) bool
		args      map[int]string
		want      []testPair[int, string]
		wantEqual bool
	}
	tests := []testcase{
		{
			name:      "true",
			yield:     func(k int, v string) bool { return true },
			args:      map[int]string{1: "a", 2: "b", 3: "x", 4: "z"},
			want:      []testPair[int, string]{{1, "a"}, {2, "b"}, {3, "x"}, {4, "z"}},
			wantEqual: true,
		},
		{
			name:      "break",
			yield:     func(k int, v string) bool { return k%2 == 0 },
			args:      map[int]string{1: "a", 2: "b", 3: "x", 4: "z"},
			want:      []testPair[int, string]{{1, "a"}, {2, "b"}, {3, "x"}, {4, "z"}},
			wantEqual: false,
		},
	}

	fnDoOneCase := func(t *testing.T, tt testcase) {
		got := testReadSeq2(Map(tt.args), tt.yield)
		slices.SortFunc(got, func(a, b testPair[int, string]) int { return a.key - b.key })
		slices.SortFunc(tt.want, func(a, b testPair[int, string]) int { return a.key - b.key })

		if tt.wantEqual || len(got) == len(tt.want) {
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestMap() = %v, want equal = %v, want values = %v", got, tt.wantEqual, tt.want)
			}
			return
		}

		for _, v := range got {
			if idx := slices.IndexFunc(tt.want, func(x testPair[int, string]) bool { return x.key >= v.key }); idx < 0 || !reflect.DeepEqual(v, tt.want[idx]) {
				t.Errorf("TestMap() = %v, want equal = %v, want values = %v", got, tt.wantEqual, tt.want)
			}
		}
	}

	for i, tt := range tests {
		name := fmt.Sprintf("%d: %s", i, tt.name)
		t.Run(name, func(t *testing.T) {
			fnDoOneCase(t, tt)
		})
	}
}

func TestToSeqByKey(t *testing.T) {
	tests := []struct {
		name  string
		yield func(k int) bool
		args  map[int]string
		want  []int
	}{
		{
			name:  "true",
			yield: func(k int) bool { return true },
			args:  map[int]string{1: "a", 2: "b", 3: "x", 4: "z"},
			want:  []int{1, 2, 3, 4},
		},
	}

	for i, tt := range tests {
		name := fmt.Sprintf("%d: %s", i, tt.name)
		t.Run(name, func(t *testing.T) {
			// got := testReadSeq(MapKeys(tt.args), tt.yield)
			got := testReadSeq(ToSeqByKey(Map(tt.args)), tt.yield)
			slices.SortFunc(got, func(a, b int) int { return a - b })
			slices.SortFunc(tt.want, func(a, b int) int { return a - b })

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestToSeqByKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToSeqByValue(t *testing.T) {
	tests := []struct {
		name  string
		yield func(v string) bool
		args  map[int]string
		want  []string
	}{
		{
			name:  "true",
			yield: func(v string) bool { return true },
			args:  map[int]string{1: "a", 2: "b", 3: "x", 4: "z"},
			want:  []string{"a", "b", "x", "z"},
		},
	}

	for i, tt := range tests {
		name := fmt.Sprintf("%d: %s", i, tt.name)
		t.Run(name, func(t *testing.T) {
			got := testReadSeq(ToSeqByValue(Map(tt.args)), tt.yield)
			slices.SortFunc(got, func(a, b string) int { return strings.Compare(a, b) })
			slices.SortFunc(tt.want, func(a, b string) int { return strings.Compare(a, b) })

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestToSeqByValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		arg   string
		yield func(i int, v rune) bool
		want  []testPair[int, rune]
	}{
		{
			arg:  "abcde",
			want: []testPair[int, rune]{{0, 'a'}, {1, 'b'}, {2, 'c'}, {3, 'd'}, {4, 'e'}},
		},
		{
			arg:   "abcdefg",
			yield: func(i int, v rune) bool { return v < 'd' },
			want:  []testPair[int, rune]{{0, 'a'}, {1, 'b'}, {2, 'c'}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			got := testReadSeq2(String(tt.arg), tt.yield)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
