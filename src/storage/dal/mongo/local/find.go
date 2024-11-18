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

package local

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Find define a find operation
type Find struct {
	*Collection

	projection map[string]int
	filter     types.Filter
	start      int64
	limit      int64
	sort       bson.D

	option types.FindOpts
}

// Fields 查询字段
func (f *Find) Fields(fields ...string) types.Find {
	for _, field := range fields {
		if len(field) <= 0 {
			continue
		}
		f.projection[field] = 1
	}
	return f
}

// Sort 查询排序
// sort支持多字段最左原则排序
// sort值为"host_id, -host_name"和sort值为"host_id:1, host_name:-1"是一样的，都代表先按host_id递增排序，再按host_name递减排序
func (f *Find) Sort(sort string) types.Find {
	if sort != "" {
		sortArr := strings.Split(sort, ",")
		f.sort = bson.D{}
		for _, sortItem := range sortArr {
			sortItemArr := strings.Split(strings.TrimSpace(sortItem), ":")
			sortKey := strings.TrimLeft(sortItemArr[0], "+-")
			if len(sortItemArr) == 2 {
				sortDescFlag := strings.TrimSpace(sortItemArr[1])
				if sortDescFlag == "-1" {
					f.sort = append(f.sort, bson.E{sortKey, -1})
					continue
				}
				f.sort = append(f.sort, bson.E{sortKey, 1})
				continue
			}
			if strings.HasPrefix(sortItemArr[0], "-") {
				f.sort = append(f.sort, bson.E{sortKey, -1})
				continue
			}
			f.sort = append(f.sort, bson.E{sortKey, 1})
		}
	}

	return f
}

// Start 查询上标
func (f *Find) Start(start uint64) types.Find {
	// change to int64,后续改成int64
	dbStart := int64(start)
	f.start = dbStart
	return f
}

// Limit 查询限制
func (f *Find) Limit(limit uint64) types.Find {
	// change to int64,后续改成int64
	dbLimit := int64(limit)
	f.limit = dbLimit
	return f
}

// All 查询多个
func (f *Find) All(ctx context.Context, result interface{}) error {
	mtc.collectOperCount(f.collName, findOper)

	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	defer func() {
		mtc.collectOperDuration(f.collName, findOper, time.Since(start))
	}()

	err := validHostType(f.collName, f.projection, result, rid)
	if err != nil {
		return err
	}

	findOpts := f.generateMongoOption()
	// 查询条件为空时候，mongodb 不返回数据
	if f.filter == nil {
		f.filter = bson.M{}
	}

	opt := getCollectionOption(ctx)

	return f.tm.AutoRunWithTxn(ctx, f.cli.Client(), func(ctx context.Context) error {
		cursor, err := f.cli.Database().Collection(f.collName, opt).Find(ctx, f.filter, findOpts)
		if err != nil {
			mtc.collectErrorCount(f.collName, findOper)
			return err
		}
		return cursor.All(ctx, result)
	})
}

// List 查询多个数据， 当分页中start值为零的时候返回满足条件总行数
func (f *Find) List(ctx context.Context, result interface{}) (int64, error) {
	mtc.collectOperCount(f.collName, findOper)

	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	defer func() {
		mtc.collectOperDuration(f.collName, findOper, time.Since(start))
	}()

	err := validHostType(f.collName, f.projection, result, rid)
	if err != nil {
		return 0, err
	}

	findOpts := f.generateMongoOption()
	// 查询条件为空时候，mongodb 不返回数据
	if f.filter == nil {
		f.filter = bson.M{}
	}

	opt := getCollectionOption(ctx)

	var total int64
	err = f.tm.AutoRunWithTxn(ctx, f.cli.Client(), func(ctx context.Context) error {
		if f.start == 0 || (f.option.WithCount != nil && *f.option.WithCount) {
			var cntErr error
			total, cntErr = f.cli.Database().Collection(f.collName, opt).CountDocuments(ctx, f.filter)
			if cntErr != nil {
				return cntErr
			}
		}
		cursor, err := f.cli.Database().Collection(f.collName, opt).Find(ctx, f.filter, findOpts)
		if err != nil {
			mtc.collectErrorCount(f.collName, findOper)
			return err
		}
		return cursor.All(ctx, result)
	})

	return total, nil
}

// One 查询一个
func (f *Find) One(ctx context.Context, result interface{}) error {
	mtc.collectOperCount(f.collName, findOper)

	start := time.Now()
	rid := ctx.Value(common.ContextRequestIDField)
	defer func() {
		mtc.collectOperDuration(f.collName, findOper, time.Since(start))
	}()

	err := validHostType(f.collName, f.projection, result, rid)
	if err != nil {
		return err
	}

	findOpts := f.generateMongoOption()

	// 查询条件为空时候，mongodb panic
	if f.filter == nil {
		f.filter = bson.M{}
	}

	opt := getCollectionOption(ctx)
	return f.tm.AutoRunWithTxn(ctx, f.cli.Client(), func(ctx context.Context) error {
		cursor, err := f.cli.Database().Collection(f.collName, opt).Find(ctx, f.filter, findOpts)
		if err != nil {
			mtc.collectErrorCount(f.collName, findOper)
			return err
		}

		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			return cursor.Decode(result)
		}
		return types.ErrDocumentNotFound
	})

}

// Count 统计数量(非事务)
func (f *Find) Count(ctx context.Context) (uint64, error) {
	mtc.collectOperCount(f.collName, countOper)

	start := time.Now()
	defer func() {
		mtc.collectOperDuration(f.collName, countOper, time.Since(start))
	}()

	if f.filter == nil {
		f.filter = bson.M{}
	}

	opt := getCollectionOption(ctx)

	sessCtx, _, useTxn, err := f.tm.GetTxnContext(ctx, f.cli.Client())
	if err != nil {
		return 0, err
	}
	if !useTxn {
		// not use transaction.
		cnt, err := f.cli.Database().Collection(f.collName, opt).CountDocuments(ctx, f.filter)
		if err != nil {
			mtc.collectErrorCount(f.collName, countOper)
			return 0, err
		}

		return uint64(cnt), err
	} else {
		// use transaction
		cnt, err := f.cli.Database().Collection(f.collName, opt).CountDocuments(sessCtx, f.filter)
		// do not release th session, otherwise, the session will be returned to the
		// session pool and will be reused. then mongodb driver will increase the transaction number
		// automatically and do read/write retry if policy is set.
		// mongo.CmdbReleaseSession(ctx, session)
		if err != nil {
			mtc.collectErrorCount(f.collName, countOper)
			return 0, err
		}
		return uint64(cnt), nil
	}
}

// Option TODO
func (f *Find) Option(opts ...*types.FindOpts) {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.WithObjectID != nil {
			f.option.WithObjectID = opt.WithObjectID
		}
		if opt.WithCount != nil {
			f.option.WithCount = opt.WithCount
		}
	}
}

func (f *Find) generateMongoOption() *options.FindOptions {
	findOpts := &options.FindOptions{}
	if f.projection == nil {
		f.projection = make(map[string]int, 0)
	}
	if f.option.WithObjectID != nil && *f.option.WithObjectID {
		// mongodb 要求，当有字段设置未1, 不设置都不显示
		// 没有设置projection 的时候，返回所有字段
		if len(f.projection) > 0 {
			f.projection["_id"] = 1
		}
	} else {
		if _, exists := f.projection["_id"]; !exists {
			f.projection["_id"] = 0
		}
	}
	if len(f.projection) != 0 {
		findOpts.Projection = f.projection
	}

	if f.start != 0 {
		findOpts.SetSkip(f.start)
	}
	if f.limit != 0 {
		findOpts.SetLimit(f.limit)
	}
	if len(f.sort) != 0 {
		findOpts.SetSort(f.sort)
	}

	return findOpts
}

var hostSpecialFieldMap = map[string]bool{
	common.BKHostInnerIPField:   true,
	common.BKHostOuterIPField:   true,
	common.BKOperatorField:      true,
	common.BKBakOperatorField:   true,
	common.BKHostInnerIPv6Field: true,
	common.BKHostOuterIPv6Field: true,
}

// validHostType valid if host query uses specified type that transforms ip & operator array to string
func validHostType(collection string, projection map[string]int, result interface{}, rid interface{}) error {
	if result == nil {
		blog.Errorf("host query result is nil, rid: %s", rid)
		return fmt.Errorf("host query result type invalid")
	}

	if collection != common.BKTableNameBaseHost {
		return nil
	}

	// check if specified fields include special fields
	if len(projection) != 0 {
		needCheck := false
		for field := range projection {
			if hostSpecialFieldMap[field] {
				needCheck = true
				break
			}
		}
		if !needCheck {
			return nil
		}
	}

	resType := reflect.TypeOf(result)
	if resType.Kind() != reflect.Ptr {
		blog.Errorf("host query result type(%v) not pointer type, rid: %v", resType, rid)
		return fmt.Errorf("host query result type invalid")
	}
	// if result is *map[string]interface{} type, it must be *metadata.HostMapStr type
	if resType.ConvertibleTo(reflect.TypeOf(&map[string]interface{}{})) {
		if resType != reflect.TypeOf(&metadata.HostMapStr{}) {
			blog.Errorf("host query result type(%v) not match *metadata.HostMapStr type, rid: %v", resType, rid)
			return fmt.Errorf("host query result type invalid")
		}
		return nil
	}

	resElem := resType.Elem()
	switch resElem.Kind() {
	case reflect.Struct:
		err := validHostStructType(resElem, rid)
		if err != nil {
			return err
		}
	case reflect.Slice:
		// check if slice item is valid type, map or struct validation is similar as before
		elem := resElem.Elem()
		if elem.ConvertibleTo(reflect.TypeOf(map[string]interface{}{})) {
			if elem != reflect.TypeOf(metadata.HostMapStr{}) {
				blog.Errorf("host query result type(%v) not match *[]metadata.HostMapStr type", resType)
				return fmt.Errorf("host query result type invalid")
			}
			return nil
		}

		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		if elem.Kind() != reflect.Struct {
			blog.Errorf("host query result type(%v) not struct pointer type or map type", resType)
			return fmt.Errorf("host query result type invalid")
		}
		err := validHostStructType(elem, rid)
		if err != nil {
			return err
		}
	default:
		blog.Errorf("host query result type(%v) not pointer of map, struct or slice, rid: %v", resType, rid)
		return fmt.Errorf("host query result type invalid")
	}
	return nil
}

// validHostStructType validate if *struct type result's special field must be metadata.StringArrayToString type
func validHostStructType(resElem reflect.Type, rid interface{}) error {
	numField := resElem.NumField()
	validType := reflect.TypeOf(metadata.StringArrayToString(""))
	for i := 0; i < numField; i++ {
		field := resElem.Field(i)
		bsonTag := field.Tag.Get("bson")
		if bsonTag == "" {
			blog.Errorf("host query result field(%s) has empty bson tag, rid: %v", field.Name, rid)
			return fmt.Errorf("host query result type invalid")
		}
		if hostSpecialFieldMap[bsonTag] && field.Type != validType {
			blog.Errorf("host query result field type(%v) not match *metadata.StringArrayToString type", field.Type)
			return fmt.Errorf("host query result type invalid")
		}
	}
	return nil
}
