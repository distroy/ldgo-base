/*
 * Copyright (C) distroy
 */

package ldtime

import "time"

func MinuteBegin(t time.Time) time.Time {
	year, month, day := t.Date()
	hour, min, _ := t.Clock()
	sec, nsec := 0, 0
	return time.Date(int(year), time.Month(month), int(day), hour, min, sec, nsec, t.Location())
}

func MinuteEnd(t time.Time) time.Time {
	return MinuteBegin(t).Add(time.Minute - 1)
}

func HourBegin(t time.Time) time.Time {
	year, month, day := t.Date()
	hour, _, _ := t.Clock()
	min, sec, nsec := 0, 0, 0
	return time.Date(int(year), time.Month(month), int(day), hour, min, sec, nsec, t.Location())
}

func HourEnd(t time.Time) time.Time {
	t.Location()
	return HourBegin(t).Add(time.Hour - 1)
}

func DateBegin(t time.Time) time.Time {
	year, month, day := t.Date()
	hour, min, sec, nsec := 0, 0, 0, 0
	return time.Date(int(year), time.Month(month), int(day), hour, min, sec, nsec, t.Location())
}

func DateEnd(t time.Time) time.Time {
	return DateBegin(t).Add(time.Hour*24 - 1)
}

func WeekBegin(t time.Time) time.Time {
	t = DateBegin(t)
	offset := int(t.Weekday() - time.Sunday)
	// offset = offset + 7
	// offset = offset % 7
	return t.Add(time.Duration(-offset) * 24 * time.Hour)
}

func WeekEnd(t time.Time) time.Time {
	return WeekBegin(t).Add(time.Hour*24*7 - 1)
}

func MonthBegin(t time.Time) time.Time {
	year, month, _ := t.Date()
	day := 1
	hour, min, sec, nsec := 0, 0, 0, 0
	return time.Date(int(year), time.Month(month), int(day), hour, min, sec, nsec, t.Location())
}

func MonthEnd(t time.Time) time.Time {
	return MonthBegin(t).AddDate(0, 1, 0).Add(-1)
}

func YearBegin(t time.Time) time.Time {
	year, _, _ := t.Date()
	month, day := time.January, 1
	hour, min, sec, nsec := 0, 0, 0, 0
	return time.Date(int(year), month, int(day), hour, min, sec, nsec, t.Location())
}

func YearEnd(t time.Time) time.Time {
	return YearBegin(t).AddDate(1, 0, 0).Add(-1)
}
