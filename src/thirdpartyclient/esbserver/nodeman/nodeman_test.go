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

package nodeman

import (
	"encoding/json"
	"testing"
)

func TestDecode(t *testing.T) {
	out := []byte(`
	{
	  "message": "success",
	  "code": "OK",
	  "data": [
		{
		  "pkg_mtime": "2018-09-19 15:51:10",
		  "cpu_arc": "xx",
		  "module": "gse_plugin",
		  "project": "netdevicebeat",
		  "pkg_size": 5142819,
		  "version": "1.0.0",
		  "pkg_name": "netdevicebeat-linux-1.0.0-x86_64.tgz",
		  "location": "http://10.167.77.15/download",
		  "pkg_ctime": "2018-09-19 15:51:10",
		  "pkg_path": "/data/bkee/miniweb/download",
		  "os": "linux",
		  "id": 5,
		  "md5": "5bbcfedc7c9cb6a5fbf1a3459fd12c24"
		}
	  ],
	  "result": true,
	  "request_id": "cd68f03a62d4463bb29fd77edac87cc4"
	}
	`)
	resp := new(SearchPluginPackageResult)
	err := json.Unmarshal(out, resp)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("resp: %+v", resp)
}
