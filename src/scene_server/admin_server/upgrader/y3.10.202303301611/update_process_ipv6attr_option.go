/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package y3_10_202303301611

import (
	"context"
	"errors"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal"
)

// updateProcessIpv6AttrOption update process ipv6 options in bind info
func updateProcessIpv6AttrOption(ctx context.Context, db dal.RDB) error {
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
	for index := range options {
		switch options[index].PropertyID {
		case "ip":
			options[index].Option = fmt.Sprintf("(^%s$)|(%s)", ipv4Regex, ipv6Regex)
		}
	}

	doc := map[string]interface{}{
		"option": options,
	}
	return db.Table(common.BKTableNameObjAttDes).Update(ctx, bindInfoCond, doc)
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
