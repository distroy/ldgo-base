/*
 * Copyright (C) distroy
 */

package ldenv

import (
	"os"
	"testing"
	"time"

	"github.com/distroy/ldgo-base/3rd/convey"
	"github.com/distroy/ldgo-base/internal/time_"
	"github.com/distroy/ldgo-base/ldptr"
	"github.com/distroy/ldgo-base/ldtime"
)

func TestParse(t *testing.T) {
	type Environment struct {
		AuthService     string          `ldenv:"name:AUTH_SERVICE; default:auth-service"`
		AuthTimeout     float64         `ldenv:"default:60"`
		ConfigApp       uint64          `ldenv:"default:1001"`
		ConfigNamespace string          `ldenv:"default:test-ns"`
		DbUrl           string          `ldenv:"name:DB_HOST; default:database.db"`
		DbPort          int             `ldenv:"default:10000"`
		DbName          *string         `ldenv:"name:DB_NAME; default:test-db"`
		Time            *time.Time      `ldenv:""`
		Duration        time.Duration   `ldenv:""`
		Duration2       ldtime.Duration `ldenv:""`
	}
	convey.Convey(t.Name(), t, func(c convey.C) {
		os.Setenv(`CONFIG_APP`, "1002")
		os.Setenv(`DB_HOST`, "database.db1")
		os.Setenv(`TIME`, "2026-06-25T23:33:28+0800")
		os.Setenv(`DURATION`, "1m23s")
		os.Setenv(`DURATION2`, "5.12")

		env := &Environment{}

		err := Parse(env)

		c.So(err, convey.ShouldBeNil)
		c.So(env, convey.ShouldResemble, &Environment{
			AuthService:     "auth-service",
			AuthTimeout:     60,
			ConfigApp:       1002,
			ConfigNamespace: "test-ns",
			DbUrl:           "database.db1",
			DbPort:          10000,
			DbName:          ldptr.New("test-db"),
			Time:            ldptr.New(firstRes(time.Parse(time_.TimeLayout, "2026-06-25T23:33:28+0800"))),
			Duration:        time.Minute*1 + time.Second*23,
			Duration2:       ldtime.Duration(time.Second*5 + time.Millisecond*120),
		})
	})
}

func firstRes[T any](a1 T, a ...any) T { return a1 }
