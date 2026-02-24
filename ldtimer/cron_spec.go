/*
 * Copyright (C) distroy
 */

package ldtimer

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/distroy/ldgo-base/ldsort"
)

// CronSpec 表示解析后的 cron 时间规格
type CronSpec struct {
	Minute     []int // 0-59
	Hour       []int // 0-23
	DayOfMonth []int // 1-31
	Month      []int // 1-12
	DayOfWeek  []int // 0-6 (0=星期日)
}

// Matches 检查给定的时间是否匹配 CronSpec（可选辅助方法）
func (s *CronSpec) Matches(minute, hour, day, month, weekday int) bool {
	return slices.Contains(s.Minute, minute) &&
		slices.Contains(s.Hour, hour) &&
		slices.Contains(s.DayOfMonth, day) &&
		slices.Contains(s.Month, month) &&
		slices.Contains(s.DayOfWeek, weekday)
}

// NextRun 计算下一次运行时间
func (s *CronSpec) NextRun(from time.Time) time.Time {
	next := from.Add(1 * time.Minute)
	next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())

	for i := 0; i < 100000; i++ { // 防止无限循环
		if s.MatchTime(next) {
			return next
		}
		next = next.Add(1 * time.Minute)
	}

	return time.Time{}
}

// Matches 检查给定时间是否匹配计划
func (s *CronSpec) MatchTime(t time.Time) bool {
	// 检查分钟
	if !slices.Contains(s.Minute, t.Minute()) {
		return false
	}

	// 检查小时
	if !slices.Contains(s.Hour, t.Hour()) {
		return false
	}

	// 检查月份
	if !slices.Contains(s.Month, int(t.Month())) {
		return false
	}

	// 检查日期和星期几的特殊逻辑
	dayMatch := slices.Contains(s.DayOfMonth, t.Day())
	weekdayMatch := slices.Contains(s.DayOfWeek, int(t.Weekday()))

	// 特殊处理最后一天 (L)
	if slices.Contains(s.DayOfMonth, -1) {
		// 获取当月的最后一天
		nextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
		lastDay := nextMonth.Add(-24 * time.Hour).Day()
		if t.Day() == lastDay {
			dayMatch = true
		}
	}

	// cron 特殊规则：如果日期和星期几都被指定，匹配任意一个即可
	// 如果只有其中一个被指定，则必须匹配该字段
	daySpecified := !(len(s.DayOfMonth) == 31 && !slices.Contains(s.DayOfMonth, -1)) // 不是全匹配且不是 L
	weekdaySpecified := len(s.DayOfWeek) < 7 || slices.Contains(s.DayOfWeek, 7)      // 不是全匹配

	if daySpecified && weekdaySpecified {
		return dayMatch || weekdayMatch
	} else if daySpecified {
		return dayMatch
	} else if weekdaySpecified {
		return weekdayMatch
	}

	return true
}

// String 返回可读的 cron 表达式
func (s *CronSpec) String() string {
	return fmt.Sprintf("minute:%v hour:%v day:%v month:%v weekday:%v",
		s.Minute, s.Hour, s.DayOfMonth, s.Month, s.DayOfWeek)
}

// ParseCronExpr 解析 crontab 的前五个时间字段，忽略命令部分
func ParseCronExpr(cronExpr string) (*CronSpec, error) {
	fields := strings.Fields(cronExpr)
	if len(fields) < 5 {
		return nil, fmt.Errorf("invalid cron expression: need at least 5 fields, got %d", len(fields))
	}

	spec := &CronSpec{}
	var err error

	// 分钟
	spec.Minute, err = parseField(fields[0], 0, 59, nil)
	if err != nil {
		return nil, fmt.Errorf("minute: %v", err)
	}
	// 小时
	spec.Hour, err = parseField(fields[1], 0, 23, nil)
	if err != nil {
		return nil, fmt.Errorf("hour: %v", err)
	}
	// 日
	spec.DayOfMonth, err = parseField(fields[2], 1, 31, nil)
	if err != nil {
		return nil, fmt.Errorf("day of month: %v", err)
	}
	// 月
	monthAlias := map[string]int{
		"jan": 1, "feb": 2, "mar": 3, "apr": 4,
		"may": 5, "jun": 6, "jul": 7, "aug": 8,
		"sep": 9, "oct": 10, "nov": 11, "dec": 12,
	}
	spec.Month, err = parseField(fields[3], 1, 12, monthAlias)
	if err != nil {
		return nil, fmt.Errorf("month: %v", err)
	}
	// 星期
	weekAlias := map[string]int{
		"sun": 0, "mon": 1, "tue": 2, "wed": 3,
		"thu": 4, "fri": 5, "sat": 6, // "sun": 7
	}
	spec.DayOfWeek, err = parseField(fields[4], 0, 7, weekAlias)
	if err != nil {
		return nil, fmt.Errorf("day of week: %v", err)
	}

	// 将星期日的两种表示法统一为 0
	spec.DayOfWeek = normalizeSunday(spec.DayOfWeek)
	return spec, nil
}

// parseField 解析单个字段（通用）
func parseField(field string, min, max int, aliases map[string]int) ([]int, error) {
	field = strings.TrimSpace(field)
	if field == "" {
		return nil, fmt.Errorf("empty field")
	}

	// 处理 *
	if field == "*" {
		return expandRange(min, max, 1), nil
	}

	// 处理步进格式 */5
	if strings.HasPrefix(field, "*/") {
		step, err := strconv.Atoi(field[2:])
		if err != nil || step <= 0 {
			return nil, fmt.Errorf("invalid step: %s", field)
		}
		return expandRange(min, max, step), nil
	}

	// 处理逗号分隔列表
	if strings.Contains(field, ",") {
		parts := strings.Split(field, ",")
		var result []int
		for _, p := range parts {
			vals, err := parseSingle(p, min, max, aliases)
			if err != nil {
				return nil, err
			}
			result = append(result, vals...)
		}
		return uniqueSort(result), nil
	}

	// 单个表达式（数字、范围、带步进的范围）
	return parseSingle(field, min, max, aliases)
}

// parseSingle 解析单个表达式（不带逗号）
func parseSingle(expr string, min, max int, aliases map[string]int) ([]int, error) {
	// 处理范围步进：1-10/2
	if strings.Contains(expr, "/") {
		return parseSingleWithStep(expr, min, max, aliases)
	}

	// 处理范围：1-5
	if strings.Contains(expr, "-") {
		return parseSingleWithRange(expr, min, max, aliases)
	}

	// 单个数字或别名
	val, err := parseValue(expr, min, max, aliases)
	if err != nil {
		return nil, err
	}
	return []int{val}, nil
}

func parseSingleWithStep(expr string, min, max int, aliases map[string]int) ([]int, error) {
	parts := strings.Split(expr, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid step expression: %s", expr)
	}
	step, err := strconv.Atoi(parts[1])
	if err != nil || step <= 0 {
		return nil, fmt.Errorf("invalid step number: %s", parts[1])
	}

	rangePart := parts[0]
	var start, end int
	if rangePart == "*" {
		start, end = min, max
	} else if strings.Contains(rangePart, "-") {
		rangeVals := strings.Split(rangePart, "-")
		if len(rangeVals) != 2 {
			return nil, fmt.Errorf("invalid range: %s", rangePart)
		}
		start, err = parseValue(rangeVals[0], min, max, aliases)
		if err != nil {
			return nil, err
		}
		end, err = parseValue(rangeVals[1], min, max, aliases)
		if err != nil {
			return nil, err
		}
	} else {
		start, err = parseValue(rangePart, min, max, aliases)
		if err != nil {
			return nil, err
		}
		end = max
	}
	if start > end {
		return nil, fmt.Errorf("range start > end: %s", rangePart)
	}
	return expandRange(start, end, step), nil
}

func parseSingleWithRange(expr string, min, max int, aliases map[string]int) ([]int, error) {
	parts := strings.Split(expr, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid range: %s", expr)
	}
	start, err := parseValue(parts[0], min, max, aliases)
	if err != nil {
		return nil, err
	}
	end, err := parseValue(parts[1], min, max, aliases)
	if err != nil {
		return nil, err
	}
	if start > end {
		return nil, fmt.Errorf("range start > end: %s", expr)
	}
	return expandRange(start, end, 1), nil
}

// parseValue 解析单个值（数字或别名）
func parseValue(s string, min, max int, aliases map[string]int) (int, error) {
	// 尝试别名
	if aliases != nil {
		if v, ok := aliases[strings.ToLower(s)]; ok {
			if v < min || v > max {
				return 0, fmt.Errorf("alias %s = %d out of range [%d,%d]", s, v, min, max)
			}
			return v, nil
		}
	}
	// 尝试数字
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", s)
	}
	if v < min || v > max {
		return 0, fmt.Errorf("value %d out of range [%d,%d]", v, min, max)
	}
	return v, nil
}

// expandRange 生成 [start, end] 步长为 step 的整数切片
func expandRange(start, end, step int) []int {
	var res []int
	for i := start; i <= end; i += step {
		res = append(res, i)
	}
	return res
}

// normalizeSunday 将星期日统一表示为 0（去除 7）
func normalizeSunday(vals []int) []int {
	// 用 0 替换 7
	for i, v := range vals {
		if v == 7 {
			vals[i] = 0
		}
	}
	return uniqueSort(vals)
}

// uniqueSort 去重并排序
func uniqueSort(nums []int) []int {
	ldsort.SortIntegers(nums)
	return ldsort.UniqIntegers(nums)
}
