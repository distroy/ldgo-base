/*
 * Copyright (C) distroy
 */

package ldredis

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/distroy/ldgo-base/ldlog"
)

const (
	ldRedisSrcPath = "/ldredis/"
	exampleSrcPath = "/ldredis/example/"
)

var goRedisSrcPathes = []string{
	"/github.com/redis/go-redis",
}

func isCallerFilePath(file string) bool {
	for _, path := range goRedisSrcPathes {
		if strings.Contains(file, path) {
			return false
		}
	}
	if !strings.Contains(file, ldRedisSrcPath) {
		return true
	}
	if strings.HasSuffix(file, "_test.go") {
		return true
	}
	if strings.Contains(file, exampleSrcPath) {
		return true
	}
	return false
}

func getCallerField() ldlog.Attr {
	// if !caller {
	// 	return zap.Skip()
	// }
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if !isCallerFilePath(file) {
			continue
		}
		return ldlog.String("caller", fmt.Sprintf("%s:%d", file, line))
	}

	return ldlog.String("caller", "overflow")
}
