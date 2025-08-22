/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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
	"errors"
	"fmt"
	"strings"
	"time"

	idgen "configcenter/pkg/id-gen"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util/table"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// OldMongo is the mongodb client for previous version that is only used in migration
// TODO remove this after migrating from 3.14 to 3.15
type OldMongo struct {
	*Mongo
}

var _ DB = new(OldMongo)

// NewOldMgo returns new mongodb client for previous version that is only used in migration
func NewOldMgo(config MongoConf, timeout time.Duration) (DB, error) {
	mgo, err := NewMgo(config, timeout)
	if err != nil {
		return nil, err
	}

	mgo.enableSharding = false
	return &OldMongo{Mongo: mgo}, nil
}

// addOldCollPrefix add cc_ prefix for old version collection name
func (c *OldMongo) addOldCollPrefix(collection string) string {
	if strings.HasPrefix(collection, "cc_") {
		return collection
	}
	return "cc_" + collection
}

// Table collection operation
func (c *OldMongo) Table(collName string) types.Table {
	return &Collection{
		collName: c.addOldCollPrefix(collName),
		Mongo:    c.Mongo,
	}
}

func (c *OldMongo) convSequenceName(sequenceName string) string {
	sequenceName = strings.TrimPrefix(sequenceName, "cc_")
	return "cc_" + c.redirectTable(sequenceName)
}

// NextSequence 获取新序列号(非事务)
func (c *OldMongo) NextSequence(ctx context.Context, sequenceName string) (uint64, error) {
	if c.enableSharding && !c.ignoreTenant {
		return 0, errors.New("next sequence do not need tenant")
	}

	sequenceName = c.convSequenceName(sequenceName)

	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	defer func() {
		blog.V(4).InfoDepthf(2, "mongo next-sequence cost %dms, rid: %v", time.Since(start)/time.Millisecond, rid)
	}()

	// 直接使用新的context，确保不会用到事务,不会因为context含有session而使用分布式事务，防止产生相同的序列号
	ctx = context.Background()

	coll := c.cli.Database().Collection("cc_idgenerator")

	Update := bson.M{
		"$inc":         bson.M{"SequenceID": c.conf.IDGenStep},
		"$setOnInsert": bson.M{"create_time": time.Now()},
		"$set":         bson.M{"last_time": time.Now()},
	}
	filter := bson.M{"_id": sequenceName}
	upsert := true
	returnChange := options.After
	opt := &options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &returnChange,
	}

	doc := Idgen{}
	err := coll.FindOneAndUpdate(ctx, filter, Update, opt).Decode(&doc)
	if err != nil {
		return 0, err
	}
	return doc.SequenceID, err
}

// NextSequences 批量获取新序列号(非事务)
func (c *OldMongo) NextSequences(ctx context.Context, sequenceName string, num int) ([]uint64, error) {
	if num == 0 {
		return make([]uint64, 0), nil
	}

	sequenceName = c.convSequenceName(sequenceName)

	if c.conf.DisableInsert && idgen.IsIDGenSeqName(sequenceName) {
		return nil, errors.New("insertion is disabled")
	}

	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	defer func() {
		blog.V(4).InfoDepthf(2, "mongo next-sequences cost %dms, rid: %v", time.Since(start)/time.Millisecond, rid)
	}()

	// 直接使用新的context，确保不会用到事务,不会因为context含有session而使用分布式事务，防止产生相同的序列号
	ctx = context.Background()

	coll := c.cli.Database().Collection("cc_idgenerator")

	Update := bson.M{
		"$inc":         bson.M{"SequenceID": num * c.conf.IDGenStep},
		"$setOnInsert": bson.M{"create_time": time.Now()},
		"$set":         bson.M{"last_time": time.Now()},
	}
	filter := bson.M{"_id": sequenceName}
	upsert := true
	returnChange := options.After
	opt := &options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &returnChange,
	}

	doc := Idgen{}
	err := coll.FindOneAndUpdate(ctx, filter, Update, opt).Decode(&doc)
	if err != nil {
		return nil, err
	}

	sequences := make([]uint64, num)
	for i := 0; i < num; i++ {
		sequences[i] = uint64((i-num+1)*c.conf.IDGenStep) + doc.SequenceID
	}

	return sequences, err
}

// HasTable 判断是否存在集合
func (c *OldMongo) HasTable(ctx context.Context, collName string) (bool, error) {
	collName = c.addOldCollPrefix(collName)

	cursor, err := c.cli.Database().ListCollections(ctx, bson.M{"name": collName, "type": "collection"})
	if err != nil {
		return false, err
	}

	defer cursor.Close(ctx)
	return cursor.Next(ctx), nil
}

// DropTable 移除集合
func (c *OldMongo) DropTable(ctx context.Context, collName string) error {
	collName = c.addOldCollPrefix(collName)
	return c.cli.Database().Collection(collName).Drop(ctx)
}

// CreateTable 创建集合
func (c *OldMongo) CreateTable(ctx context.Context, collName string) error {
	collName = c.addOldCollPrefix(collName)
	return c.cli.Database().CreateCollection(ctx, collName)
}

// RenameTable 更新集合名称
func (c *OldMongo) RenameTable(ctx context.Context, prevName, currName string) error {
	if !strings.HasPrefix(currName, "cc_") && table.NeedPreImageTable(strings.TrimPrefix(currName, "default_")) {
		cmd := bson.D{
			{"collMod", prevName},
			{"changeStreamPreAndPostImages", bson.M{"enabled": true}},
		}
		if err := c.cli.Client().Database(c.GetDBName()).RunCommand(ctx, cmd).Err(); err != nil {
			return fmt.Errorf("enable change stream pre and post image failed: %v", err)
		}
	}

	cmd := bson.D{
		{"renameCollection", c.cli.DBName() + "." + prevName},
		{"to", c.cli.DBName() + "." + currName},
	}
	return c.cli.Client().Database("admin").RunCommand(ctx, cmd).Err()
}
