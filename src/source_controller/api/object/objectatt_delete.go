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

// DeleteMetaObject, objID 如果为0， 则必须要设置obj，如果设置了 objID,那么会自动忽略obj 的值
func (cli *Client) DeleteMetaObjectAtt(forward *ForwardParam, objAttID int, data []byte) error {

	if 0 >= objAttID {
		if len(data) == 0 {
			return Err_Not_Set_Input
		}
	}

	rst, err := cli.base.HttpCli.DELETE(fmt.Sprintf("%s/object/v1/meta/objectatt/%d", cli.address, objAttID), forward.Header, data)

	if nil != err {
		return Err_Request_Object
	}
	blog.Debug("the result is %s, the input is %s", string(rst), string(data))
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
