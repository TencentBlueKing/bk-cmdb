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

package resources

import (
	"sync"

	"gopkg.in/redis.v5"

	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/storage/dal"
)

// db db operator
var db dal.RDB

// txn start transcation
var txn dal.Transcation

// cache redis client
var cache *redis.Client

// err error info
var err errors.CCErrorIf

// lang language info
var lang language.CCLanguageIf

// 锁主要是防止并发对数据的读写
var dbLock sync.RWMutex

// 锁主要是防止并发对数据的读写
var txnLock sync.RWMutex

// 锁主要是防止并发对数据的读写
var cacheLock sync.RWMutex

// 锁主要是防止并发对数据的读写
var errLock sync.RWMutex

// 锁主要是防止并发对数据的读写
var langLock sync.RWMutex

// SetDB set new db client
func SetDB(newDB dal.RDB) {
	dbLock.Lock()
	defer dbLock.Unlock()
	db = newDB
}

// GetDB 要提前自己做初始化，要不会painc
func GetDB() dal.RDB {
	if db == nil {
		panic("db uninitialized")
	}
	dbLock.RLock()
	defer dbLock.RUnlock()

	return db
}

// SetTranscation set new db transcation client
func SetTranscation(newTxn dal.Transcation) {
	txnLock.Lock()
	defer txnLock.Unlock()
	txn = newTxn
}

// GetTranscation 要提前自己做初始化，要不会painc
func GetTranscation() dal.Transcation {
	if txn == nil {
		panic("db Transcation uninitialized")
	}
	txnLock.RLock()
	defer txnLock.RUnlock()

	return txn
}

// SetCache set new redis client
func SetCache(newCache *redis.Client) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	cache = newCache
}

// GetCache 要提前自己做初始化，要不会painc
func GetCache() *redis.Client {
	if cache == nil {
		panic("cache uninitialized")
	}
	cacheLock.RLock()
	defer cacheLock.RUnlock()

	return cache
}

// SetError set new error handle,
func SetError(newErr errors.CCErrorIf) {
	errLock.Lock()
	defer errLock.Unlock()
	err = newErr
}

// GetError 要提前自己做初始化，要不会painc, backone 会做初始化
func GetError() errors.CCErrorIf {
	if err == nil {
		panic("error uninitialized")
	}
	errLock.RLock()
	defer errLock.RUnlock()

	return err
}

// SetLanguage set new language handle, backone 会做初始化
func SetLanguage(newLang language.CCLanguageIf) {
	langLock.Lock()
	defer langLock.Unlock()
	lang = newLang
}

// GetLanguage 要提前自己做初始化，要不会painc
func GetLanguage() language.CCLanguageIf {
	if lang == nil {
		panic("language uninitialized")
	}
	langLock.RLock()
	defer langLock.RUnlock()

	return lang
}
