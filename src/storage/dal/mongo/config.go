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

package mongo

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
)

const (
	// if maxOpenConns isn't configured, use default value
	DefaultMaxOpenConns = 1000
	// if maxOpenConns exceeds maximum value, use maximum value
	MaximumMaxOpenConns = 3000
	// if maxIDleConns is less than minimum value, use minimum value
	MinimumMaxIdleOpenConns = 50
)

// Config config
type Config struct {
	Connect      string
	Address      string
	User         string
	Password     string
	Port         string
	Database     string
	Mechanism    string
	MaxOpenConns uint64
	MaxIdleConns uint64
	RsName       string
}

// BuildURI return mongo uri according to  https://docs.mongodb.com/manual/reference/connection-string/
// format example: mongodb://[username:password@]host1[:port1][,host2[:port2],...[,hostN[:portN]]][/[database][?options]]
func (c Config) BuildURI() string {
	if c.Connect != "" {
		return c.Connect
	}

	if !strings.Contains(c.Address, ":") && len(c.Port) > 0 {
		c.Address = c.Address + ":" + c.Port
	}

	c.User = strings.Replace(c.User, "@", "%40", -1)
	c.Password = strings.Replace(c.Password, "@", "%40", -1)
	c.User = strings.Replace(c.User, ":", "%3a", -1)
	c.Password = strings.Replace(c.Password, ":", "%3a", -1)
	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s?authMechanism=%s", c.User, c.Password, c.Address, c.Database, c.Mechanism)
	return uri
}

// ParseConfigFromKV returns a new config
func ParseConfigFromKV(prefix string, configmap map[string]string) Config {
	c := Config{
		Address:   configmap[prefix+".host"],
		Port:      configmap[prefix+".port"],
		User:      configmap[prefix+".usr"],
		Password:  configmap[prefix+".pwd"],
		Database:  configmap[prefix+".database"],
		Mechanism: configmap[prefix+".mechanism"],
		RsName:    configmap[prefix+".rsName"],
	}

	if c.RsName == "" {
		blog.Errorf("rsName not set")
	}
	if c.Mechanism == "" {
		c.Mechanism = "SCRAM-SHA-1"
	}
	blog.ErrorJSON("xxx c: %s, configmap: %s", c, configmap)

	maxOpenConns, err := strconv.ParseUint(configmap[prefix+".maxOpenConns"], 10, 64)
	if err != nil {
		blog.Errorf("parse mongo.maxOpenConns config error: %s, use default value: %d", err.Error(), DefaultMaxOpenConns)
		maxOpenConns = DefaultMaxOpenConns
	}
	if maxOpenConns > MaximumMaxOpenConns {
		blog.Errorf("mongo.maxOpenConns config %d exceeds maximum value, use maximum value %d", maxOpenConns, MaximumMaxOpenConns)
		maxOpenConns = MaximumMaxOpenConns
	}
	c.MaxOpenConns = maxOpenConns

	maxIdleConns, err := strconv.ParseUint(configmap[prefix+".maxIdleConns"], 10, 64)
	if err != nil || maxIdleConns < MinimumMaxIdleOpenConns {
		blog.Errorf("parse mongo.maxIdleConns config encounters error %v or %d less than minimum value, use minimum value %d", err, maxIdleConns, MinimumMaxIdleOpenConns)
		maxIdleConns = MinimumMaxIdleOpenConns
	}
	c.MaxIdleConns = maxIdleConns

	return c
}

func (c Config) GetMongoConf() local.MongoConf {
	return local.MongoConf{
		MaxOpenConns: c.MaxOpenConns,
		MaxIdleConns: c.MaxIdleConns,
		URI:          c.BuildURI(),
		RsName:       c.RsName,
	}
}

func (c Config) GetMongoClient() (db dal.RDB, err error) {
	mongoConf := local.MongoConf{
		MaxOpenConns: c.MaxOpenConns,
		MaxIdleConns: c.MaxIdleConns,
		URI:          c.BuildURI(),
		RsName:       c.RsName,
	}
	db, err = local.NewMgo(mongoConf, time.Minute)
	if err != nil {
		return nil, fmt.Errorf("connect mongo server failed %s", err.Error())
	}
	return
}
