/*
 * Copyright (C) distroy
 */

package ldtimer

import (
	"reflect"
	"testing"
)

// 测试完整的 Parse 函数
func TestParseCronExpr(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		want    *CronSpec
		wantErr bool
	}{
		{
			name: "standard five fields with command (ignored)",
			expr: "*/15 0 1,15 * 1-5 /usr/bin/backup",
			want: &CronSpec{
				Minute:     []int{0, 15, 30, 45},
				Hour:       []int{0},
				DayOfMonth: []int{1, 15},
				Month:      expandRange(1, 12, 1),
				DayOfWeek:  []int{1, 2, 3, 4, 5},
			},
			wantErr: false,
		},
		{
			name: "every minute",
			expr: "* * * * *",
			want: &CronSpec{
				Minute:     expandRange(0, 59, 1),
				Hour:       expandRange(0, 23, 1),
				DayOfMonth: expandRange(1, 31, 1),
				Month:      expandRange(1, 12, 1),
				DayOfWeek:  expandRange(0, 6, 1),
			},
			wantErr: false,
		},
		{
			name: "every 5 minutes",
			expr: "*/5 * * * *",
			want: &CronSpec{
				Minute:     []int{0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55},
				Hour:       expandRange(0, 23, 1),
				DayOfMonth: expandRange(1, 31, 1),
				Month:      expandRange(1, 12, 1),
				DayOfWeek:  expandRange(0, 6, 1),
			},
			wantErr: false,
		},
		{
			name: "range and list",
			expr: "0-10/2 9-17 1,15 * 1-5",
			want: &CronSpec{
				Minute:     []int{0, 2, 4, 6, 8, 10},
				Hour:       []int{9, 10, 11, 12, 13, 14, 15, 16, 17},
				DayOfMonth: []int{1, 15},
				Month:      expandRange(1, 12, 1),
				DayOfWeek:  []int{1, 2, 3, 4, 5},
			},
			wantErr: false,
		},
		{
			name: "month names",
			expr: "0 12 1 jan,dec 0",
			want: &CronSpec{
				Minute:     []int{0},
				Hour:       []int{12},
				DayOfMonth: []int{1},
				Month:      []int{1, 12},
				DayOfWeek:  []int{0},
			},
			wantErr: false,
		},
		{
			name: "weekday names",
			expr: "30 14 * * mon,fri",
			want: &CronSpec{
				Minute:     []int{30},
				Hour:       []int{14},
				DayOfMonth: expandRange(1, 31, 1),
				Month:      expandRange(1, 12, 1),
				DayOfWeek:  []int{1, 5},
			},
			wantErr: false,
		},
		{
			name: "sunday both 0 and 7",
			expr: "0 0 * * 0,7",
			want: &CronSpec{
				Minute:     []int{0},
				Hour:       []int{0},
				DayOfMonth: expandRange(1, 31, 1),
				Month:      expandRange(1, 12, 1),
				DayOfWeek:  []int{0}, // 7 被统一为 0
			},
			wantErr: false,
		},
		{
			name:    "less than 5 fields",
			expr:    "* * * *",
			wantErr: true,
		},
		{
			name:    "invalid minute step",
			expr:    "*/0 * * * *",
			wantErr: true,
		},
		{
			name:    "minute out of range",
			expr:    "60 * * * *",
			wantErr: true,
		},
		{
			name:    "invalid month name",
			expr:    "* * * xxx *",
			wantErr: true,
		},
		{
			name:    "range start > end",
			expr:    "* * 31-1 * *",
			wantErr: true,
		},
		{
			name:    "invalid step in range",
			expr:    "* * 1-5/0 * *",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCronExpr(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Parse() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

// 测试 parseField 独立功能
func Test_parseField(t *testing.T) {
	aliases := map[string]int{"jan": 1, "feb": 2}
	tests := []struct {
		name    string
		field   string
		min     int
		max     int
		aliases map[string]int
		want    []int
		wantErr bool
	}{
		{"star", "*", 1, 3, nil, []int{1, 2, 3}, false},
		{"step", "*/2", 0, 5, nil, []int{0, 2, 4}, false},
		{"list", "1,3,5", 0, 6, nil, []int{1, 3, 5}, false},
		{"range", "2-4", 0, 6, nil, []int{2, 3, 4}, false},
		{"range with step", "1-5/2", 0, 6, nil, []int{1, 3, 5}, false},
		{"star with step", "*/3", 1, 7, nil, []int{1, 4, 7}, false},
		{"alias single", "jan", 1, 12, aliases, []int{1}, false},
		{"alias in list", "jan,feb", 1, 12, aliases, []int{1, 2}, false},
		{"invalid step", "*/a", 0, 5, nil, nil, true},
		{"step zero", "*/0", 0, 5, nil, nil, true},
		{"invalid range", "1-a", 0, 5, nil, nil, true},
		{"out of range", "6", 0, 5, nil, nil, true},
		{"unknown alias", "mar", 1, 12, aliases, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseField(tt.field, tt.min, tt.max, tt.aliases)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseField() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试 parseSingle 函数
func Test_parseSingle(t *testing.T) {
	aliases := map[string]int{"sun": 0}
	tests := []struct {
		name    string
		expr    string
		min     int
		max     int
		aliases map[string]int
		want    []int
		wantErr bool
	}{
		{"single number", "5", 0, 9, nil, []int{5}, false},
		{"alias", "sun", 0, 6, aliases, []int{0}, false},
		{"range", "1-3", 0, 5, nil, []int{1, 2, 3}, false},
		{"range with step", "0-6/2", 0, 6, nil, []int{0, 2, 4, 6}, false},
		{"star with step", "*/4", 0, 8, nil, []int{0, 4, 8}, false},
		{"invalid range", "5-1", 0, 9, nil, nil, true},
		{"invalid number", "x", 0, 9, nil, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSingle(tt.expr, tt.min, tt.max, tt.aliases)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSingle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSingle() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试 parseValue 函数
func Test_parseValue(t *testing.T) {
	aliases := map[string]int{"mon": 1}
	tests := []struct {
		name    string
		s       string
		min     int
		max     int
		aliases map[string]int
		want    int
		wantErr bool
	}{
		{"number", "5", 0, 9, nil, 5, false},
		{"alias", "mon", 0, 6, aliases, 1, false},
		{"out of range", "10", 0, 5, nil, 0, true},
		{"invalid number", "abc", 0, 5, nil, 0, true},
		{"alias out of range", "mon", 2, 5, aliases, 0, true}, // mon=1 < min=2
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseValue(tt.s, tt.min, tt.max, tt.aliases)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试 expandRange
func Test_expandRange(t *testing.T) {
	tests := []struct {
		start, end, step int
		want             []int
	}{
		{1, 5, 1, []int{1, 2, 3, 4, 5}},
		{0, 6, 2, []int{0, 2, 4, 6}},
		{5, 5, 1, []int{5}},
		{1, 3, 3, []int{1, 4}}, // 注意：1+3=4 > 3，实际只输出1（实现是i+=step，i<=end，所以会输出1和4? 需要验证原代码：for i:=start; i<=end; i+=step，所以当start=1, end=3, step=3时，i=1 (<=3), i=4 (<=3? 4>3, stop) 所以输出[1]。所以期望是[1]）
	}
	// 修正用例
	tests[3].want = []int{1}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := expandRange(tt.start, tt.end, tt.step)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expandRange(%d,%d,%d) = %v, want %v", tt.start, tt.end, tt.step, got, tt.want)
			}
		})
	}
}

// 测试 normalizeSunday
func Test_normalizeSunday(t *testing.T) {
	tests := []struct {
		input []int
		want  []int
	}{
		{[]int{0}, []int{0}},
		{[]int{7}, []int{0}},
		{[]int{0, 7}, []int{0}},
		{[]int{1, 7, 2}, []int{0, 1, 2}},
		{[]int{0, 1, 7, 2}, []int{0, 1, 2}},
		{[]int{}, []int{}},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := normalizeSunday(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("normalizeSunday(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// 测试 uniqueSort
func Test_uniqueSort(t *testing.T) {
	tests := []struct {
		input []int
		want  []int
	}{
		{[]int{3, 1, 2, 1, 3}, []int{1, 2, 3}},
		{[]int{}, []int{}},
		{[]int{5}, []int{5}},
		{[]int{2, 2, 2}, []int{2}},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := uniqueSort(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uniqueSort(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// 测试 Matches 方法
func TestCronSpec_Matches(t *testing.T) {
	spec := &CronSpec{
		Minute:     []int{0, 30},
		Hour:       []int{9, 17},
		DayOfMonth: []int{1, 15},
		Month:      []int{1, 6, 12},
		DayOfWeek:  []int{1, 3, 5}, // 周一、三、五
	}
	tests := []struct {
		minute, hour, day, month, weekday int
		want                              bool
	}{
		{0, 9, 1, 1, 1, true},    // 1月1日周一09:00
		{30, 17, 15, 6, 3, true}, // 6月15日周三17:30
		{0, 9, 2, 1, 1, false},   // day 2 不在列表中
		{0, 9, 1, 2, 1, false},   // month 2 不在列表中
		{0, 8, 1, 1, 1, false},   // hour 8 不在列表中
		{0, 9, 1, 1, 2, false},   // weekday 2 不在列表中（周二）
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := spec.Matches(tt.minute, tt.hour, tt.day, tt.month, tt.weekday)
			if got != tt.want {
				t.Errorf("Matches(%d,%d,%d,%d,%d) = %v, want %v", tt.minute, tt.hour, tt.day, tt.month, tt.weekday, got, tt.want)
			}
		})
	}
}
