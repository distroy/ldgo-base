/*
 * Copyright (C) distroy
 */

package ldtrie

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"testing"
)

var (
	_blacklist           []string
	_blacklistIgnoreCase []string
	_testcases           []string
)

func strToLower(s string) string {
	return strings.ToLower(s)
}

func testGetBlacklist(t testing.TB, ignoreCase bool) []string {
	if ignoreCase {
		return getTestBlacklistIgnoreCase(t)
	}
	return getTestBlacklist(t)
}

func getTestBlacklist(t testing.TB) []string {
	if len(_blacklist) != 0 {
		return _blacklist
	}

	f, err := os.Open("./blacklist.csv")
	if err != nil {
		t.Fatalf("open blacklist fail. err:%s", err.Error())
	}
	defer f.Close()

	res := testReadStrings(t, f)
	_blacklist = res
	return res
}

func testReadStrings(t testing.TB, r io.Reader) []string {
	reader := csv.NewReader(r)
	lines, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("read blacklist fail. err:%s", err.Error())
	}

	res := make([]string, 0, len(lines))
	for _, line := range lines {
		res = append(res, line[0])
	}

	return res
}

func getTestBlacklistIgnoreCase(t testing.TB) []string {
	if len(_blacklistIgnoreCase) != 0 {
		return _blacklistIgnoreCase
	}
	txts := getTestBlacklist(t)
	res := make([]string, 0, len(txts))
	m := make(map[string]struct{}, len(txts))
	for _, text := range txts {
		text = strToLower(text)
		if _, ok := m[text]; ok {
			continue
		}
		res = append(res, text)
	}

	_blacklistIgnoreCase = res
	return res
}

func getTestcases(t testing.TB) []string {
	if len(_testcases) != 0 {
		return _testcases
	}

	f, err := os.Open("./testcase.csv")
	if err != nil {
		t.Fatalf("open blacklist fail. err:%s", err.Error())
	}
	defer f.Close()

	res := testReadStrings(t, f)
	_testcases = res
	return res
}

func newTestcaseParallel(testcases []string) *testcaseParallel {
	return &testcaseParallel{
		Testcases: testcases,
		Length:    uint32(len(testcases)),
	}
}

type testcaseParallel struct {
	Testcases []string
	Length    uint32
	Position  uint32
}

func (that *testcaseParallel) Get() string {
	n := atomic.AddUint32(&that.Position, 1) - 1
	n = n % that.Length
	return that.Testcases[n]
}

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
}
