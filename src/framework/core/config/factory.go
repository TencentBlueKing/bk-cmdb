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

package config

import (
	"bufio"
	"bytes"
	"configcenter/src/common/blog"
	"io"
	"os"
	"strings"
)

const separator = "."

// ParseFromBytes parse config from data
func ParseFromBytes(data []byte) (Config, error) {
	buf := bytes.NewBuffer(data)
	return Parse(buf)
}

// ParseFromFile parse config from file
func ParseFromFile(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		blog.Infof("path: %s,config file not exits;", path)
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

// Parse parse config from reader
func Parse(rd io.Reader) (Config, error) {
	strcet := ""
	r := bufio.NewReader(rd)
	c := Config{}
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		s := strings.TrimSpace(string(b))
		//fmt.Println(s)
		if strings.Index(s, "#") == 0 {
			continue
		}

		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			strcet = strings.TrimSpace(s[n1+1 : n2])
			continue
		}

		if len(strcet) == 0 {
			continue
		}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		frist := strings.TrimSpace(s[:index])
		if len(frist) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}

		key := strcet + separator + frist
		c[key] = strings.TrimSpace(second)
	}
	return c, nil
}
