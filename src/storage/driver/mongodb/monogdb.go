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

package mongodb

import (
	"strings"
	"time"

	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metric"
	"configcenter/src/common/types"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	dbType "configcenter/src/storage/dal/types"
)

/*
 暂时不支持，多个mongodb实例连接， 暂时不值热更新，所以没有加锁
*/

var (
	dbMap = make(map[string]dal.DB)

	// 在并发的情况下，这里存在panic的问题
	lastInitErr   errors.CCErrorCoder
	lastConfigErr errors.CCErrorCoder
)

// Client  get default error
func Client(prefix ...string) dal.DB {
	var pre string
	if len(prefix) > 0 {
		pre = prefix[0]
	}
	return dbMap[pre]
}

// Table 获取操作db table的对象
func Table(name string) dbType.Table {
	return Client().Table(name)
}

// ParseConfig parse mongodb configuration
func ParseConfig(prefix string, configMap map[string]string) (*mongo.Config, errors.CCErrorCoder) {
	lastConfigErr = nil
	config, err := cc.Mongo(prefix)
	if err != nil {
		return nil, errors.NewCCError(common.CCErrCommConfMissItem, "can't find mongo configuration")
	}
	if config.Address == "" {
		lastConfigErr = errors.NewCCError(common.CCErrCommConfMissItem,
			"Configuration file missing ["+prefix+".host] configuration item")
		return nil, lastConfigErr
	}
	if config.User == "" {
		lastConfigErr = errors.NewCCError(common.CCErrCommConfMissItem,
			"Configuration file missing ["+prefix+".usr] configuration item")
		return nil, lastConfigErr
	}
	if config.Password == "" {
		lastConfigErr = errors.NewCCError(common.CCErrCommConfMissItem,
			"Configuration file missing ["+prefix+".pwd] configuration item")
		return nil, lastConfigErr
	}
	if config.Database == "" {
		lastConfigErr = errors.NewCCError(common.CCErrCommConfMissItem,
			"Configuration file missing ["+prefix+".database] configuration item")
		return nil, lastConfigErr
	}

	return &config, nil
}

// InitClient init mongodb client
func InitClient(prefix string, config *mongo.Config) errors.CCErrorCoder {
	lastInitErr = nil
	var dbErr error
	dbMap[prefix], dbErr = local.NewMgo(config.GetMongoConf(), time.Minute)
	if dbErr != nil {
		blog.Errorf("failed to connect the mongo server, error info is %s", dbErr.Error())
		lastInitErr = errors.NewCCError(common.CCErrCommResourceInitFailed,
			"'"+prefix+".mongodb' initialization failed")
		return lastInitErr
	}
	return nil
}

// UpdateConfig update mongodb configuration
func UpdateConfig(prefix string, config mongo.Config) {
	// 不支持热更新
	return
}

// Healthz check db health status
func Healthz() []metric.HealthItem {
	items := make([]metric.HealthItem, 0)

	for prefix, db := range dbMap {
		if db == nil {
			items = append(items, metric.HealthItem{
				IsHealthy: false,
				Message:   prefix + " db not initialized",
				Name:      types.CCFunctionalityMongo,
			})
			continue
		}

		if err := db.Ping(); err != nil {
			items = append(items, metric.HealthItem{
				IsHealthy: false,
				Message:   prefix + " db connect error. err: " + err.Error(),
				Name:      types.CCFunctionalityMongo,
			})
			continue
		}

		items = append(items, metric.HealthItem{
			IsHealthy: true,
			Name:      types.CCFunctionalityMongo,
		})
	}

	return items
}

// IsDuplicatedError check duplicated error
func IsDuplicatedError(err error) bool {
	if err != nil {
		if strings.Contains(err.Error(), "The existing index") {
			return true
		}
		if strings.Contains(err.Error(), "There's already an index with name") {
			return true
		}
		if strings.Contains(err.Error(), "E11000 duplicate") {
			return true
		}
		if strings.Contains(err.Error(), "IndexOptionsConflict") {
			return true
		}
		if strings.Contains(err.Error(), "all indexes already exist") {
			return true
		}
		if strings.Contains(err.Error(), "already exists with a different name") {
			return true
		}
	}
	return err == dbType.ErrDuplicated
}

// IsNotFoundError check the not found error
func IsNotFoundError(err error) bool {
	return err == dbType.ErrDocumentNotFound
}

// GetDuplicateKey get duplicate key from error, if the error is not a duplicate error, returns the raw error message
// mongodb raw error format example:
// ...{E11000 duplicate key error collection: cmdb.cc_ObjectBase_... index: bkcc_unique_... dup key:
// { bk_inst_name: \"xxx\" }}]},...
func GetDuplicateKey(err error) string {
	if err == nil {
		return ""
	}

	errString := err.Error()
	if !strings.Contains(errString, "E11000 duplicate") {
		return errString
	}

	start := strings.Index(errString, "dup key: ")
	if start == -1 {
		return errString
	}
	start += len("dup key: ") + 1

	end := strings.LastIndex(errString, "}]")
	if end == -1 || end < start {
		return errString
	}

	return errString[start:end]
}

// GetDuplicateValue get duplicate Value from error, if the error is not a duplicate error, returns the raw error message
// mongodb raw error format example:
// Index build failed: ... E11000 duplicate key error collection: cmdb.cc_ObjectBase_0_pub_...:  dup key:
// dup key: { field: "xxxx" }
func GetDuplicateValue(field string, err error) string {
	if field == "" {
		return ""
	}
	if err == nil {
		return ""
	}

	errString := err.Error()
	if !strings.Contains(errString, "E11000 duplicate") {
		return errString
	}

	start := strings.Index(errString, "dup key: ")
	if start == -1 {
		return errString
	}
	start += len("dup key: { " + field + ": ")

	end := strings.LastIndex(errString, " }")
	if end == -1 || end < start {
		return errString
	}

	return errString[start:end]
}
