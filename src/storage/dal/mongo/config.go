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

	"configcenter/src/common/backbone"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
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
	MaxOpenConns string
	MaxIdleConns string
	TxnEnabled    string
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
	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s", c.User, c.Password, c.Address, c.Database)
	if c.Mechanism != "" {
		uri += "?authMechanism=" + c.Mechanism
	}
	return uri
}

func (c Config) GetMaxOpenConns() int {
	max, err := strconv.Atoi(c.MaxOpenConns)
	if err != nil {
		return 0
	}
	return max
}

func (c Config) GetMaxIdleConns() int {
	max, err := strconv.Atoi(c.MaxIdleConns)
	if err != nil {
		return 0
	}
	return max
}

// ParseConfigFromKV returns a new config
func ParseConfigFromKV(prefix string, conifgmap map[string]string) Config {
	return Config{
		Address:      conifgmap[prefix+".host"],
		Port:         conifgmap[prefix+".port"],
		User:         conifgmap[prefix+".usr"],
		Password:     conifgmap[prefix+".pwd"],
		Database:     conifgmap[prefix+".database"],
		MaxOpenConns: conifgmap[prefix+".maxOpenConns"],
		MaxIdleConns: conifgmap[prefix+".maxIDleConns"],
		Mechanism:    conifgmap[prefix+".mechanism"],
		TxnEnabled:    conifgmap[prefix+".txnEnabled"],
	}
}

func (c Config) GetMongoClient(engine *backbone.Engine) (db dal.RDB, err error) {
	db, err = local.NewMgo(c.BuildURI(), time.Minute)
	if err != nil {
		return nil, fmt.Errorf("connect mongo server failed %s", err.Error())
	}
	return
}

func (c Config) GetTransactionClient(engine *backbone.Engine) (client dal.Transaction, err error) {
	client, err = local.NewMgo(c.BuildURI(), time.Minute)

	if err != nil {
		return nil, fmt.Errorf("connect mongo server failed %s", err.Error())
	}
	return
}
