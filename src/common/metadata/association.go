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

import types "configcenter/src/common/mapstr"

const (
	// AssociationFieldObjectID the association data field definition
	AssociationFieldObjectID = "bk_obj_id"
	// AssociationFieldObjectAttributeID the association data field definition
	AssociationFieldObjectAttributeID = "bk_object_att_id"
	// AssociationFieldSupplierAccount the association data field definition
	AssociationFieldSupplierAccount = "bk_supplier_account"
	// AssociationFieldAssociationForward the association data field definition
	AssociationFieldAssociationForward = "bk_asst_forward"
	// AssociationFieldAssociationObjectID the association data field definition
	AssociationFieldAssociationObjectID = "bk_asst_obj_id"
	// AssociationFieldAssociationName the association data field definition
	AssociationFieldAssociationName = "bk_asst_name"
)

// Association define object association struct
type Association struct {
	ID          int    `field:"id"`
	ObjectID    string `field:"bk_obj_id"`
	ObjectAttID string `field:"bk_object_att_id"`
	OwnerID     string `field:"bk_supplier_account"`
	AsstForward string `field:"bk_asst_forward"`
	AsstObjID   string `field:"bk_asst_obj_id"`
	AsstName    string `field:"bk_asst_name"`
}

// Parse load the data from mapstr attribute into attribute instance
func (cli *Association) Parse(data types.MapStr) (*Association, error) {

	err := SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}

// ToMapStr to mapstr
func (cli *Association) ToMapStr() types.MapStr {
	return SetValueToMapStrByTags(cli)
}
