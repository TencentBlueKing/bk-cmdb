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
)

func (cli *Client) SelectPropertyGroup(forward *ForwardParam, data []byte) ([]ObjAttGroupDes, error) {
	rst, err := cli.base.HttpCli.POST("/meta/objectatt/group/search", forward.Header, data)
	if nil != err {
		blog.Error("request failed, error:%v", err)
		return nil, Err_Request_Object
	}

	var rstRes ObjAttGroupRsp
	if jserr := json.Unmarshal(rst, &rstRes); nil != jserr {
		blog.Error("can not unmarshal the result , error information is %v", jserr)
		return nil, jserr
	}

	if rstRes.Code != common.CCSuccess {
		return nil, fmt.Errorf("%v", rstRes.Message)
	}

	return rstRes.Data, nil
}

// SelectPropertyGroupByObjectID 查询元数据对象集合
func (cli *Client) SelectPropertyGroupByObjectID(forward *ForwardParam, ownerID, objectID string, data []byte) ([]ObjAttGroupDes, error) {

	url := fmt.Sprintf("%s/object/v1/meta/objectatt/group/property/owner/%s/object/%s", cli.address, ownerID, objectID)
	rst, err := cli.base.HttpCli.POST(url, forward.Header, data)
	if nil != err {
		blog.Error("request failed, error:%v", err)
		return nil, Err_Request_Object
	}
	blog.Debug("search property group by objectid, url(%s),result:%s", url, string(rst))
	var rstRes ObjAttGroupRsp
	if jserr := json.Unmarshal(rst, &rstRes); nil != jserr {
		blog.Error("can not unmarshal the result , error information is %v", jserr)
		return nil, jserr
	}

	if rstRes.Code != common.CCSuccess {
		return nil, fmt.Errorf("%v", rstRes.Message)
	}

	return rstRes.Data, nil
}
