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

package y3_10_202203031512

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

// addHostAgentIDAndIPv6Attr add host agent id and ipv6 attributes
func addHostAgentIDAndIPv6Attr(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	// NOCC:tosa/linelength(IPv6正则表达式的长度超过限制，忽略该检查)
	const ipv6Regex = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	multiIpv6Regex := fmt.Sprintf(`^%s(,%s)*$`, ipv6Regex, ipv6Regex)

	agentIDAndIPv6Attrs := []attribute{
		{
			PropertyID:   "bk_host_innerip_v6",
			PropertyName: "内网IPv6",
			IsEditable:   false,
			PropertyType: common.FieldTypeSingleChar,
			Option:       multiIpv6Regex,
		},
		{
			PropertyID:   "bk_host_outerip_v6",
			PropertyName: "外网IPv6",
			IsEditable:   true,
			PropertyType: common.FieldTypeSingleChar,
			Option:       multiIpv6Regex,
		},
		{
			PropertyID:   "bk_agent_id",
			PropertyName: "GSE Agent ID",
			IsEditable:   false,
			PropertyType: common.FieldTypeSingleChar,
			Option:       "^[a-zA-Z0-9]{0,64}$",
		},
	}

	now := time.Now()
	attrIDs := make([]string, 0)
	for index, attr := range agentIDAndIPv6Attrs {
		agentIDAndIPv6Attrs[index].OwnerID = conf.OwnerID
		agentIDAndIPv6Attrs[index].ObjectID = common.BKInnerObjIDHost
		agentIDAndIPv6Attrs[index].PropertyGroup = "default"
		agentIDAndIPv6Attrs[index].IsPre = true
		agentIDAndIPv6Attrs[index].Creator = common.CCSystemOperatorUserName
		agentIDAndIPv6Attrs[index].CreateTime = &now
		agentIDAndIPv6Attrs[index].LastTime = &now

		attrIDs = append(attrIDs, attr.PropertyID)
	}

	return insertHostAgentIDAndIPv6Attr(ctx, db, attrIDs, agentIDAndIPv6Attrs)
}

// insertHostAgentIDAndIPv6Attr insert not exist host agent id and ipv6 attributes in db
func insertHostAgentIDAndIPv6Attr(ctx context.Context, db dal.RDB, attrIDs []string, attrArr []attribute) error {
	// check if the attributes to add are already exist
	attrFilter := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: map[string]interface{}{common.BKDBIN: attrIDs},
	}

	existAttrs := make([]attribute, 0)
	err := db.Table(common.BKTableNameObjAttDes).Find(attrFilter).Fields(common.BKPropertyIDField).All(ctx, &existAttrs)
	if err != nil {
		blog.Errorf("check if to insert host attribute exists failed, filter: %v, err: %v", attrFilter, err)
		return err
	}

	existAttrMap := make(map[string]struct{})
	for _, attr := range existAttrs {
		existAttrMap[attr.PropertyID] = struct{}{}
	}

	toInsertAttrs := make([]attribute, 0)
	for _, attr := range attrArr {
		if _, exists := existAttrMap[attr.PropertyID]; !exists {
			toInsertAttrs = append(toInsertAttrs, attr)
		}
	}

	attrLen := len(toInsertAttrs)
	if attrLen == 0 {
		return nil
	}

	// add attributes that are not exist, generate new id and index for them
	newAttrIDs, err := db.NextSequences(ctx, common.BKTableNameObjAttDes, attrLen)
	if err != nil {
		blog.Errorf("get new attributes ids failed, err: %v", err)
		return err
	}

	attrIdxFilter := map[string]interface{}{
		common.BKObjIDField: common.BKInnerObjIDHost,
	}
	sort := common.BKPropertyIndexField + ":-1"
	maxIdxAttr := new(attribute)
	if err := db.Table(common.BKTableNameObjAttDes).Find(attrIdxFilter).Sort(sort).One(ctx, maxIdxAttr); err != nil {
		blog.Errorf("get max host attribute index failed, filter: %v, err: %v", attrIdxFilter, err)
		return err
	}

	for index, _ := range toInsertAttrs {
		toInsertAttrs[index].ID = int64(newAttrIDs[index])
		toInsertAttrs[index].PropertyIndex = maxIdxAttr.PropertyIndex + int64(index) + 1
	}

	if err := db.Table(common.BKTableNameObjAttDes).Insert(ctx, toInsertAttrs); err != nil {
		blog.Errorf("insert host attributes(%#v) failed, err: %v", toInsertAttrs, err)
		return err
	}

	return nil
}

// addHostAgentIDAndIPv6Unique add host agent id and ipv6 attributes related unique rule
func addHostAgentIDAndIPv6Unique(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	// find inner ipv6 and cloud id attribute ids
	attrCond := map[string]interface{}{
		common.BKObjIDField: common.BKInnerObjIDHost,
		common.BKPropertyIDField: map[string]interface{}{common.BKDBIN: []string{"bk_host_innerip_v6",
			common.BKCloudIDField}},
	}
	attrArr := make([]attribute, 0)
	err := db.Table(common.BKTableNameObjAttDes).Find(attrCond).All(ctx, &attrArr)
	if err != nil {
		blog.Errorf("get host unique fields failed, filter: %v, err: %v", attrCond, err)
		return err
	}

	attrIDMap := make(map[string]uint64)
	for _, attr := range attrArr {
		attrIDMap[attr.PropertyID] = uint64(attr.ID)
	}

	unique := objectUnique{
		Keys: []uniqueKey{
			{Kind: "property", ID: attrIDMap["bk_host_innerip_v6"]},
			{Kind: "property", ID: attrIDMap[common.BKCloudIDField]},
		},
		ObjID:    common.BKInnerObjIDHost,
		IsPre:    true,
		OwnerID:  conf.OwnerID,
		LastTime: time.Now(),
	}

	keyMap := make(map[uint64]string)
	for _, key := range unique.Keys {
		keyMap[key.ID] = key.Kind
	}

	// get already exist host uniques for comparison
	uniqueCond := map[string]interface{}{common.BKObjIDField: common.BKInnerObjIDHost}
	existUniques := make([]objectUnique, 0)
	err = db.Table(common.BKTableNameObjUnique).Find(uniqueCond).All(ctx, &existUniques)
	if err != nil {
		blog.Errorf("get host unique failed, filter: %v, err: %v", uniqueCond, err)
		return err
	}

	// check if inner ipv6 + cloud id unique exists, skip the creation if it already exists
	for _, existUnique := range existUniques {
		if len(unique.Keys) != len(existUnique.Keys) {
			continue
		}

		isEqual := true
		for _, key := range existUnique.Keys {
			if key.Kind != keyMap[key.ID] {
				isEqual = false
				break
			}
		}

		if isEqual {
			return nil
		}
	}

	// insert ipv6 + cloud id unique when it is not exist
	newUniqueID, err := db.NextSequence(ctx, common.BKTableNameObjUnique)
	if err != nil {
		blog.Errorf("get new unique ids failed, err: %v", err)
		return err
	}

	unique.ID = newUniqueID

	if err := db.Table(common.BKTableNameObjUnique).Insert(ctx, unique); err != nil {
		blog.Errorf("insert host unique(%#v) failed, err: %v", unique, err)
		return err
	}
	return nil
}

// addHostAgentIDAndIPv6Index add host agent id and ipv6 attributes db index
func addHostAgentIDAndIPv6Index(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	agentIDAndIPv6Indexes := []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bkHostInnerIPv6_bkCloudID",
			Keys: bson.D{{
				"bk_host_innerip_v6", 1},
				{common.BKCloudIDField, 1},
			},
			Unique:     true,
			Background: true,
			PartialFilterExpression: map[string]interface{}{
				"bk_host_innerip_v6":  map[string]string{common.BKDBType: "string"},
				common.BKCloudIDField: map[string]string{common.BKDBType: "number"},
			},
		},
	}

	existIndexArr, err := db.Table(common.BKTableNameBaseHost).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist index for host table failed, err: %v", err)
		return err
	}

	existIdxMap := make(map[string]struct{})
	for _, index := range existIndexArr {
		existIdxMap[index.Name] = struct{}{}
	}

	for _, index := range agentIDAndIPv6Indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}

		err = db.Table(common.BKTableNameBaseHost).CreateIndex(ctx, index)
		if err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("create index(%#v) failed, err: %v", index, err)
			return err
		}
	}

	return nil
}

type attribute struct {
	BizID         int64       `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id" mapstructure:"bk_biz_id"`
	ID            int64       `field:"id" json:"id" bson:"id" mapstructure:"id"`
	OwnerID       string      `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account" mapstructure:"bk_supplier_account"`
	ObjectID      string      `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id" mapstructure:"bk_obj_id"`
	PropertyID    string      `field:"bk_property_id" json:"bk_property_id" bson:"bk_property_id" mapstructure:"bk_property_id"`
	PropertyName  string      `field:"bk_property_name" json:"bk_property_name" bson:"bk_property_name" mapstructure:"bk_property_name"`
	PropertyGroup string      `field:"bk_property_group" json:"bk_property_group" bson:"bk_property_group" mapstructure:"bk_property_group"`
	PropertyIndex int64       `field:"bk_property_index" json:"bk_property_index" bson:"bk_property_index" mapstructure:"bk_property_index"`
	Unit          string      `field:"unit" json:"unit" bson:"unit" mapstructure:"unit"`
	Placeholder   string      `field:"placeholder" json:"placeholder" bson:"placeholder" mapstructure:"placeholder"`
	IsEditable    bool        `field:"editable" json:"editable" bson:"editable" mapstructure:"editable"`
	IsPre         bool        `field:"ispre" json:"ispre" bson:"ispre" mapstructure:"ispre"`
	IsRequired    bool        `field:"isrequired" json:"isrequired" bson:"isrequired" mapstructure:"ispre"`
	IsReadOnly    bool        `field:"isreadonly" json:"isreadonly" bson:"isreadonly" mapstructure:"isreadonly"`
	IsOnly        bool        `field:"isonly" json:"isonly" bson:"isonly" mapstructure:"isonly"`
	IsSystem      bool        `field:"bk_issystem" json:"bk_issystem" bson:"bk_issystem" mapstructure:"bk_issystem"`
	IsAPI         bool        `field:"bk_isapi" json:"bk_isapi" bson:"bk_isapi" mapstructure:"bk_isapi"`
	PropertyType  string      `field:"bk_property_type" json:"bk_property_type" bson:"bk_property_type" mapstructure:"bk_property_type"`
	Option        interface{} `field:"option" json:"option" bson:"option" mapstructure:"option"`
	Description   string      `field:"description" json:"description" bson:"description" mapstructure:"description"`
	Creator       string      `field:"creator" json:"creator" bson:"creator" mapstructure:"creator"`
	CreateTime    *time.Time  `json:"create_time" bson:"create_time" mapstructure:"create_time"`
	LastTime      *time.Time  `json:"last_time" bson:"last_time" mapstructure:"last_time"`
}

type objectUnique struct {
	ID       uint64      `json:"id" bson:"id"`
	ObjID    string      `json:"bk_obj_id" bson:"bk_obj_id"`
	Keys     []uniqueKey `json:"keys" bson:"keys"`
	IsPre    bool        `json:"ispre" bson:"ispre"`
	OwnerID  string      `json:"bk_supplier_account" bson:"bk_supplier_account"`
	LastTime time.Time   `json:"last_time" bson:"last_time"`
}

type uniqueKey struct {
	Kind string `json:"key_kind" bson:"key_kind"`
	ID   uint64 `json:"key_id" bson:"key_id"`
}
