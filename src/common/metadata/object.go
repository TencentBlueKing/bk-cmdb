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

package metadata

import (
	"time"

	types "configcenter/src/common/mapstr"
)

const (
	ModelFieldObjCls      = "bk_classification_id"
	ModelFieldObjIcon     = "bk_obj_icon"
	ModelFieldObjectID    = "bk_obj_id"
	ModelFieldObjectName  = "bk_obj_name"
	ModelFieldIsPre       = "ispre"
	ModelFieldIsPaused    = "bk_ispaused"
	ModelFieldPosition    = "position"
	ModelFieldOwnerID     = "bk_supplier_account"
	ModelFieldDescription = "description"
	ModelFieldCreator     = "creator"
	ModelFieldModifier    = "modifier"
	ModelFieldCreateTime  = "create_time"
	ModelFieldLastTime    = "last_time"
)

// Object object metadata definition
type Object struct {
	ID          int        `field:"id"`
	ObjCls      string     `field:"bk_classification_id"`
	ObjIcon     string     `field:"bk_obj_icon"`
	ObjectID    string     `field:"bk_obj_id"`
	ObjectName  string     `field:"bk_obj_name"`
	IsPre       bool       `field:"ispre"`
	IsPaused    bool       `field:"bk_ispaused"`
	Position    string     `field:"position"`
	OwnerID     string     `field:"bk_supplier_account"`
	Description string     `field:"description"`
	Creator     string     `field:"creator"`
	Modifier    string     `field:"modifier"`
	CreateTime  *time.Time `field:"create_time"`
	LastTime    *time.Time `field:"last_time"`
}

// Parse load the data from mapstr object into object instance
func (cli *Object) Parse(data types.MapStr) (*Object, error) {

	err := SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}

// ToMapStr to mapstr
func (cli *Object) ToMapStr() types.MapStr {
	return SetValueToMapStrByTags(cli)
}
