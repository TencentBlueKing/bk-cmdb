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

// Package config defines the config for cmdb ctl tool
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"configcenter/src/common/core/cc/config"
	"configcenter/src/common/cryptor"
	"configcenter/src/common/zkclient"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/mongo/sharding"
	"configcenter/src/storage/dal/redis"

	"github.com/spf13/cobra"
)

// Conf is the global config
var Conf *Config

// Config is the config for cmdb ctl tool
type Config struct {
	Zk        config.ZkConfig
	MongoConf *MongoConfig
	RedisConf redis.Config
}

// MongoConfig is the mongodb config for cmdb ctl tool
type MongoConfig struct {
	MongoURI    string
	MongoRsName string
	CryptoConf  string
}

// AddFlags add flags
func (c *Config) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&c.Zk.Addr, "zk-addr", os.Getenv("ZK_ADDR"),
		"the ip address and port for the zookeeper hosts, separated by comma, corresponding environment variable is ZK_ADDR")
	cmd.PersistentFlags().StringVar(&c.Zk.User, "zk-user", os.Getenv("ZK_USER"),
		"the user name for zookeeper auth, corresponding environment variable is ZK_USER")
	cmd.PersistentFlags().StringVar(&c.Zk.Password, "zk-password", os.Getenv("ZK_PASSWORD"),
		"the password for zookeeper auth, corresponding environment variable is ZK_PASSWORD")
	cmd.PersistentFlags().StringVar(&c.Zk.TLS.CAFile, "zk-tls-ca-file", os.Getenv("ZK_TLS_CA_FILE"),
		"the path of TLS CA file for the zookeeper hosts, corresponding environment variable is ZK_TLS_CA_FILE")
	cmd.PersistentFlags().BoolVar(&c.Zk.TLS.InsecureSkipVerify,
		"zk-tls-skip-verify", os.Getenv("ZK_TLS_SKIP_VERIFY") == "true",
		"the flag of TLS certificate skip verify for zookeeper, corresponding environment variable is ZK_TLS_SKIP_VERIFY")
	cmd.PersistentFlags().StringVar(&c.Zk.TLS.CertFile, "zk-tls-certfile", os.Getenv("ZK_TLS_CERT_FILE"),
		"the path of TLS cert file for zookeeper, corresponding environment variable is ZK_TLS_CERT_FILE")
	cmd.PersistentFlags().StringVar(&c.Zk.TLS.KeyFile, "zk-tls-keyfile", os.Getenv("ZK_TLS_KEY_FILE"),
		"the path of TLS key file for zookeeper, corresponding environment variable is ZK_TLS_KEY_FILE")
	cmd.PersistentFlags().StringVar(&c.Zk.TLS.Password, "zk-tls-password", os.Getenv("ZK_TLS_PASSWORD"),
		"the password of TLS for zookeeper, corresponding environment variable is ZK_TLS_PASSWORD")
	c.MongoConf = new(MongoConfig)
	cmd.PersistentFlags().StringVar(&c.MongoConf.MongoURI, "mongo-uri", os.Getenv("MONGO_URI"),
		"the mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb, corresponding environment variable is MONGO_URI")
	cmd.PersistentFlags().StringVar(&c.MongoConf.MongoRsName, "mongo-rs-name", "rs0", "mongodb replica set name")
	cmd.PersistentFlags().StringVar(&c.MongoConf.CryptoConf, "crypto-config", "", "mongo crypto config in json format")
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

// Service is the common service for cmdb ctl tool
type Service struct {
	ZkCli   *zkclient.ZkClient
	DbProxy dal.Dal
}

// NewZkService creates a ZK service using the provided ZkConfig.
func NewZkService(zkConf config.ZkConfig) (*Service, error) {
	if zkConf.Addr == "" {
		return nil, errors.New("zk-addr must set via flag or environment variable")
	}
	service := &Service{
		ZkCli: zkclient.NewZkClient(zkConf),
	}
	if err := service.ZkCli.Connect(); err != nil {
		return nil, err
	}
	return service, nil
}

// NewMongoService new mongodb service
func NewMongoService(conf *MongoConfig) (*Service, error) {
	if conf.MongoURI == "" {
		return nil, errors.New("mongo-uri must set via flag or environment variable")
	}

	cryptoConf := new(cryptor.Config)
	if conf.CryptoConf != "" {
		err := json.Unmarshal([]byte(conf.CryptoConf), cryptoConf)
		if err != nil {
			return nil, fmt.Errorf("parse mongodb crypto config failed, err: %v", err)
		}
	}
	crypto, err := cryptor.NewCrypto(cryptoConf)
	if err != nil {
		return nil, fmt.Errorf("new crypto failed, err: %v", err)
	}

	mongoConfig := local.MongoConf{
		MaxOpenConns: mongo.MinimumMaxIdleOpenConns,
		MaxIdleConns: mongo.MinimumMaxIdleOpenConns,
		URI:          conf.MongoURI,
		RsName:       conf.MongoRsName,
	}
	db, err := sharding.NewShardingMongo(mongoConfig, time.Minute, crypto)
	if err != nil {
		return nil, err
	}
	return &Service{
		DbProxy: db,
	}, nil
}
