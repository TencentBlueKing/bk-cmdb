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

package actions_test

import (
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/core/cc/config"
	"configcenter/src/common/errors"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/topo_service"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"
)

var (
	testserverlock = sync.Mutex{}
	testserver     *httptest.Server
)

func CCAPITester(t *testing.T) *httptest.Server {
	testserverlock.Lock()
	defer testserverlock.Unlock()
	if testserver == nil {
		InitTest(t)
		app, err := ccapi.NewCCAPIServer(&config.CCAPIConfig{ExConfig: findConf()})
		require.NoError(t, err)
		app.InitHttpServ()
		testserver = httptest.NewServer(app.HttpServ().GetWebContainer())
	}
	return testserver
}

func findConf() string {
	cur, _ := filepath.Abs(".")
	for len(cur) > 1 {
		if tmp := cur + "/api.conf"; util.FileExists(tmp) {
			return tmp
		} else if tmp = cur + "/conf/api.conf"; util.FileExists(tmp) {
			return tmp
		}
		cur = filepath.Dir(cur)
	}
	return ""
}

func InitTest(t *testing.T) {
	a := api.GetAPIResource()
	require.NotNil(t, a)

	if a.Error == nil {
		errif, err := errors.New("gopath/src/configcenter/src/error_conf")
		require.NoError(t, err)
		a.Error = errif
	}
}
