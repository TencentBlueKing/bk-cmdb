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
	"bytes"
	"configcenter/src/common/errors"
	"reflect"
	"testing"
)

func TestNewConfCenter(t *testing.T) {
	if s := NewConfCenter("127.0.0.1"); s == nil {
		t.Errorf("conf center is nil")
	}
}

func TestConfCenter_Start(t *testing.T) {

}

func TestConfCenter_Stop(t *testing.T) {
}

func TestConfCenter_GetConfigureCxt(t *testing.T) {
	s := NewConfCenter("127.0.0.1")

	cxtFake := []byte("fake cxt")
	s.dealConfChangeEvent(cxtFake)
	if cxt := s.GetConfigureCxt(); !bytes.Equal(cxt, cxtFake) {
		t.Errorf("context not as expected: %v", cxt)
	}
}

func TestConfCenter_GetErrorCxt(t *testing.T) {
	s := NewConfCenter("127.0.0.1")
	langFake := []byte("{\"fake\":{\"foo\":\"bar\"}}")
	errorCodeFake := map[string]errors.ErrorCode{
		"fake": {"foo": "bar"},
	}
	s.dealErrorResEvent(langFake)
	if lang := s.GetErrorCxt(); !reflect.DeepEqual(lang, errorCodeFake) {
		t.Errorf("language not as expected: %v", lang)
	}
}
