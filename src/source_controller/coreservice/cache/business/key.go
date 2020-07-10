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

package business

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/json"
)

const (
	bizNamespace = common.BKCacheKeyV3Prefix + "biz"
)

var bizKey = keyGenerator{
	namespace:            bizNamespace,
	name:                 bizKeyName,
	listExpireDuration:   5 * time.Minute,
	detailExpireDuration: 5 * time.Minute,
}

var moduleKey = keyGenerator{
	namespace:            bizNamespace,
	name:                 moduleKeyName,
	listExpireDuration:   5 * time.Minute,
	detailExpireDuration: 5 * time.Minute,
}

var setKey = keyGenerator{
	namespace:            bizNamespace,
	name:                 setKeyName,
	listExpireDuration:   5 * time.Minute,
	detailExpireDuration: 5 * time.Minute,
}

type keyName string

const (
	bizKeyName    keyName = "biz"
	moduleKeyName keyName = "module"
	setKeyName            = "set"
)

type keyGenerator struct {
	namespace            string
	name                 keyName
	listExpireDuration   time.Duration
	detailExpireDuration time.Duration
}

// return a key which can indicate whether the resources has already been listed.
func (k keyGenerator) listDoneKey() string {
	return k.namespace + ":" + string(k.name) + ":listdone"
}

func (k keyGenerator) listKeyWithBiz(bizID int64) string {
	if k.name == bizKeyName {
		return k.namespace + ":" + string(k.name) + "_list"
	}
	return fmt.Sprintf("%s:%s_list:%d", k.namespace, string(k.name), bizID)
}

func (k keyGenerator) listLockKeyWithBiz(bizID int64) string {
	if k.name == bizKeyName {
		return k.namespace + ":" + string(k.name) + "_list:lock"
	}
	return fmt.Sprintf("%s:%s_list:lock:%d", k.namespace, string(k.name), bizID)
}

func (k keyGenerator) listExpireKeyWithBiz(bizID int64) string {
	if k.name == bizKeyName {
		return k.namespace + ":" + string(k.name) + "_list:expire"
	}
	return fmt.Sprintf("%s:%s_list:expire:%d", k.namespace, string(k.name), bizID)
}

func (k keyGenerator) genListKeyValue(instID int64, parentID int64, instName string) string {
	return strconv.FormatInt(instID, 10) + ":" + strconv.FormatInt(parentID, 10) + ":" + instName
}

func (k keyGenerator) parseListKeyValue(key string) (int64, int64, string, error) {
	fields := strings.SplitN(key, ":", 3)
	if len(fields) != 3 {
		return 0, 0, "", fmt.Errorf("invalid key: %s", key)
	}

	instID, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, 0, "", fmt.Errorf("key: %s with invalid inst id, err: %v", key, err)
	}

	parentID, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return instID, 0, "", fmt.Errorf("key: %s with invalid parent id, err: %v", key, err)
	}

	name := fields[2]
	if len(name) == 0 {
		return instID, parentID, "", fmt.Errorf("key: %s with empty name, err: %v", key, err)
	}

	return instID, parentID, name, nil
}

func (k keyGenerator) detailKey(instID int64) string {
	return fmt.Sprintf("%s:%s_detail:%d", k.namespace, k.name, instID)
}

func (k keyGenerator) detailLockKey(instID int64) string {
	return fmt.Sprintf("%s:%s_detail:lock:%d", k.namespace, k.name, instID)
}

func (k keyGenerator) detailExpireKey(instID int64) string {
	return fmt.Sprintf("%s:%s_detail:expire:%d", k.namespace, k.name, instID)
}

// this key is to save the document object id(as is _id) relations with the instance id
func (k keyGenerator) objectIDKey() string {
	return k.namespace + ":oid"
}

var customKey = customKeyGen{
	namespace:            bizNamespace + ":custom",
	listExpireDuration:   5 * time.Minute,
	detailExpireDuration: 5 * time.Minute,
}

// for business custom level cache storage usage.
// to generate the custom object level key
type customKeyGen struct {
	namespace            string
	listExpireDuration   time.Duration
	detailExpireDuration time.Duration
}

func (c customKeyGen) listDoneKey(objectID string) string {
	return c.namespace + ":listdone:" + objectID
}

// key to store the mainline topology, except the host object.
func (c customKeyGen) topologyKey() string {
	return c.namespace + ":topology"
}

func (c customKeyGen) topologyValue(rank []string) string {
	return strings.Join(rank, ",")
}

func (c customKeyGen) parseTopologyValue(v string) []string {
	return strings.Split(v, ",")
}

func (c customKeyGen) topologyExpireKey() string {
	return c.namespace + ":topology:expire"
}

// store a business's and object instance list to a same key.
func (c customKeyGen) objListKeyWithBiz(objectID string, bizID int64) string {
	return fmt.Sprintf("%s:%s_list:%d", c.namespace, objectID, bizID)
}

func (c customKeyGen) objListLockKeyWithBiz(objectID string, bizID int64) string {
	return fmt.Sprintf("%s:%s_list:lock:%d", c.namespace, objectID, bizID)
}

func (c customKeyGen) objListExpireKeyWithBiz(objectID string, bizID int64) string {
	return fmt.Sprintf("%s:%s_list:expire:%d", c.namespace, objectID, bizID)
}

func (c customKeyGen) genListKeyValue(objInstID int64, parentID int64, objInstName string) string {
	return strconv.FormatInt(objInstID, 10) + ":" + strconv.FormatInt(parentID, 10) + ":" + objInstName
}

func (c customKeyGen) parseListKeyValue(key string) (int64, int64, string, error) {
	fields := strings.SplitN(key, ":", 3)
	if len(fields) != 3 {
		return 0, 0, "", fmt.Errorf("invalid key: %s", key)
	}

	instID, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, 0, "", fmt.Errorf("key: %s with invalid inst id, err: %v", key, err)
	}

	parentID, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return instID, 0, "", fmt.Errorf("key: %s with invalid parent id, err: %v", key, err)
	}

	name := fields[2]
	if len(name) == 0 {
		return instID, parentID, "", fmt.Errorf("key: %s with empty name, err: %v", key, err)
	}

	return instID, parentID, name, nil
}

func (c customKeyGen) detailKey(objectID string, id int64) string {
	return fmt.Sprintf("%s:%s_detail:%d", c.namespace, objectID, id)
}

func (c customKeyGen) detailLockKey(objectID string, id int64) string {
	return fmt.Sprintf("%s:%s_detail:lock:%d", c.namespace, objectID, id)
}

func (c customKeyGen) detailExpireKey(objectID string, id int64) string {
	return fmt.Sprintf("%s:%s_detail:expire:%d", c.namespace, objectID, id)
}

// this key is to save the document object id(as is _id) relations with the instance id and biz id
func (c customKeyGen) objectIDKey() string {
	return c.namespace + ":oid"
}

func (c customKeyGen) mainlineListDoneKey() string {
	return c.namespace + ":" + mainlineTopologyListDoneKey
}

type oidValue struct {
	biz    int64  `json:"biz"`
	instID int64  `json:"inst_id"`
	obj    string `json:"obj"`
}

// generate the value to save in oid key.
func (c customKeyGen) genObjectIDKeyValue(bizID, objInstID int64, objectID string) string {
	js, _ := json.Marshal(oidValue{
		biz:    bizID,
		instID: objInstID,
		obj:    objectID,
	})
	return string(js)
}

// parse the object key saved in cache to business id and instance id
func (c customKeyGen) parseObjectIDKeyValue(value string) (*oidValue, error) {
	v := new(oidValue)
	if err := json.Unmarshal([]byte(value), v); err != nil {
		return nil, err
	}

	return v, nil
}
