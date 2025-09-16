/*
 * Copyright (C) distroy
 */

package ldrcfg

import (
	"context"
	"testing"

	"github.com/distroy/ldgo-base/3rd/convey"
	"github.com/distroy/ldgo-base/3rd/yaml"
	"github.com/distroy/ldgo-base/ldctx"
	"github.com/distroy/ldgo-base/lderr"
	"github.com/distroy/ldgo-base/ldlog"
	"github.com/distroy/ldgo-base/ldptr"
)

type testJsonConfig struct {
	Str string `json:"str"`
	I64 int64  `json:"i64"`
}

type testYamlConfig struct {
	Str string `yaml:"str"`
	I64 int64  `yaml:"i64"`
}

type testParserConfig struct {
	Str string `yaml:"str"`
	I64 int64  `yaml:"i64"`
}

func (c *testParserConfig) Parse(ctx context.Context, v []byte) error {
	if err := yaml.Unmarshal(v, c); err != nil {
		ldctx.LogE(ctx, "parse object fail", ldlog.Error(err))
		return lderr.ErrUnmarshal
	}

	return nil
}

func TestClient_Register(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		type Config struct {
			Str     Event[string]            `ldrcfg:"key:str; type:string"`
			Json    Event[testJsonConfig]    `ldrcfg:"key:json; type:json"`
			Yaml    Event[testYamlConfig]    `ldrcfg:"key:yaml; type:yaml"`
			Parser  Event[testParserConfig]  `ldrcfg:"key:parser; type:parser"`
			Strp    Event[*string]           `ldrcfg:"key:strp"`
			Jsonp   Event[*testJsonConfig]   `ldrcfg:"key:jsonp"`
			Yamlp   Event[*testYamlConfig]   `ldrcfg:"key:yamlp"`
			Parserp Event[*testParserConfig] `ldrcfg:"key:parserp"`
		}

		var (
			// ctx     = ldctx.Discard()
			ctx     = ldctx.Default()
			cli     = NewClient(testAdaptor{})
			ns      = "test-ns"
			trigger = func(key string, val string) { testTriggerEvent(cli, ns, key, val) }
			cfg     = &Config{}
			vv      any
		)

		cfg.Str.OnChange(ctx, func(c context.Context, v string) { vv = v })
		cfg.Json.OnChange(ctx, func(c context.Context, v testJsonConfig) { vv = v })
		cfg.Yaml.OnChange(ctx, func(c context.Context, v testYamlConfig) { vv = v })
		cfg.Parser.OnChange(ctx, func(c context.Context, v testParserConfig) { vv = v })
		cfg.Strp.OnChange(ctx, func(c context.Context, v *string) { vv = v })
		cfg.Jsonp.OnChange(ctx, func(c context.Context, v *testJsonConfig) { vv = v })
		cfg.Yamlp.OnChange(ctx, func(c context.Context, v *testYamlConfig) { vv = v })
		cfg.Parserp.OnChange(ctx, func(c context.Context, v *testParserConfig) { vv = v })

		cli.logger = ldlog.Discard()
		cli.Register(ns, cfg)

		c.Convey("str", func(c convey.C) {
			var (
				vk = "str"
				pk = "strp"
				v  = &cfg.Str
				p  = &cfg.Strp
			)

			trigger(vk, `1234`)
			c.So(v.Load(), convey.ShouldEqual, "1234")
			c.So(v.Load(), convey.ShouldEqual, vv)

			trigger(pk, `abcd`)
			c.So(p.Load(), convey.ShouldResemble, ldptr.New("abcd"))
			c.So(p.Load(), convey.ShouldResemble, vv)
		})

		c.Convey("json", func(c convey.C) {
			var (
				vk = "json"
				pk = "jsonp"
				v  = &cfg.Json
				p  = &cfg.Jsonp
			)

			trigger(vk, `abc`)
			c.So(v.Load(), convey.ShouldResemble, testJsonConfig{})

			trigger(vk, `{"str":"abc"}`)
			c.So(v.Load(), convey.ShouldResemble, testJsonConfig{Str: "abc"})
			c.So(v.Load(), convey.ShouldResemble, vv)

			trigger(vk, `{"i64":123}`)
			c.So(v.Load(), convey.ShouldResemble, testJsonConfig{I64: 123})
			c.So(v.Load(), convey.ShouldResemble, vv)

			trigger(pk, `abc`)
			c.So(p.Load(), convey.ShouldResemble, &testJsonConfig{})

			trigger(pk, `{"str":"abc"}`)
			c.So(p.Load(), convey.ShouldResemble, &testJsonConfig{Str: "abc"})
			c.So(p.Load(), convey.ShouldResemble, vv)

			trigger(pk, `{"i64":123}`)
			c.So(p.Load(), convey.ShouldResemble, &testJsonConfig{I64: 123})
			c.So(p.Load(), convey.ShouldResemble, vv)

			trigger(pk, `str: xyz`)
			c.So(p.Load(), convey.ShouldResemble, &testJsonConfig{Str: "xyz"})
			c.So(p.Load(), convey.ShouldResemble, vv)
		})

		c.Convey("yaml", func(c convey.C) {
			var (
				vk = "yaml"
				pk = "yamlp"
				v  = &cfg.Yaml
				p  = &cfg.Yamlp
			)

			trigger(vk, `abc`)
			c.So(v.Load(), convey.ShouldResemble, testYamlConfig{})

			trigger(vk, `str: abc`)
			c.So(v.Load(), convey.ShouldResemble, testYamlConfig{Str: "abc"})
			c.So(v.Load(), convey.ShouldResemble, vv)

			trigger(vk, `i64: 123`)
			c.So(v.Load(), convey.ShouldResemble, testYamlConfig{I64: 123})
			c.So(v.Load(), convey.ShouldResemble, vv)

			trigger(pk, `abc`)
			c.So(p.Load(), convey.ShouldResemble, &testYamlConfig{})

			trigger(pk, `str: abc`)
			c.So(p.Load(), convey.ShouldResemble, &testYamlConfig{Str: "abc"})
			c.So(p.Load(), convey.ShouldResemble, vv)

			trigger(pk, `i64: 123`)
			c.So(p.Load(), convey.ShouldResemble, &testYamlConfig{I64: 123})
			c.So(p.Load(), convey.ShouldResemble, vv)
		})

		c.Convey("parser", func(c convey.C) {
			var (
				vk = "parser"
				pk = "parserp"
				v  = &cfg.Parser
				p  = &cfg.Parserp
			)

			trigger(vk, `abc`)
			c.So(v.Load(), convey.ShouldResemble, testParserConfig{})

			trigger(vk, `str: abc`)
			c.So(v.Load(), convey.ShouldResemble, testParserConfig{Str: "abc"})
			c.So(v.Load(), convey.ShouldResemble, vv)

			trigger(vk, `i64: 123`)
			c.So(v.Load(), convey.ShouldResemble, testParserConfig{I64: 123})
			c.So(v.Load(), convey.ShouldResemble, vv)

			trigger(pk, `abc`)
			c.So(p.Load(), convey.ShouldResemble, &testParserConfig{})

			trigger(pk, `str: abc`)
			c.So(p.Load(), convey.ShouldResemble, &testParserConfig{Str: "abc"})
			c.So(p.Load(), convey.ShouldResemble, vv)

			trigger(pk, `i64: 123`)
			c.So(p.Load(), convey.ShouldResemble, &testParserConfig{I64: 123})
			c.So(p.Load(), convey.ShouldResemble, vv)
		})
	})
}
