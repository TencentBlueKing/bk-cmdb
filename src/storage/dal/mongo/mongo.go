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
)

// Config config
type Config struct {
	Connect      string
	Address      string
	User         string
	Password     string
	Database     string
	Mechanism    string
	MaxOpenConns string
	MaxIdleConns string
}

// BuildURI return mongo uri according to  https://docs.mongodb.com/manual/reference/connection-string/
// format example: mongodb://[username:password@]host1[:port1][,host2[:port2],...[,hostN[:portN]]][/[database][?options]]
func (c Config) BuildURI() string {
	if c.Connect != "" {
		return c.Connect
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s", c.User, c.Password, c.Address, c.Database)
	if c.Mechanism != "" {
		uri += "?authMechanism=" + c.Mechanism
	}
	return uri
}
