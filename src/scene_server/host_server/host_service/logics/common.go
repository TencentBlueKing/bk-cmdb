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

package logics

import (
	"encoding/json"

	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	sourceAPI "configcenter/src/source_controller/api/object"
)

// GetHttpResult get http result
func GetHttpResult(req *restful.Request, url, method string, params interface{}) (bool, string, interface{}) {
	var strParams []byte
	switch params.(type) {
	case string:
		strParams = []byte(params.(string))
	default:
		strParams, _ = json.Marshal(params)

	}
	blog.Info("get request url:%s", url)
	blog.Info("get request info  params:%v", string(strParams))
	reply, err := httpcli.ReqHttp(req, url, method, []byte(strParams))
	blog.Info("get request result:%v", string(reply))
	if err != nil {
		blog.Error("http do error, params:%s, error:%s", strParams, err.Error())
		return false, err.Error(), nil
	}

	addReply, err := simplejson.NewJson([]byte(reply))
	if err != nil {
		blog.Error("http do error, params:%s, reply:%s, error:%s", strParams, reply, err.Error())
		return false, err.Error(), nil
	}
	isSuccess, err := addReply.Get("result").Bool()
	if nil != err || !isSuccess {
		errMsg, _ := addReply.Get(common.HTTPBKAPIErrorMessage).String()
		blog.Error("http do error, url:%s, params:%s, error:%s", url, strParams, errMsg)
		return false, errMsg, addReply.Get("data").Interface()
	}
	return true, "", addReply.Get("data").Interface()
}

//GetObjectFields get object fields
func GetObjectFields(forward *sourceAPI.ForwardParam, ownerID, objID, ObjAddr, sort string) ([]sourceAPI.ObjAttDes, error) {
	data := make(map[string]interface{})
	data[common.BKOwnerIDField] = ownerID
	data[common.BKObjIDField] = objID
	data["page"] = common.KvMap{
		"start": 0,
		"limit": common.BKNoLimit,
		"sort":  sort,
	}
	info, _ := json.Marshal(data)
	client := sourceAPI.NewClient(ObjAddr)
	atts, err := client.SearchMetaObjectAtt(forward, []byte(info))
	if nil != err {
		return nil, err
	}

	for idx, a := range atts {
		if !util.IsAssocateProperty(a.PropertyType) {
			continue
		}
		// read property group
		condition := map[string]interface{}{
			"bk_object_att_id":    a.PropertyID, // tmp.PropertyGroup,
			common.BKOwnerIDField: a.OwnerID,
			"bk_obj_id":           a.ObjectID,
		}
		objasstval, jserr := json.Marshal(condition)
		if nil != jserr {
			blog.Error("mashar json failed, error information is %v", jserr)
			return nil, jserr
		}
		asstMsg, err := client.SearchMetaObjectAsst(forward, objasstval)
		if nil != err {
			return nil, err
		}
		if 0 < len(asstMsg) {
			atts[idx].AssociationID = asstMsg[0].AsstObjID // by the rules, only one id
			atts[idx].AsstForward = asstMsg[0].AsstForward // by the rules, only one id
		}
	}
	return atts, nil
}
