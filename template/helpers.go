// Copyright (c) 2018, Drone.IO Inc
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package template

import (
	"fmt"
	"math"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/Masterminds/sprig/v3"
	"github.com/flowchartsman/handlebars/v3"
)

func init() {
	funcs := map[string]interface{}{
		"duration":       toDuration,
		"datetime":       toDatetime,
		"success":        isSuccess,
		"failure":        isFailure,
		"truncate":       truncate,
		"urlencode":      urlencode,
		"since":          since,
		"uppercasefirst": uppercaseFirst,
		"uppercase":      strings.ToUpper,
		"lowercase":      strings.ToLower,
		"regexReplace":   regexReplace,
	}
	for name, f := range sprig.GenericFuncMap() {
		if _, ok := funcs[name]; ok || !validHelper(f) {
			continue
		}
		funcs[name] = f
	}
	handlebars.RegisterHelpers(funcs)
}

func toDuration(started, finished int64) string {
	return fmt.Sprint(time.Duration(finished-started) * time.Second)
}

func toDatetime(timestamp int64, layout, zone string) string {
	if len(zone) == 0 {
		return time.Unix(timestamp, 0).Format(layout)
	}

	loc, err := time.LoadLocation(zone)
	if err != nil {
		return time.Unix(timestamp, 0).Local().Format(layout)
	}

	return time.Unix(timestamp, 0).In(loc).Format(layout)
}

func isSuccess(conditional bool, options *handlebars.Options) string {
	if !conditional {
		return options.Inverse()
	}

	switch options.ParamStr(0) {
	case "success":
		return options.Fn()
	default:
		return options.Inverse()
	}
}

func isFailure(conditional bool, options *handlebars.Options) string {
	if !conditional {
		return options.Inverse()
	}

	switch options.ParamStr(0) {
	case "failure", "error", "killed":
		return options.Fn()
	default:
		return options.Inverse()
	}
}

func truncate(s string, len int) string {
	if utf8.RuneCountInString(s) <= int(math.Abs(float64(len))) {
		return s
	}

	runes := []rune(s)

	if len < 0 {
		len = -len
		return string(runes[len:])
	}

	return string(runes[:len])
}

func urlencode(options *handlebars.Options) string {
	return url.QueryEscape(options.Fn())
}

func since(start int64) string {
	now := time.Unix(time.Now().Unix(), 0)
	return fmt.Sprint(now.Sub(time.Unix(start, 0)))
}

func uppercaseFirst(s string) string {
	a := []rune(s)

	a[0] = unicode.ToUpper(a[0])
	s = string(a)

	return s
}

func regexReplace(pattern, input, replacement string) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(input, replacement)
}

func validHelper(f interface{}) bool {
	typ := reflect.TypeOf(f)
	if typ.NumOut() != 1 {
		return false
	}
	v := reflect.Zero(typ.Out(0))
	switch v.Interface().(type) {
	case
		bool,
		float64,
		[]int,
		int,
		int64,
		[][]interface{},
		[]interface{},
		map[string]interface{},
		map[string]string,
		[]string,
		string,
		time.Time,
		interface{}:
		return true
	}
	return false
}
