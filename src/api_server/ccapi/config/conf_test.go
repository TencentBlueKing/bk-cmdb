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
	"testing"
	"context"
	"encoding/json"
	"configcenter/src/common/confregdiscover"
	"github.com/stretchr/testify/assert"
	"configcenter/src/common/errors"
)

var conf *ConfCenter
var errCode map[string]errors.ErrorCode

const (
	TestConfEvent = "test_conf_event"
)

func init(){
	conf = &ConfCenter{
		confRegDiscv:confregdiscover.NewConfRegDiscover("127.0.0.1:2181"),
	}
	conf.rootCtx, conf.cancel = context.WithCancel(context.Background())
	errCode = map[string]errors.ErrorCode{
		"test_err":errors.ErrorCode{},
	}
}

func TestDealConfChangeEvent(t *testing.T){
	event := []byte(TestConfEvent)
	err := conf.dealConfChangeEvent(event)
	assert.Nil(t,err)
}

func TestGetConfigureCxt(t *testing.T){
	cxt := conf.GetConfigureCxt()
	assert.Equal(t,[]byte(TestConfEvent),cxt)
}

func TestDealLanguageEvent(t *testing.T){
	by,_ := json.Marshal(errCode)
	err := conf.dealLanguageEvent(by)
	assert.Nil(t,err)
}

func TestGetLanguageCxt(t *testing.T){
	errC := conf.GetLanguageCxt()
	byC,_ := json.Marshal(errC)
	byE,_ := json.Marshal(errCode)
	assert.Equal(t,byE,byC)
}