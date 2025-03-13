/*
Tencent is pleased to support the open source community by making Blueking Container Service available.
Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
Licensed under the MIT License (the "License"); you may not use this file except
in compliance with the License. You may obtain a copy of the License at
http://opensource.org/licenses/MIT
Unless required by applicable law or agreed to in writing, software distributed under
the License is distributed on an "as IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
either express or implied. See the License for the specific language governing permissions and
limitations under the License.
*/

package mainline

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/http/rest"
)

const (
	topologyKey       = common.BKCacheKeyV3Prefix + "biz:custom:topology"
	detailTTLDuration = 180 * time.Minute
)

func genTopologyKey(kit *rest.Kit) string {
	return topologyKey + ":" + kit.TenantID
}

type mainlineAssociation struct {
	AssociateTo string `json:"bk_asst_obj_id" bson:"bk_asst_obj_id"`
	ObjectID    string `json:"bk_obj_id" bson:"bk_obj_id"`
}
