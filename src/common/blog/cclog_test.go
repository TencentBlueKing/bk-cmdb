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

package blog

import (
	"net/http"
	"testing"

	"configcenter/src/common"
)

func getHeader() http.Header {
	header := make(http.Header, 0)
	header.Set(common.BKHTTPCCRequestID, "xxx-log-id")
	return header
}

func TestErrorf(t *testing.T) {
	CCLogHeader(nil).Error("int:%d,str:%s, array:%v", 1, "str", []int{1, 2})
	CCLogHeader(nil).Errorf("int:%d,str:%s, array:%v", 2, "strf", []string{"errorf", "errorf"})

	CCLogHeader(getHeader()).Error("int:%d,str:%s, array:%v", 1, "str", []int{1, 2})
	CCLogHeader(getHeader()).Errorf("int:%d,str:%s, array:%v", 2, "strf", []string{"errorf", "errorf"})
}

func TestWarn(t *testing.T) {
	CCLogHeader(nil).Warn("int:%d,str:%s, array:%v", 1, "str", []int{1, 2})
	CCLogHeader(nil).Warnf("int:%d,str:%s, array:%v", 2, "strf", []string{"errorf", "errorf"})

	CCLogHeader(getHeader()).Warn("int:%d,str:%s, array:%v", 1, "str", []int{1, 2})
	CCLogHeader(getHeader()).Warnf("int:%d,str:%s, array:%v", 2, "strf", []string{"errorf", "errorf"})
}

func TestDebug(t *testing.T) {
	CCLogHeader(nil).Debug("int:%d,str:%s, array:%v", 1, "str", []int{1, 2})

	CCLogHeader(getHeader()).Debug("int:%d,str:%s, array:%v", 1, "str", []int{1, 2})
}

func TestInfo(t *testing.T) {
	CCLogHeader(nil).Info("int:%d,str:%s, array:%v", 1, "str", []int{1, 2})
	CCLogHeader(nil).Infof("int:%d,str:%s, array:%v", 2, "strf", []string{"errorf", "errorf"})

	CCLogHeader(getHeader()).Info("int:%d,str:%s, array:%v", 1, "str", []int{1, 2})
	CCLogHeader(getHeader()).Infof("int:%d,str:%s, array:%v", 2, "strf", []string{"errorf", "errorf"})
}

func TestInfoJSON(t *testing.T) {
	CCLogHeader(nil).InfoJSON("int:%d,str:%s, array:%v", 1, "str", []int{1, 2})

	CCLogHeader(getHeader()).InfoJSON("int:%d,str:%s, array:%v", 1, "str", []int{1, 2})
}
