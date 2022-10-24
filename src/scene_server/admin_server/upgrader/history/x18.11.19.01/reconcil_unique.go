// Package x18_11_19_01 TODO
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
package x18_11_19_01

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func createObjectUnitTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tablenames := []string{common.BKTableNameObjUnique}
	for _, tablename := range tablenames {
		exists, err := db.HasTable(ctx, tablename)
		if err != nil {
			return err
		}
		if !exists {
			if err = db.CreateTable(ctx, tablename); err != nil && !db.IsDuplicatedError(err) {
				return err
			}
		}
	}
	return nil
}

// Attribute TODO
type Attribute struct {
	ID                int64       `json:"id" bson:"id"`
	OwnerID           string      `json:"bk_supplier_account" bson:"bk_supplier_account"`
	ObjectID          string      `json:"bk_obj_id" bson:"bk_obj_id"`
	PropertyID        string      `json:"bk_property_id" bson:"bk_property_id"`
	PropertyName      string      `json:"bk_property_name" bson:"bk_property_name"`
	PropertyGroup     string      `json:"bk_property_group" bson:"bk_property_group"`
	PropertyGroupName string      `json:"bk_property_group_name" bson:"-"`
	PropertyIndex     int64       `json:"bk_property_index" bson:"bk_property_index"`
	Unit              string      `json:"unit" bson:"unit"`
	Placeholder       string      `json:"placeholder" bson:"placeholder"`
	IsEditable        bool        `json:"editable" bson:"editable"`
	IsPre             bool        `json:"ispre" bson:"ispre"`
	IsRequired        bool        `json:"isrequired" bson:"isrequired"`
	IsReadOnly        bool        `json:"isreadonly" bson:"isreadonly"`
	IsOnly            bool        `json:"isonly" bson:"isonly"`
	IsSystem          bool        `json:"bk_issystem" bson:"bk_issystem"`
	IsAPI             bool        `json:"bk_isapi" bson:"bk_isapi"`
	PropertyType      string      `json:"bk_property_type" bson:"bk_property_type"`
	Option            interface{} `json:"option" bson:"option"`
	Description       string      `json:"description" bson:"description"`
	Creator           string      `json:"creator" bson:"creator"`
	CreateTime        *time.Time  `json:"create_time" bson:"create_time"`
	LastTime          *time.Time  `json:"last_time" bson:"last_time"`
}

func reconcilUnique(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	oldAttributes := []Attribute{}
	err := db.Table(common.BKTableNameObjAttDes).Find(nil).All(ctx, &oldAttributes)
	if err != nil {
		return err
	}
	var obj2IsOnlyProperty = map[string][]Attribute{}
	var propertyIDToProperty = map[string]Attribute{}

	var keyfunc = func(a, b string) string { return a + ":" + b }
	for _, oldAttr := range oldAttributes {
		if oldAttr.IsOnly {
			obj2IsOnlyProperty[oldAttr.ObjectID] = append(obj2IsOnlyProperty[oldAttr.ObjectID], oldAttr)
		}
		propertyIDToProperty[keyfunc(oldAttr.ObjectID, oldAttr.PropertyID)] = oldAttr
	}

	shouldCheck := []string{keyfunc(common.BKInnerObjIDHost, common.BKAssetIDField),
		keyfunc(common.BKInnerObjIDHost, common.BKCloudIDField),
		keyfunc(common.BKInnerObjIDHost, common.BKHostInnerIPField),
		keyfunc(common.BKInnerObjIDProc, common.BKAppIDField),
		keyfunc(common.BKInnerObjIDProc, common.BKProcNameField),
		keyfunc(common.BKInnerObjIDProc, common.BKAppIDField),
		keyfunc(common.BKInnerObjIDProc, common.BKFuncIDField),
		keyfunc(common.BKInnerObjIDApp, common.BKAppNameField),
		keyfunc(common.BKInnerObjIDSet, common.BKAppIDField),
		keyfunc(common.BKInnerObjIDSet, common.BKSetNameField),
		keyfunc(common.BKInnerObjIDSet, common.BKInstParentStr),
		keyfunc(common.BKInnerObjIDModule, common.BKAppIDField),
		keyfunc(common.BKInnerObjIDModule, common.BKSetIDField),
		keyfunc(common.BKInnerObjIDModule, common.BKModuleNameField),
	}

	checkKeysShouldExists(propertyIDToProperty, shouldCheck)

	uniques := []objectUnique{
		// host
		// 产品调整，回撤
		// {
		// 	ObjID:     common.BKInnerObjIDHost,
		// 	MustCheck: true,
		// 	Keys: []metadata.UniqueKey{
		// 		{
		// 			Kind: metadata.UniqueKeyKindProperty,
		// 			ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDHost, common.BKAssetIDField)].ID),
		// 		},
		// 	},
		// 	Ispre:    true,
		// 	OwnerID:  conf.OwnerID,
		// 	LastTime: metadata.Now(),
		// },
		{
			ObjID:     common.BKInnerObjIDHost,
			MustCheck: true,
			Keys: []UniqueKey{
				{
					Kind: UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDHost, common.BKCloudIDField)].ID),
				},
				{
					Kind: UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDHost, common.BKHostInnerIPField)].ID),
				},
			},
			Ispre:    true,
			OwnerID:  conf.OwnerID,
			LastTime: Now(),
		},
		// process
		{
			ObjID:     common.BKInnerObjIDProc,
			MustCheck: true,
			Keys: []UniqueKey{
				{
					Kind: UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDProc, common.BKAppIDField)].ID),
				},
				{
					Kind: UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDProc, common.BKProcNameField)].ID),
				},
			},
			Ispre:    true,
			OwnerID:  conf.OwnerID,
			LastTime: Now(),
		},
		{
			ObjID:     common.BKInnerObjIDProc,
			MustCheck: false,
			Keys: []UniqueKey{
				{
					Kind: UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDProc, common.BKAppIDField)].ID),
				},
				{
					Kind: UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDProc, common.BKFuncIDField)].ID),
				},
			},
			Ispre:    true,
			OwnerID:  conf.OwnerID,
			LastTime: Now(),
		},
		// biz
		{
			ObjID:     common.BKInnerObjIDApp,
			MustCheck: true,
			Keys: []UniqueKey{
				{
					Kind: UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDApp, common.BKAppNameField)].ID),
				},
			},
			Ispre:    true,
			OwnerID:  conf.OwnerID,
			LastTime: Now(),
		},
		// set
		{
			ObjID:     common.BKInnerObjIDSet,
			MustCheck: true,
			Keys: []UniqueKey{
				{
					Kind: UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDSet, common.BKAppIDField)].ID),
				},
				{
					Kind: UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDSet, common.BKSetNameField)].ID),
				},
				{
					Kind: UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDSet, common.BKInstParentStr)].ID),
				},
			},
			Ispre:    true,
			OwnerID:  conf.OwnerID,
			LastTime: Now(),
		},
		// module
		{
			ObjID:     common.BKInnerObjIDModule,
			MustCheck: true,
			Keys: []UniqueKey{
				{
					Kind: UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDModule, common.BKAppIDField)].ID),
				},
				{
					Kind: UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDModule, common.BKSetIDField)].ID),
				},
				{
					Kind: UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDModule, common.BKModuleNameField)].ID),
				},
			},
			Ispre:    true,
			OwnerID:  conf.OwnerID,
			LastTime: Now(),
		},
		// cloud area
		{
			ObjID:     common.BKInnerObjIDPlat,
			MustCheck: true,
			Keys: []UniqueKey{
				{
					Kind: UniqueKeyKindProperty,
					ID:   uint64(propertyIDToProperty[keyfunc(common.BKInnerObjIDPlat, common.BKCloudNameField)].ID),
				},
			},
			Ispre:    true,
			OwnerID:  conf.OwnerID,
			LastTime: Now(),
		},
	}

	for objID, oldAttrs := range obj2IsOnlyProperty {
		keys := []UniqueKey{}
		ownerID := conf.OwnerID
		allPreset := true
		for _, oldAttr := range oldAttrs {
			keys = append(keys, UniqueKey{
				Kind: UniqueKeyKindProperty,
				ID:   uint64(oldAttr.ID),
			})
			ownerID = oldAttr.OwnerID
			if !oldAttr.IsPre || (oldAttr.IsPre && oldAttr.PropertyID == common.BKInstNameField) {
				allPreset = false
			}
		}
		if allPreset {
			continue
		}

		unique := objectUnique{
			ObjID:     objID,
			MustCheck: true,
			Keys:      keys,
			Ispre:     false,
			OwnerID:   ownerID,
			LastTime:  Now(),
		}
		uniques = append(uniques, unique)
	}

	uniqueIDs, err := db.NextSequences(ctx, common.BKTableNameObjUnique, len(uniques))
	if err != nil {
		return err
	}

	for index, unique := range uniques {
		exists, err := isUniqueExists(ctx, db, conf, unique)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		unique.ID = uniqueIDs[index]
		if err := db.Table(common.BKTableNameObjUnique).Insert(ctx, unique); err != nil {
			return err
		}
	}

	if err := db.Table(common.BKTableNameObjAttDes).DropColumn(ctx, common.BKIsOnlyField); err != nil {
		return err
	}

	return nil
}

func checkKeysShouldExists(m map[string]Attribute, shouldExistKeys []string) []string {
	notValidKeys := []string{}
	for _, k := range shouldExistKeys {
		if _, ok := m[k]; !ok {
			notValidKeys = append(notValidKeys, k)
		}
	}
	return notValidKeys
}

func isUniqueExists(ctx context.Context, db dal.RDB, conf *upgrader.Config, unique objectUnique) (bool, error) {
	keyhash := unique.KeysHash()
	uniqueCond := condition.CreateCondition()
	uniqueCond.Field(common.BKObjIDField).Eq(unique.ObjID)
	uniqueCond.Field(common.BKOwnerIDField).Eq(conf.OwnerID)
	existUniques := []objectUnique{}

	err := db.Table(common.BKTableNameObjUnique).Find(uniqueCond.ToMapStr()).All(ctx, &existUniques)
	if err != nil {
		return false, err
	}

	for _, uni := range existUniques {
		if uni.KeysHash() == keyhash {
			return true, nil
		}
	}
	return false, nil

}

type objectUnique struct {
	ID        uint64      `json:"id" bson:"id"`
	ObjID     string      `json:"bk_obj_id" bson:"bk_obj_id"`
	MustCheck bool        `json:"must_check" bson:"must_check"`
	Keys      []UniqueKey `json:"keys" bson:"keys"`
	Ispre     bool        `json:"ispre" bson:"ispre"`
	OwnerID   string      `json:"bk_supplier_account" bson:"bk_supplier_account"`
	LastTime  time.Time   `json:"last_time" bson:"last_time"`
}

// Time TODO
type Time metadata.Time

// Now TODO
func Now() time.Time {
	return time.Now().UTC()
}

// UniqueKey TODO
type UniqueKey struct {
	Kind string `json:"key_kind" bson:"key_kind"`
	ID   uint64 `json:"key_id" bson:"key_id"`
}

// KeysHash TODO
func (o objectUnique) KeysHash() string {
	keys := []string{}
	for _, key := range o.Keys {
		keys = append(keys, fmt.Sprintf("%s:%d", key.Kind, key.ID))
	}
	sort.Strings(keys)
	return strings.Join(keys, "#")
}

const (
	// UniqueKeyKindProperty TODO
	UniqueKeyKindProperty = "property"
)
