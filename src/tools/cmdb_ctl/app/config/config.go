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

package config

import (
	"errors"
	"os"
	"strings"
	"time"

	"configcenter/src/common/zkclient"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"

	"github.com/spf13/cobra"
)

var Conf *Config

type Config struct {
	ZkAddr      string
	MongoURI    string
	MongoRsName string
}

// AddFlags add flags
func (c *Config) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&c.ZkAddr, "zk-addr", os.Getenv("ZK_ADDR"), "the ip address and port for the zookeeper hosts, separated by comma, corresponding environment variable is ZK_ADDR")
	// TODO add zkuser and zkpwd
	cmd.PersistentFlags().StringVar(&c.MongoURI, "mongo-uri", os.Getenv("MONGO_URI"), "the mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb, corresponding environment variable is MONGO_URI")
	cmd.PersistentFlags().StringVar(&c.MongoRsName, "mongo-rs-name", "rs0", "mongodb replica set name")
}

type Service struct {
	ZkCli   *zkclient.ZkClient
	DbProxy dal.RDB
}

func NewZkService(zkAddr string) (*Service, error) {
	if zkAddr == "" {
		return nil, errors.New("zk-addr must set via flag or environment variable")
	}
	service := &Service{
		ZkCli: zkclient.NewZkClient(strings.Split(zkAddr, ",")),
	}
	if err := service.ZkCli.Connect(); err != nil {
		return nil, err
	}
	return service, nil
}

func NewMongoService(mongoURI string, mongoRsName string) (*Service, error) {
	if mongoURI == "" {
		return nil, errors.New("mongo-uri must set via flag or environment variable")
	}
	mongoConfig := local.MongoConf{
		MaxOpenConns: mongo.DefaultMaxOpenConns,
		MaxIdleConns: mongo.MinimumMaxIdleOpenConns,
		URI:          mongoURI,
		RsName:       mongoRsName,
	}
	db, err := local.NewMgo(mongoConfig, time.Minute)
	if err != nil {
		return nil, err
	}
	return &Service{
		DbProxy: db,
	}, nil
}
