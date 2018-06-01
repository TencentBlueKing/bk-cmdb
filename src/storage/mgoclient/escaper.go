/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mgoclient

import (
	"bytes"
	"io"
	"reflect"
	"strings"
)

func EscapeHtml(origins ...interface{}) {
	for _, origin := range origins {
		ov := reflect.ValueOf(origin)
		for ov.Kind() == reflect.Ptr {
			ov = ov.Elem()
		}
		if ov.Kind() == reflect.Map {
			for _, k := range ov.MapKeys() {
				tv := ov.MapIndex(k)
				if !tv.IsValid() {
					continue
				}
				if tv.CanInterface() {
					fv := tv.Interface()
					switch v := fv.(type) {
					case string:
						ts := HTMLEscapeString(v)
						ov.SetMapIndex(k, reflect.ValueOf(ts))
					case *string:
						if nil != v {
							ts := HTMLEscapeString(*v)
							ov.SetMapIndex(k, reflect.ValueOf(ts))
						}
					}
				}
			}
		}
		if ov.Kind() == reflect.Struct {
			index := ov.NumField() - 1
			for ; index >= 0; index-- {
				tv := ov.Field(index)
				for ov.Kind() == reflect.Ptr {
					ov = ov.Elem()
				}
				if tv.Kind() == reflect.Struct {
					escapeEmbed(ov, []int{index})
				}
				if !tv.IsValid() {
					continue
				}
				if tv.CanInterface() && tv.CanSet() {
					fv := tv.Interface()
					switch v := fv.(type) {
					case string:
						ts := HTMLEscapeString(v)
						tv.Set(reflect.ValueOf(ts))
					case *string:
						if nil != v {
							ts := HTMLEscapeString(*(v))
							tv.Set(reflect.ValueOf(&ts))
						}
					}
				}
			}
		}
	}
}

func escapeEmbed(ov reflect.Value, oindex []int) {
	fv := ov.FieldByIndex(oindex)
	index := fv.NumField() - 1
	for ; index >= 0; index-- {
		tv := ov.FieldByIndex(append(oindex, index))
		for ov.Kind() == reflect.Ptr {
			ov = ov.Elem()
		}
		if tv.Kind() == reflect.Struct {
			escapeEmbed(ov, append(oindex, index))
		}
		if !tv.IsValid() {
			continue
		}
		if tv.CanInterface() && tv.CanSet() {
			fv := tv.Interface()
			switch fv.(type) {
			case string:
				ts := HTMLEscapeString(fv.(string))
				tv.SetString(ts)
			case *string:
				ts := HTMLEscapeString(*(fv.(*string)))
				tv.Set(reflect.ValueOf(&ts))
			}
		}
	}
}

// HTML escaping.

var (
	htmlQuot = []byte("&#34;") // shorter than "&quot;"
	htmlApos = []byte("&#39;") // shorter than "&apos;" and apos was not in HTML until HTML5
	htmlAmp  = []byte("&amp;")
	htmlLt   = []byte("&lt;")
	htmlGt   = []byte("&gt;")
	htmlNull = []byte("\uFFFD")
)

// HTMLEscape writes to w the escaped HTML equivalent of the plain text data b.
func HTMLEscape(w io.Writer, b []byte) {
	last := 0
	for i, c := range b {
		var html []byte
		switch c {
		case '\000':
			html = htmlNull
		// case '"':
		// 	html = htmlQuot
		// case '\'':
		// 	html = htmlApos
		case '&':
			html = htmlAmp
		case '<':
			html = htmlLt
		case '>':
			html = htmlGt
		default:
			continue
		}
		w.Write(b[last:i])
		w.Write(html)
		last = i + 1
	}
	w.Write(b[last:])
}

// HTMLEscapeString returns the escaped HTML equivalent of the plain text data s.
func HTMLEscapeString(s string) string {
	// Avoid allocation if we can.
	if !strings.ContainsAny(s, "'&<>\000") {
		return s
	}
	var b bytes.Buffer
	HTMLEscape(&b, []byte(s))
	return b.String()
}
