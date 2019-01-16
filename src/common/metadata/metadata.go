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
	"errors"
	"fmt"
	"strconv"

	"configcenter/src/common/mapstr"
)

type ModelKind string

const (
	// LabelTrue and LabelFalse is used to define if the label value is a bool value,
	// which is true or false.
	LabelTrue  string = "true"
	LabelFalse string = "false"

	// LabelBusinessID is used to define a label key which value is the business number value.
	LabelBusinessID string = "bk_biz_id"

	// LabelModelKind is used to define a label key which describe what kind of this object is.
	// this label key is always used with the value of PublicModelKindValue or BusinessModelKindValue.
	LabelModelKind         string    = "object_kind"
	PublicModelKindValue   ModelKind = "public_object"
	BusinessModelKindValue ModelKind = "business_object"
)

var (
	LabelKeyNotExistError = errors.New("label key does not exist")
)

var BizLabelNotExist = mapstr.MapStr{"metadata.label.bk_biz_id": mapstr.MapStr{"$exists": false}}

func PublicAndBizCondition(meta Metadata) mapstr.MapStr {
	exist, bizID := meta.Label.Get(LabelBusinessID)
	if false == exist {
		bizID = ""
	}
	condArr := make([]mapstr.MapStr, 0)
	condArr = append(condArr, BizLabelNotExist, mapstr.MapStr{"metadata.label.bk_biz_id": bizID})
	return mapstr.MapStr{"$or": condArr}
}

const (
	BKMetadata string = "metadata"
	BKLabel    string = "label"
)

// Label define the Label type used to limit the scope of application of resources
type Label map[string]string

func NewMetaDataFromBusinessID(value string) Metadata {
	label := make(Label)
	label[LabelBusinessID] = value
	meta := Metadata{Label: label}
	return meta
}

func GetBusinessIDFromMeta(data interface{}) string {
	if nil == data {
		return ""
	}
	tmp, ok := data.(map[string]interface{})
	if !ok {
		return ""
	}
	label, ok := tmp[BKLabel].(map[string]interface{})
	if !ok {
		return ""
	}
	bizID, ok := label[LabelBusinessID].(string)
	if !ok {
		return ""
	}
	return bizID
}

func NewMetaDataFromMap(mapData mapstr.MapStr) *Metadata {
	data, exsit := mapData.Get(BKMetadata)
	if !exsit {
		return nil
	}

	tmp, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}
	label, ok := tmp[BKLabel].(map[string]interface{})
	if !ok {
		return nil
	}
	bizID, ok := label[LabelBusinessID].(string)
	if !ok {
		return nil
	}

	return &Metadata{Label: Label{LabelBusinessID: bizID}}
}

// Metadata  used to define the metadata for the resources
type Metadata struct {
	Label Label `field:"label" json:"label" bson:"label"`
}

func (label Label) Set(key, value string) {
	label[key] = value
}

func (label Label) Get(key string) (exist bool, value string) {
	value, exist = label[key]
	return
}

// isTrue is used to check if the label key is a true value or not.
// if the key does not exist, it will return with a LabelKeyNotExistError,
// which can be used to check it, if you do care about it.
func (label Label) IsTrue(key string) (bool, error) {
	value, exist := label[key]
	if !exist {
		return false, LabelKeyNotExistError
	}

	return strconv.ParseBool(value)
}

// int64 is used to get a int64 value from a label key.
// if the key does not exist, it will return with a LabelKeyNotExistError,
// which can be used to check it, if you do care about it.
func (label Label) Int64(key string) (int64, error) {
	value, exist := label[key]
	if !exist {
		return 0, LabelKeyNotExistError
	}

	return strconv.ParseInt(value, 10, 64)
}

func (label Label) GetBusinessID() (int64, error) {
	return label.Int64(LabelBusinessID)
}

func (label Label) SetBusinessID(id int64) {
	label[LabelBusinessID] = strconv.FormatInt(id, 10)
}

func (label Label) GetModelKind() (ModelKind, error) {
	kind, exist := label[LabelModelKind]
	if !exist {
		return "", LabelKeyNotExistError
	}

	switch ModelKind(kind) {
	case PublicModelKindValue, BusinessModelKindValue:
		return ModelKind(kind), nil
	default:
		return "", fmt.Errorf("unknown model kind %s", kind)
	}
}

func (label Label) SetModelKind(kind ModelKind) {
	label[LabelModelKind] = string(kind)
}
