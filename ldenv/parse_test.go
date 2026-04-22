/*
 * Copyright (C) distroy
 */

package ldenv

import (
	"os"
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
	"github.com/distroy/ldgo-base/ldptr"
)

func TestParse(t *testing.T) {
	type Environment struct {
		AuthService     string  `ldenv:"name:AUTH_SERVICE; default:auth-service"`
		AuthTimeout     float64 `ldenv:"default:60"`
		ConfigApp       uint64  `ldenv:"default:1001"`
		ConfigNamespace string  `ldenv:"default:test-ns"`
		DbUrl           string  `ldenv:"name:DB_HOST; default:database.db"`
		DbPort          int     `ldenv:"default:10000"`
		DbName          *string `ldenv:"name:DB_NAME; default:test-db"`
	}
	convey.Convey(t.Name(), t, func(c convey.C) {
		os.Setenv(`CONFIG_APP`, "1002")
		os.Setenv(`DB_HOST`, "database.db1")

		env := &Environment{}
		Parse(env)

		c.So(env, convey.ShouldResemble, &Environment{
			AuthService:     "auth-service",
			AuthTimeout:     60,
			ConfigApp:       1002,
			ConfigNamespace: "test-ns",
			DbUrl:           "database.db1",
			DbPort:          10000,
			DbName:          ldptr.New("test-db"),
		})
	})
}
