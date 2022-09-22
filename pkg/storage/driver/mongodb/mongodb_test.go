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

package mongodb

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMonogodb(t *testing.T) {

	prefix := "mongodb"
	configMap := map[string]string{
		prefix + ".host":         os.Getenv("mongo_host"),
		prefix + ".usr":          os.Getenv("mongo_usr"),
		prefix + ".pwd":          os.Getenv("mongo_pwd"),
		prefix + ".database":     os.Getenv("mongo_db"),
		prefix + ".mechanism":    os.Getenv("mongo_mechanism=SCRAM-SHA-1"),
		prefix + ".maxOpenConns": "300",
		prefix + ".maxIDleConns": "100",
	}
	config, err := ParseConfig(prefix, configMap)
	require.NoError(t, err)

	err = InitClient(prefix, config)
	require.NoError(t, err)

	dbErr := Client().Ping()
	require.NoError(t, dbErr)
	_, dbErr = Client().Table("tmptest").Find(nil).Count(context.Background())
	require.NoError(t, dbErr)
}
