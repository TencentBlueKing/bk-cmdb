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
	"context"
	"testing"
	"bytes"
	"reflect"

	"configcenter/src/common/confregdiscover"
	"configcenter/src/common/errors"
)

var cc *ConfCenter

func init() {
	cc = &ConfCenter{
		ctx:          []byte("test"),
		confRegDiscv: confregdiscover.NewConfRegDiscover("127.0.0.1:2181"),
		rootCtx:      context.Background(),
	}
}

func TestGetConfigureCxt(t *testing.T) {
	cxt := cc.GetConfigureCxt()
	if "test" != string(cxt) {
		t.Error("configure cxt is not equal")
	}

}

func TestConfCenter_GetConfigureCxt(t *testing.T) {
	cc := NewConfCenter("127.0.0.1")

	cxtFake := []byte("fake cxt")
	cc.dealConfChangeEvent(cxtFake)
	if cxt := cc.GetConfigureCxt(); !bytes.Equal(cxt, cxtFake) {
		t.Errorf("context not as expected: %v", cxt)
	}
}

func TestConfCenter_GetLanguageCxt(t *testing.T) {
	cc := NewConfCenter("127.0.0.1")
	langFake := []byte("{\"fake\":{\"foo\":\"bar\"}}")
	errorCodeFake := map[string]errors.ErrorCode{
		"fake": {"foo": "bar"},
	}
	cc.dealLanguageEvent(langFake)
	if lang := cc.GetLanguageCxt(); !reflect.DeepEqual(lang, errorCodeFake) {
		t.Errorf("language not as expected: %v", lang)
	}
}
