/*
 * Copyright (C) distroy
 */

package ldslice

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSplitByCount(t *testing.T) {
	type args struct {
		list []int
		size int
	}
	tests := []struct {
		args args
		want [][]int
	}{
		{
			args: args{list: []int{}, size: 2},
			want: nil,
		},
		{
			args: args{list: []int{1, 2}, size: 4},
			want: [][]int{{1, 2}},
		},
		{
			args: args{list: []int{1, 2, 3, 4, 5, 6, 7}, size: 2},
			want: [][]int{{1, 2}, {3, 4}, {5, 6}, {7}},
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%v-%d", tt.args.list, tt.args.size)
		t.Run(name, func(t *testing.T) {
			got := SplitByCount(tt.args.list, tt.args.size)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`Split(%v, %d) = %v, want=%v`, tt.args.list, tt.args.size, got, tt.want)
				return
			}
		})
	}
}

func TestSplitByCountFunc(t *testing.T) {
	type args struct {
		list []int
		size int
	}
	type want struct {
		count int
		index []int
		slice [][]int
	}
	tests := []struct {
		args args
		want want
	}{
		{
			args: args{list: []int{}, size: 2},
			want: want{
				count: 0,
				index: nil,
				slice: nil,
			},
		},
		{
			args: args{list: []int{1, 2}, size: 4},
			want: want{
				count: 1,
				index: []int{0},
				slice: [][]int{{1, 2}},
			},
		},
		{
			args: args{list: []int{1, 2, 3, 4, 5, 6, 7}, size: 2},
			want: want{
				count: 4,
				index: []int{0, 1, 2, 3},
				slice: [][]int{{1, 2}, {3, 4}, {5, 6}, {7}},
			},
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%v-%d", tt.args.list, tt.args.size)
		t.Run(name, func(t *testing.T) {
			cnt := SplitByCountFunc(tt.args.list, tt.args.size, nil)
			if cnt != tt.want.count {
				t.Errorf(`SplitFunc(%v, %d, nil) = %v, want=%v`, tt.args.list, tt.args.size, cnt, tt.want.count)
				return
			}
			got := want{}
			got.count = SplitByCountFunc(tt.args.list, tt.args.size, func(idx int, val []int) bool {
				got.index = append(got.index, idx)
				got.slice = append(got.slice, val)
				return true
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`SplitFunc(%v, %d, func) = %v, want=%v`, tt.args.list, tt.args.size, got, tt.want)
				return
			}
		})
	}
}

func TestSplitByCountIter(t *testing.T) {
	type args struct {
		list []int
		size int
	}
	type want struct {
		index []int
		slice [][]int
	}
	tests := []struct {
		args args
		want want
	}{
		{
			args: args{list: []int{}, size: 2},
			want: want{
				index: nil,
				slice: nil,
			},
		},
		{
			args: args{list: []int{1, 2}, size: 4},
			want: want{
				index: []int{0},
				slice: [][]int{{1, 2}},
			},
		},
		{
			args: args{list: []int{1, 2, 3, 4, 5, 6, 7}, size: 2},
			want: want{
				index: []int{0, 1, 2, 3},
				slice: [][]int{{1, 2}, {3, 4}, {5, 6}, {7}},
			},
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%v-%d", tt.args.list, tt.args.size)
		t.Run(name, func(t *testing.T) {
			got := want{}
			for idx, val := range SplitByCountIter(tt.args.list, tt.args.size) {
				got.index = append(got.index, idx)
				got.slice = append(got.slice, val)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`SplitFunc(%v, %d, func) = %v, want=%v`, tt.args.list, tt.args.size, got, tt.want)
				return
			}
		})
	}
}
