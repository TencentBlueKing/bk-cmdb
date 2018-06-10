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

package util

import (
	"errors"
	"reflect"
)

// ParseInterface parse interface use struct
type ParseInterface struct {
	data interface{}
	err  error
}

// NewParseInterface  return a struct
func NewParseInterface(data interface{}) *ParseInterface {
	return &ParseInterface{
		data: data,
		err:  nil,
	}
}

// Get get key from interface to ParseInterface
func (p *ParseInterface) Get(key string) *ParseInterface {
	err := p.err
	if nil != err {
		return &ParseInterface{
			data: nil,
			err:  err,
		}
	}
	if nil == p.data {
		if nil == err {
			err = errors.New(key + " not found")
		}
		return &ParseInterface{
			data: nil,
			err:  err,
		}
	}

	mapData, ok := p.data.(map[string]interface{})
	if false == ok {
		return &ParseInterface{
			data: nil,
			err:  errors.New(key + " data not  map"),
		}
	}
	keyData, ok := mapData[key]
	if false == ok {
		return &ParseInterface{
			data: nil,
			err:  errors.New(key + " not found"),
		}
	}

	return &ParseInterface{
		data: keyData,
		err:  nil,
	}

}

// Interface  return interface val
func (p *ParseInterface) Interface() (interface{}, error) {
	return p.data, p.err
}

// String return string val
func (p *ParseInterface) String() (string, error) {
	if nil != p.err {
		return "", p.err
	}

	val, ok := p.data.(string)
	if true == ok {
		return val, nil
	}
	return "", errors.New("data type not string, is  " + reflect.TypeOf(p.data).String())

}

// ArrayInterface  return interface val
func (p *ParseInterface) ArrayInterface() ([]interface{}, error) {
	if nil != p.err {
		return nil, p.err
	}

	val, ok := p.data.([]interface{})
	if true == ok {
		return val, nil
	}
	return nil, errors.New("data type not string, is  " + reflect.TypeOf(p.data).String())

}
