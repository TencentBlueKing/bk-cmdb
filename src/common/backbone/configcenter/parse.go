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

package configcenter

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"configcenter/src/common/errors"
	"configcenter/src/common/language"
)

const middle = "."

func LoadConfigFromLocalFile(confPath string, handler *CCHandler) error {
	fileConf, err := ParseConfigWithFile(confPath)
	if err != nil {
		return err
	}

	lang, ok := fileConf.ConfigMap["language.res"]
	if !ok {
		return fmt.Errorf("load config from file[%s], but can not found language config", confPath)
	}

	langC, err := language.LoadLanguageResourceFromDir(lang)
	if err != nil {
		return fmt.Errorf("load config from file[%s], but load language failed, err: %v", confPath, err)
	}

	errCode, ok := fileConf.ConfigMap["errors.res"]
	if !ok {
		return fmt.Errorf("load config from file[%s], but can not found error code config", confPath)
	}

	errC, err := errors.LoadErrorResourceFromDir(errCode)
	if err != nil {
		return fmt.Errorf("load config from file[%s], but load error code failed, err: %v", confPath, err)
	}

	if len(fileConf.ConfigMap) == 0 {
		return fmt.Errorf("load config from file[%s], but can not found process config", confPath)
	}

	handler.OnProcessUpdate(ProcessConfig{}, ProcessConfig{ConfigMap: fileConf.ConfigMap})
	handler.OnLanguageUpdate(nil, langC)
	handler.OnErrorUpdate(nil, errC)

	return nil
}

type ProcessConfig struct {
	ConfigMap map[string]string
}

func ParseConfigWithFile(filePath string) (*ProcessConfig, error) {
	c := config{
		configmap: map[string]string{},
	}
	if err := c.parseWithFile(filePath); err != nil {
		return nil, err
	}

	return &ProcessConfig{ConfigMap: c.configmap}, nil
}

func ParseConfigWithData(data []byte) (*ProcessConfig, error) {
	c := config{
		configmap: map[string]string{},
	}
	if err := c.parseWithData(data); err != nil {
		return nil, err
	}

	return &ProcessConfig{ConfigMap: c.configmap}, nil
}

type config struct {
	configmap map[string]string
	strcet    string
}

func (p *config) parse(input io.Reader) error {
	r := bufio.NewReader(input)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		s := strings.TrimSpace(string(b))
		if strings.Index(s, "#") == 0 {
			continue
		}

		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			p.strcet = strings.TrimSpace(s[n1+1 : n2])
			continue
		}

		if len(p.strcet) == 0 {
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

		key := p.strcet + middle + frist
		p.configmap[key] = strings.TrimSpace(second)
	}
	return nil
}

func (c *config) parseWithFile(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := c.parse(bufio.NewReader(f)); err != nil {
		return fmt.Errorf("parse config failed, err: %v", err)
	}
	return nil
}

func (c *config) parseWithData(data []byte) error {
	return c.parse(bytes.NewReader(data))
}
