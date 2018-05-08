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
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/option"
	"errors"
	"io"
	"os"
	"strings"
	"time"
)

const separator = "."

// Init init config
func Init(opt *option.Options) error {
	if nil == config {
		config = Config{}
	}
	if "" != opt.Config {
		return ParseFromFile(opt.Config)
	}
	if "" != opt.Regdiscv {
		cc := NewConfCenter(opt.Regdiscv)
		dataCh := make(chan []byte)
		errCh := make(chan error)
		go func() {
			defer cc.Stop()
			err := cc.Start()
			log.Errorf("configure center module start failed!. err:%s", err.Error())
			errCh <- err
		}()

		go func() {
			var data []byte
			for {
				data = cc.GetConfigureCxt()
				if len(data) <= 0 {
					log.Warningf("faile to get config from center, we will retry after 2 seconds")
					time.Sleep(time.Second * 2)
					continue
				}
				dataCh <- data
			}
		}()
		select {
		case data := <-dataCh:
			if err := ParseFromBytes(data); err != nil {
				return err
			}
		case err := <-errCh:
			return err
		}
	}
	return errors.New("not config source specified")
}

// ParseFromBytes parse config from data
func ParseFromBytes(data []byte) error {
	buf := bytes.NewBuffer(data)
	return Parse(buf)
}

// ParseFromFile parse config from file
func ParseFromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		log.Infof("path: %s,config file not exits;", path)
		return err
	}
	defer f.Close()
	return Parse(f)
}

// Parse parse config from reader
func Parse(rd io.Reader) error {
	strcet := ""
	r := bufio.NewReader(rd)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
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
		config[key] = strings.TrimSpace(second)
	}
	return nil
}
