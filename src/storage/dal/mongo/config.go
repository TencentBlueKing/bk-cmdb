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
	"strings"
	"time"

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
	// if timeout isn't configured, use default value
	DefaultSocketTimeout = 10
	// if timeout exceeds maximum value, use maximum value
	MaximumSocketTimeout = 30
	// if timeout less than the minimum value, use minimum value
	MinimumSocketTimeout = 5
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
	SocketTimeout   int
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

func (c Config) GetMongoConf() local.MongoConf {
	return local.MongoConf{
		MaxOpenConns: c.MaxOpenConns,
		MaxIdleConns: c.MaxIdleConns,
		URI:          c.BuildURI(),
		RsName:       c.RsName,
		SocketTimeout:  c.SocketTimeout,
	}
}

func (c Config) GetMongoClient() (db dal.RDB, err error) {
	mongoConf := local.MongoConf{
		MaxOpenConns: c.MaxOpenConns,
		MaxIdleConns: c.MaxIdleConns,
		URI:          c.BuildURI(),
		RsName:       c.RsName,
		SocketTimeout: c.SocketTimeout,
	}
	db, err = local.NewMgo(mongoConf, time.Minute)
	if err != nil {
		return nil, fmt.Errorf("connect mongo server failed %s", err.Error())
	}
	return
}
