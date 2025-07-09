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

// Package config TODO
package config

import (
	"errors"
	"os"
	"strings"
	"time"

	"configcenter/src/common/ssl"
	"configcenter/src/common/zkclient"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/redis"

	"github.com/spf13/cobra"
)

// Conf TODO
var Conf *Config

// Config TODO
type Config struct {
	ZkAddr      string
	ZkTLS       ssl.TLSClientConfig
	MongoURI    string
	MongoRsName string
	RedisConf   redis.Config
}

// AddFlags add flags
func (c *Config) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&c.ZkAddr, "zk-addr", os.Getenv("ZK_ADDR"),
		"the ip address and port for the zookeeper hosts, separated by comma, corresponding environment variable is ZK_ADDR")
	cmd.PersistentFlags().StringVar(&c.ZkTLS.CAFile, "zk-tls-ca-file", os.Getenv("ZK_TLS_CA_FILE"),
		"the path of TLS CA file for the zookeeper hosts, corresponding environment variable is ZK_TLS_CA_FILE")
	cmd.PersistentFlags().BoolVar(&c.ZkTLS.InsecureSkipVerify,
		"zk-tls-skip-verify", os.Getenv("ZK_TLS_SKIP_VERIFY") == "true",
		"the flag of TLS certificate skip verify for zookeeper, corresponding environment variable is ZK_TLS_SKIP_VERIFY")
	cmd.PersistentFlags().StringVar(&c.ZkTLS.CertFile, "zk-tls-certfile", os.Getenv("ZK_TLS_CERT_FILE"),
		"the path of TLS cert file for zookeeper, corresponding environment variable is ZK_TLS_CERT_FILE")
	cmd.PersistentFlags().StringVar(&c.ZkTLS.KeyFile, "zk-tls-keyfile", os.Getenv("ZK_TLS_KEY_FILE"),
		"the path of TLS key file for zookeeper, corresponding environment variable is ZK_TLS_KEY_FILE")
	cmd.PersistentFlags().StringVar(&c.ZkTLS.Password, "zk-tls-password", os.Getenv("ZK_TLS_PASSWORD"),
		"the password of TLS for zookeeper, corresponding environment variable is ZK_TLS_PASSWORD")
	cmd.PersistentFlags().StringVar(&c.MongoURI, "mongo-uri", os.Getenv("MONGO_URI"),
		"the mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb, corresponding environment variable is MONGO_URI")
	cmd.PersistentFlags().StringVar(&c.MongoRsName, "mongo-rs-name", "rs0", "mongodb replica set name")
	cmd.PersistentFlags().StringVar(&c.RedisConf.Address, "redis-addr", "127.0.0.1:6379",
		"assign redis server address default is 127.0.0.1:6379")
	cmd.PersistentFlags().StringVar(&c.RedisConf.MasterName, "redis-mastername", "",
		"assign redis server master name defalut is null")
	cmd.PersistentFlags().StringVar(&c.RedisConf.Password, "redis-pwd", "",
		"assign redis server password default is null")
	cmd.PersistentFlags().StringVar(&c.RedisConf.SentinelPassword, "redis-sentinelpwd", "",
		"assign the redis sentinel password  default is null")
	cmd.PersistentFlags().StringVar(&c.RedisConf.Database, "redis-database", "0",
		"assign the redis database  default is 0")
	return
}

// Service TODO
type Service struct {
	ZkCli   *zkclient.ZkClient
	DbProxy dal.RDB
}

// NewZkService TODO
func NewZkService(zkAddr string, tlsConfig *ssl.TLSClientConfig) (*Service, error) {
	if zkAddr == "" {
		return nil, errors.New("zk-addr must set via flag or environment variable")
	}
	service := &Service{
		ZkCli: zkclient.NewZkClient(strings.Split(zkAddr, ","), tlsConfig),
	}
	if err := service.ZkCli.Connect(); err != nil {
		return nil, err
	}
	return service, nil
}

// NewMongoService TODO
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
