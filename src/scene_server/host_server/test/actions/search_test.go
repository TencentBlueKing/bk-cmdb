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
	"bytes"
	"configcenter/src/common"
	hostParse "configcenter/src/common/paraparse"

	_ "configcenter/src/scene_server/topo_server/actions/inst"   // import inst
	_ "configcenter/src/scene_server/topo_server/actions/object" // import object actions
	_ "configcenter/src/scene_server/topo_server/actions/openapi"
	_ "configcenter/src/scene_server/topo_server/actions/privilege"
	_ "configcenter/src/scene_server/topo_server/logics/object"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func testHostSearch(t *testing.T) []string {
	server := CCAPITester(t)
	condiction := hostParse.HostCommonSearch{
		Condition: []hostParse.SearchCondition{
			hostParse.SearchCondition{
				ObjectID: "Host",
				Condition: []interface{}{
					map[string]interface{}{
						"field":    "InnerIP",
						"operator": common.BKDBEQ,
						"value":    "127.0.0.1",
					},
				},
			},
		},
	}
	reqbody, err := json.Marshal(condiction)
	r, err := http.NewRequest("POST", server.URL+"/host/v1/search", bytes.NewBuffer(reqbody))
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	respbody, err := ioutil.ReadAll(resp.Body)
	t.Logf("search result: %s", respbody)
	require.NoError(t, err)
	body := gjson.ParseBytes(respbody)
	require.Equal(t, "true", body.Get("result").String(), "response:[%s]", respbody)
	ids := body.Get("data.info.#.Host.HostID").Array()
	nids := []string{}
	for _, id := range ids {
		nids = append(nids, id.String())
	}
	return nids
}
