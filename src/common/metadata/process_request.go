/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */
package metadata

type CreateServiceCategoryOption struct {
	Metadata *Metadata `field:"metadata" json:"metadata" bson:"metadata"`
	BizID    int64     `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	Name     string    `field:"name" json:"name,omitempty" bson:"name"`
	ParentID int64     `field:"bk_parent_id" json:"bk_parent_id,omitempty" bson:"bk_parent_id"`
}

type CreateServiceTemplateOption struct {
	Metadata          *Metadata `field:"metadata" json:"metadata" bson:"metadata"`
	BizID             int64     `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	Name              string    `field:"name" json:"name,omitempty" bson:"name"`
	ServiceCategoryID int64     `field:"service_category_id" json:"service_category_id,omitempty" bson:"service_category_id"`
}

type UpdateServiceTemplateOption struct {
	Metadata          *Metadata `field:"metadata" json:"metadata" bson:"metadata"`
	BizID             int64     `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	ID                int64     `field:"id" json:"id,omitempty" bson:"id"`
	Name              string    `field:"name" json:"name,omitempty" bson:"name"`
	ServiceCategoryID int64     `field:"service_category_id" json:"service_category_id,omitempty" bson:"service_category_id"`
}

type RemoveFromModuleHost struct {
	MoveToIdle bool    `field:"move_to_idle" json:"move_to_idle"`
	HostID     int64   `field:"bk_host_id" json:"bk_host_id"`
	Modules    []int64 `field:"bk_module_ids" json:"bk_module_ids"`
}

type ServiceInstanceDeletePreview struct {
	ToMoveModuleHosts []RemoveFromModuleHost `field:"to_move_module_hosts" json:"to_move_module_hosts"`
}
