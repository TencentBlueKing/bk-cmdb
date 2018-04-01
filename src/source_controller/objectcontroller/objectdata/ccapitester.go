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
 
package ccapi

import (
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/core/cc/config"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testServerLock = &sync.Mutex{}
	testServer     *httptest.Server
)

func CCAPITester(t *testing.T) *httptest.Server {
	testServerLock.Lock()
	defer testServerLock.Unlock()
	if testServer == nil {
		app, err := NewCCAPIServer(&config.CCAPIConfig{ExConfig: findConf()})
		require.NoError(t, err)
		a := api.NewAPIResource()
		config, _ := a.ParseConfig()
		app.InitHttpServ(config)
		testServer = httptest.NewServer(app.HttpServ().GetWebContainer())
	}
	return testServer
}

func findConf() string {
	cur, _ := filepath.Abs(".")
	for len(cur) > 1 {
		if tmp := cur + "/api.conf"; exists(tmp) {
			return tmp
		} else if tmp = cur + "/conf/api.conf"; exists(tmp) {
			return tmp
		}
		cur = filepath.Dir(cur)
	}
	return ""
}

func exists(name string) bool {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
