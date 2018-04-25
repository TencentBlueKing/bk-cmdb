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
	"configcenter/src/common/core/cc/api"
	"encoding/json"
	"fmt"
)

// DeleteMetaObjectAttGroup delete the property group
func (cli *Client) DeleteMetaObjectAttGroup(forward *ForwardParam, id int, data []byte) error {

	if 0 >= id {
		if len(data) == 0 {
			return Err_Not_Set_Input
		}
	}

	rst, err := cli.base.HttpCli.DELETE(fmt.Sprintf("%s/object/v1/meta/objectatt/group/groupid/%d", cli.address, id), forward.Header, data)

	if nil != err {
		return Err_Request_Object
	}

	var rstRes api.APIRsp
	if jserr := json.Unmarshal(rst, &rstRes); nil != jserr {
		blog.Error("can not unmarshal the result , error information is %v", jserr)
		return jserr
	}

	if rstRes.Code != common.CCSuccess {
		return fmt.Errorf("%v", rstRes.Message)
	}

	return nil
}

// DeleteMetaObjectAttGroupProperty delete the group property
func (cli *Client) DeleteMetaObjectAttGroupProperty(forward *ForwardParam, ownerID, objectID, propertyID, groupID string) error {

	rst, err := cli.base.HttpCli.DELETE(fmt.Sprintf("%s/object/v1/meta/objectatt/group/owner/%s/object/%s/propertyids/%s/groupids/%s", cli.address, ownerID, objectID, propertyID, groupID), forward.Header, nil)

	if nil != err {
		return Err_Request_Object
	}

	var rstRes api.APIRsp
	if jserr := json.Unmarshal(rst, &rstRes); nil != jserr {
		blog.Error("can not unmarshal the result , error information is %v", jserr)
		return jserr
	}

	if rstRes.Code != common.CCSuccess {
		return fmt.Errorf("%v", rstRes.Message)
	}

	return nil
}
