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
	"time"

	"configcenter/src/common/registerdiscover"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/redis"

	"github.com/spf13/cobra"
)

var Conf *Config

// Config is data structure of cmdb ctl config
type Config struct {
	RegDiscv    string
	RdUser		string
	RdPassword  string
	RdCertFile  string
	RdKeyFile   string
	RdCaFile    string
	MongoURI    string
	MongoRsName string
	RedisConf   redis.Config
}

// AddFlags add flags
func (c *Config) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&c.RegDiscv, "regdiscv", os.Getenv("REGDISCV_ADDR"),
		"the regdiscv address, separated by comma, corresponding environment variable is REGDISCV_ADDR")
	cmd.PersistentFlags().StringVar(&c.RdUser, "rduser", "",
		"user name for authentication in register and discover")
	cmd.PersistentFlags().StringVar(&c.RdPassword, "rdpwd", "",
		"password for authentication in register and discover")
	cmd.PersistentFlags().StringVar(&c.RdCertFile, "rdcert", "",
		"cert file in register and discover")
	cmd.PersistentFlags().StringVar(&c.RdKeyFile, "rdkey", "", "key file in register and discover")
	cmd.PersistentFlags().StringVar(&c.RdCaFile, "rdca", "", "CA file in register and discover")
	cmd.PersistentFlags().StringVar(&c.MongoURI, "mongo-uri", os.Getenv("MONGO_URI"),
		"the mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb, corresponding environment variable is MONGO_URI")
	cmd.PersistentFlags().StringVar(&c.MongoRsName, "mongo-rs-name", "rs0",
		"mongodb replica set name")
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

// Service is register and discover interface and db proxy
type Service struct {
	RegDiscv   *registerdiscover.RegDiscv
	DbProxy dal.RDB
}

// NewRegDiscv creates a service object with register and discover
func NewRegDiscv(cfg *Config) (*Service, error) {
	regdiscvConf := &registerdiscover.Config{
		Host:   cfg.RegDiscv,
		User:   cfg.RdUser,
		Passwd: cfg.RdPassword,
		Cert:   cfg.RdCertFile,
		Key:    cfg.RdKeyFile,
		Ca:     cfg.RdCaFile,
	}
	rd, err := registerdiscover.NewRegDiscv(regdiscvConf)
	if err != nil {
		return nil, err
	}

	service := &Service{
		RegDiscv: rd,
	}
	if err := service.RegDiscv.Ping(); err != nil {
		return nil, err
	}
	return service, nil
}

// NewMongoService creates a service object with db proxy
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
