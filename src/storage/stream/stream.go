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

package stream

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/stream/event"
	"configcenter/src/storage/stream/types"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// Stream Interface defines all the functionality it have.
type Interface interface {
	List(ctx context.Context, opts *types.ListOptions) (ch chan *types.Event, err error)
	Watch(ctx context.Context, opts *types.WatchOptions) (*types.Watcher, error)
	ListWatch(ctx context.Context, opts *types.ListWatchOptions) (*types.Watcher, error)
}

func NewStream(conf local.MongoConf) (Interface, error) {
	connStr, err := connstring.Parse(conf.URI)
	if nil != err {
		return nil, err
	}
	if conf.RsName == "" {
		return nil, fmt.Errorf("rsName not set")
	}

	timeout := 15 * time.Second
	conOpt := options.ClientOptions{
		MaxPoolSize:    &conf.MaxOpenConns,
		MinPoolSize:    &conf.MaxIdleConns,
		ConnectTimeout: &timeout,
		ReplicaSet:     &conf.RsName,
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(conf.URI), &conOpt)
	if nil != err {
		return nil, err
	}
	if err := client.Connect(context.TODO()); nil != err {
		return nil, err
	}

	event, err := event.NewEvent(client, connStr.Database)
	if err != nil {
		return nil, fmt.Errorf("new event failed, err: %v", err)
	}
	return event, nil
}
