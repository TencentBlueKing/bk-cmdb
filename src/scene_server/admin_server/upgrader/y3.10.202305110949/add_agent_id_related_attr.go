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

package y3_10_202305110949

import (
	"context"
	"errors"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

// changeHostIpv4RequireAttr change the value of the isrequired attribute of the host bk_host_innerip to false.
func changeHostIpv4RequireAttr(ctx context.Context, db dal.RDB) error {

	// 1、get host bk_host_innerip property.
	filter := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: "bk_host_innerip",
	}

	attrs := make([]attribute, 0)
	err := db.Table(common.BKTableNameObjAttDes).Find(filter).All(ctx, &attrs)
	if err != nil {
		blog.Errorf("get bk_host_innerip attribute failed, filter: %v, err: %v", filter, err)
		return err
	}

	if len(attrs) > 1 {
		return errors.New("get multiple bk_host_innerip fields")
	}

	if len(attrs) == 0 {
		return errors.New("no bk_host_innerip fields founded")
	}

	if !attrs[0].IsRequired {
		blog.Infof("bk_host_innerip isrequired attribute is already false")
		return nil
	}
	// 2、change the property value isrequired to false.
	data := map[string]interface{}{
		"isrequired": false,
	}

	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, data); err != nil {
		blog.Errorf("change bk_host_innerip required attr to false failed, filter: %+v, data: %+v, err: %v",
			filter, data, err)
		return err
	}
	return nil
}

// addHostAddressingAttr add host addressing attribute
func addHostAddressingAttr(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	now := time.Now()
	addressingAttr := attribute{
		OwnerID:       conf.TenantID,
		ObjectID:      common.BKInnerObjIDHost,
		PropertyID:    "bk_addressing",
		PropertyName:  "寻址方式",
		PropertyGroup: "default",
		IsEditable:    false,
		IsPre:         true,
		PropertyType:  common.FieldTypeEnum,
		Option: []enumVal{{ID: "static", Name: "静态", Type: "text", IsDefault: true},
			{ID: "dynamic", Name: "动态", Type: "text"}},
		Description: "",
		Creator:     common.CCSystemOperatorUserName,
		CreateTime:  &now,
		LastTime:    &now,
	}

	// check if the addressing attribute is already exist
	attrFilter := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: "bk_addressing",
	}

	attrs := make([]attribute, 0)
	err := db.Table(common.BKTableNameObjAttDes).Find(attrFilter).All(ctx, &attrs)
	if err != nil {
		blog.Errorf("check if addressing attribute exists failed, filter: %v, err: %v", attrFilter, err)
		return err
	}

	if len(attrs) > 1 {
		return errors.New("get multiple bk_addressing fields")
	}

	if len(attrs) == 0 {
		// add attributes that are not exist, generate new id and index for them
		newAttrID, err := db.NextSequence(ctx, common.BKTableNameObjAttDes)
		if err != nil {
			blog.Errorf("get new attribute id failed, err: %v", err)
			return err
		}

		attrIdxFilter := map[string]interface{}{
			common.BKObjIDField: common.BKInnerObjIDHost,
		}
		sort := common.BKPropertyIndexField + ":-1"
		maxIdxAttr := new(attribute)
		if err := db.Table(common.BKTableNameObjAttDes).Find(attrIdxFilter).Sort(sort).One(ctx,
			maxIdxAttr); err != nil {
			blog.Errorf("get max host attribute index failed, filter: %v, err: %v", attrIdxFilter, err)
			return err
		}

		addressingAttr.ID = int64(newAttrID)
		addressingAttr.PropertyIndex = maxIdxAttr.PropertyIndex + 1

		if err := db.Table(common.BKTableNameObjAttDes).Insert(ctx, addressingAttr); err != nil {
			blog.Errorf("insert addressing attribute(%#v) failed, err: %v", addressingAttr, err)
			return err
		}

		return addDefaultHostAddressingField(ctx, db, conf)
	}

	// check if exist addressing attribute is the same with the to be added one by name and creator
	existAttr := attrs[0]
	if existAttr.PropertyName != addressingAttr.PropertyName || existAttr.Creator != addressingAttr.Creator {
		return errors.New("get unexpected bk_addressing field")
	}

	return addDefaultHostAddressingField(ctx, db, conf)
}

// addDefaultHostAddressingField update all hosts' addressing field to static if it does not have the field
func addDefaultHostAddressingField(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	hostFilter := map[string]interface{}{
		"bk_addressing": map[string]interface{}{common.BKDBExists: false},
	}
	updateData := map[string]interface{}{
		"bk_addressing": "static",
	}

	for {
		hosts := make([]map[string]interface{}, 0)
		err := db.Table(common.BKTableNameBaseHost).Find(hostFilter).Limit(common.BKMaxPageSize).
			Fields(common.BKHostIDField).All(ctx, &hosts)
		if err != nil {
			blog.Errorf("get host ids to add addressing field failed, err: %v", err)
			return err
		}

		if len(hosts) == 0 {
			break
		}

		hostIDs := make([]int64, len(hosts))
		for index, host := range hosts {
			hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
			if err != nil {
				blog.Errorf("get host id failed, host: %+v, err: %v", host, err)
				return err
			}
			hostIDs[index] = hostID
		}

		filter := map[string]interface{}{
			common.BKHostIDField: map[string]interface{}{
				common.BKDBIN: hostIDs,
			},
		}

		if err := db.Table(common.BKTableNameBaseHost).Update(ctx, filter, updateData); err != nil {
			blog.Errorf("add host instances(%+v) bk_addressing field failed, err: %v", hostIDs, err)
			return err
		}

		if len(hosts) < common.BKMaxPageSize {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

// addProcessIpv6AttrOption add process ipv6 options in bind info
func addProcessIpv6AttrOption(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	bindInfoCond := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDProc,
		common.BKPropertyIDField: common.BKProcBindInfo,
	}

	procAttr := make([]attribute, 0)
	if err := db.Table(common.BKTableNameObjAttDes).Find(bindInfoCond).All(ctx, &procAttr); err != nil {
		blog.Errorf("get process bind info attribute failed, cond: %v, err: %v", bindInfoCond, err)
		return err
	}

	if len(procAttr) != 1 {
		blog.Errorf("process bind info attribute has %d, should be one", len(procAttr))
		return errors.New("process has not one bind info attribute")
	}

	options, err := metadata.ParseSubAttribute(ctx, procAttr[0].Option)
	if err != nil {
		blog.Errorf("parse process bind info attribute's sub attribute(%+v) failed, err: %v", procAttr[0].Option, err)
		return err
	}

	const ipv4Regex = `((1?\d{1,2}|2[0-4]\d|25[0-5])[.]){3}(1?\d{1,2}|2[0-4]\d|25[0-5])`
	// NOCC:tosa/linelength(IPv6正则表达式的长度超过限制，忽略该检查)
	const ipv6Regex = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`

	for index := range options {
		switch options[index].PropertyID {
		case "ip":
			options[index].Option = fmt.Sprintf("(^%s$)|(^%s$)", ipv4Regex, ipv6Regex)
		case "protocol":
			options[index].Option = []enumVal{{ID: "1", Name: "TCP", Type: "text", IsDefault: true},
				{ID: "2", Name: "UDP", Type: "text"}, {ID: "3", Name: "TCP6", Type: "text"},
				{ID: "4", Name: "UDP6", Type: "text"}}
		}
	}

	doc := map[string]interface{}{
		"option": options,
	}
	return db.Table(common.BKTableNameObjAttDes).Update(ctx, bindInfoCond, doc)
}

// adjustHostUnique add agentID unique, do not adjust innerIP + cloudID unique since it can't specify the partial index
func adjustHostUnique(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	// find agent id attribute id
	attrCond := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: common.BKAgentIDField,
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

	// get already exist host uniques
	uniqueCond := map[string]interface{}{common.BKObjIDField: common.BKInnerObjIDHost}
	existUniques := make([]objectUnique, 0)
	err = db.Table(common.BKTableNameObjUnique).Find(uniqueCond).All(ctx, &existUniques)
	if err != nil {
		blog.Errorf("get host unique failed, filter: %v, err: %v", uniqueCond, err)
		return err
	}

	// add agentID unique if it is not exist
	for _, existUnique := range existUniques {
		if len(existUnique.Keys) == 1 && existUnique.Keys[0].ID == attrIDMap[common.BKAgentIDField] {
			return nil
		}
	}

	agentIDUnique := objectUnique{
		Keys: []uniqueKey{
			{Kind: "property", ID: attrIDMap[common.BKAgentIDField]},
		},
		ObjID:    common.BKInnerObjIDHost,
		IsPre:    true,
		OwnerID:  conf.TenantID,
		LastTime: time.Now(),
	}

	newUniqueID, err := db.NextSequence(ctx, common.BKTableNameObjUnique)
	if err != nil {
		blog.Errorf("get new unique id failed, err: %v", err)
		return err
	}

	agentIDUnique.ID = newUniqueID

	if err := db.Table(common.BKTableNameObjUnique).Insert(ctx, agentIDUnique); err != nil {
		blog.Errorf("insert host unique(%#v) failed, err: %v", agentIDUnique, err)
		return err
	}

	return nil
}

// adjustHostUniqueIndex adjust host innerIP + cloudID unique index adding static condition, add agentID unique index
func adjustHostUniqueIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	existIndexArr, err := db.Table(common.BKTableNameBaseHost).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist index for host table failed, err: %v", err)
		return err
	}

	// drop exist host innerIP + cloudID unique index
	for _, index := range existIndexArr {
		if !index.Unique {
			continue
		}

		if len(index.Keys) != 2 {
			continue
		}

		keyMap := make(map[string]struct{})
		for _, v := range index.Keys {
			keyMap[v.Key] = struct{}{}
		}

		if _, exists := keyMap[common.BKCloudIDField]; !exists {
			continue
		}

		_, ipv4Exists := keyMap[common.BKHostInnerIPField]
		_, ipv6Exists := keyMap[common.BKHostInnerIPv6Field]

		if ipv4Exists || ipv6Exists {
			err = db.Table(common.BKTableNameBaseHost).DropIndex(ctx, index.Name)
			if err != nil {
				blog.Errorf("drop index(%#v) failed, err: %v", index, err)
				return err
			}
			continue
		}
	}

	// add new host innerIP + cloudID + addressing == "static" and agentID index
	newIndexes := []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bkHostInnerIP_bkCloudID",
			Keys: bson.D{
				{common.BKHostInnerIPField, 1},
				{common.BKCloudIDField, 1},
			},
			Unique:     true,
			Background: true,
			PartialFilterExpression: map[string]interface{}{
				common.BKHostInnerIPField: map[string]string{common.BKDBType: "string"},
				common.BKCloudIDField:     map[string]string{common.BKDBType: "number"},
				common.BKAddressingField:  "static",
			},
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bkHostInnerIPv6_bkCloudID",
			Keys: bson.D{
				{common.BKHostInnerIPv6Field, 1},
				{common.BKCloudIDField, 1},
			},
			Unique:     true,
			Background: true,
			PartialFilterExpression: map[string]interface{}{
				common.BKHostInnerIPv6Field: map[string]string{common.BKDBType: "string"},
				common.BKCloudIDField:       map[string]string{common.BKDBType: "number"},
				common.BKAddressingField:    "static",
			},
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bkAgentID",
			Keys: bson.D{
				{common.BKAgentIDField, 1},
			},
			Unique:     true,
			Background: true,
			PartialFilterExpression: map[string]interface{}{
				common.BKAgentIDField: map[string]string{
					common.BKDBType: "string",
					common.BKDBGT:   "",
				},
			},
		},
	}

	for _, index := range newIndexes {
		err = db.Table(common.BKTableNameBaseHost).CreateIndex(ctx, index)
		if err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("create index(%#v) failed, err: %v", index, err)
			return err
		}
	}

	return nil
}

type attribute struct {
	BizID         int64       `bson:"bk_biz_id"`
	ID            int64       `bson:"id"`
	OwnerID       string      `bson:"bk_supplier_account"`
	ObjectID      string      `bson:"bk_obj_id"`
	PropertyID    string      `bson:"bk_property_id"`
	PropertyName  string      `bson:"bk_property_name"`
	PropertyGroup string      `bson:"bk_property_group"`
	PropertyIndex int64       `bson:"bk_property_index"`
	Unit          string      `bson:"unit"`
	Placeholder   string      `bson:"placeholder"`
	IsEditable    bool        `bson:"editable"`
	IsPre         bool        `bson:"ispre"`
	IsRequired    bool        `bson:"isrequired"`
	IsReadOnly    bool        `bson:"isreadonly"`
	IsOnly        bool        `bson:"isonly"`
	IsSystem      bool        `bson:"bk_issystem"`
	IsAPI         bool        `bson:"bk_isapi"`
	PropertyType  string      `bson:"bk_property_type"`
	Option        interface{} `bson:"option"`
	Description   string      `bson:"description"`
	Creator       string      `bson:"creator"`
	CreateTime    *time.Time  `bson:"create_time"`
	LastTime      *time.Time  `bson:"last_time"`
}

// enumVal enum option val
type enumVal struct {
	ID        string `bson:"id"`
	Name      string `bson:"name"`
	Type      string `bson:"type"`
	IsDefault bool   `bson:"is_default"`
}

type objectUnique struct {
	ID       uint64      `bson:"id"`
	ObjID    string      `bson:"bk_obj_id"`
	Keys     []uniqueKey `bson:"keys"`
	IsPre    bool        `bson:"ispre"`
	OwnerID  string      `bson:"bk_supplier_account"`
	LastTime time.Time   `bson:"last_time"`
}

type uniqueKey struct {
	Kind string `bson:"key_kind"`
	ID   uint64 `bson:"key_id"`
}
