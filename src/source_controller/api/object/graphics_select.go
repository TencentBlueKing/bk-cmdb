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

package object

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
)

// SearchTopoGraphics search topo graphics
func (cli *Client) SearchTopoGraphics(forward *ForwardParam, params *TopoGraphics) ([]TopoGraphics, error) {

	out, err := json.Marshal(params)
	if err != nil {
		blog.Errorf("SearchTopoGraphics marshal error %v", err)
		return nil, err
	}

	rst, err := cli.base.HttpCli.POST(fmt.Sprintf("%s/object/v1/topographics/search", cli.address), forward.Header, out)

	if nil != err {
		blog.Error("request failed, error:%v", err)
		return nil, Err_Request_Object
	}

	rstRes := gjson.ParseBytes(rst)

	if rstRes.Get(common.HTTPBKAPIErrorCode).Int() != common.CCSuccess {
		return nil, fmt.Errorf("%v", rstRes.Get(common.HTTPBKAPIErrorMessage))
	}

	var ret []TopoGraphics
	if jserr := json.Unmarshal([]byte(rstRes.Get("data").String()), &ret); nil != jserr {
		blog.Error("can not unmarshal the result , error information is %v, %s", jserr, rstRes.String())
		return nil, jserr
	}

	return ret, nil
}
