/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"

	"github.com/graphql-go/graphql"
)

type Host struct {
	HostID   int64  `json:"bk_host_id" bson:"bk_host_id"`           // 主机ID(host_id)								数字
	HostName string `json:"bk_host_name" bson:"bk_host_name"`       // 主机名称
	InnerIP  string `json:"bk_host_innerip" bson:"bk_host_innerip"` // 内网IP
	OuterIP  string `json:"bk_host_outerip" bson:"bk_host_outerip"` // 外网IP
}

var HostType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Host",
		Fields: graphql.Fields{
			common.BKHostIDField: &graphql.Field{
				Type: graphql.Int,
			},
			common.BKHostNameField: &graphql.Field{
				Type: graphql.String,
			},
			common.BKHostInnerIPField: &graphql.Field{
				Type: graphql.String,
			},
			common.BKHostOuterIPField: &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var ObjectConfig = graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"host": &graphql.Field{
			Type: HostType,
			Args: graphql.FieldConfigArgument{
				common.BKHostIDField: &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				hostIDIf, isOK := p.Args[common.BKHostIDField]
				if isOK == false {
					return nil, fmt.Errorf("key: %s not found", common.BKHostIDField)
				}
				hostID, err := util.GetInt64ByInterface(hostIDIf)
				if err != nil {
					return nil, fmt.Errorf("value %+v of key %s not found", hostIDIf, common.BKHostIDField)
				}

				hostFilter := map[string]interface{}{
					common.BKHostIDField: hostID,
				}
				host := Host{}
				if err := DbProxy.Table(common.BKTableNameBaseHost).Find(hostFilter).One(context.Background(), &host); err != nil {
					return nil, fmt.Errorf("select host failed, err: %+v", err)
				}
				return host, nil
			},
		},
	},
}

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: graphql.NewObject(ObjectConfig),
	},
)

var DbProxy dal.RDB

func (s *coreService) SearchGraphql(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	DbProxy = s.db
	queryIf, ok := data["query"]
	if ok == false {
		return nil, fmt.Errorf("query not found, data: %+v", data)
	}
	query := util.GetStrByInterface(queryIf)
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		blog.Infof("wrong result, unexpected errors: %v", result.Errors)
	}
	return result, nil
}
