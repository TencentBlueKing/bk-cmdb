/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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

package etcd_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	cc "github.com/TencentBlueKing/bk-cmdb/pkg/config-center"
	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/etcd"
	etcdcli "github.com/TencentBlueKing/bk-cmdb/pkg/etcd"
)

func TestEtcdConfigCenter(t *testing.T) {
	ctx := context.Background()

	tempDir := prepareConfigFiles(t)
	registry, discovery := generateRegDisc(t)

	// test register event handler for writer
	pgsql := testRegisterEventHandler[pgsqlConf](t, "pgsql")

	// test write initial config
	writer := cc.NewWriter(registry, tempDir)
	if err := writer.RunConfigWrite(ctx); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, pgsql, &pgsqlConf{
		Host:     "pgsql.example.com",
		Port:     5432,
		Database: "test",
	})

	// test writer update config
	writeConfigFile(t, tempDir, "pgsql", `
pgsql:
  host: pgsql1.example.com
  port: 15432
  database: test1
`)
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, pgsql, &pgsqlConf{
		Host:     "pgsql1.example.com",
		Port:     15432,
		Database: "test1",
	})

	// test register event handler for reader
	apigw := testRegisterEventHandler[apigwConf](t, "apigw")
	testConf := ""
	err := cc.RegisterEventHandler("test", func(event *cc.Event) error {
		switch event.Type {
		case cc.UpsertEvent:
			var err error
			testConf, err = cc.ConvertBasic[string](event.Data)
			if err != nil {
				t.Fatal(err)
			}
		case cc.DeleteEvent:
			testConf = ""
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// test read initial config
	reader := cc.NewReader(discovery)
	if err := reader.RunConfigRead(ctx); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, apigw, &apigwConf{
		Address:   "https://apigw.example.com",
		AppCode:   "test",
		AppSecret: "test",
	})
	assert.Equal(t, testConf, "test")

	// test reader update config
	writeConfigFile(t, tempDir, "common", `
apigw:
  address: https://apigw1.example.com
  app_code: test1
  app_secret: test1

test: test1
`)
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, apigw, &apigwConf{
		Address:   "https://apigw1.example.com",
		AppCode:   "test1",
		AppSecret: "test1",
	})
	assert.Equal(t, testConf, "test1")

	// test reader delete config
	writeConfigFile(t, tempDir, "common", `
apigw:
  address: https://apigw2.example.com
  app_code: test2
  app_secret: test2
`)
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, apigw, &apigwConf{
		Address:   "https://apigw2.example.com",
		AppCode:   "test2",
		AppSecret: "test2",
	})
	assert.Equal(t, testConf, "")
}

// prepareConfigFiles prepare config files
func prepareConfigFiles(t *testing.T) string {
	tempDir := t.TempDir()

	writeConfigFile(t, tempDir, "pgsql", `
pgsql:
  host: pgsql.example.com
  port: 5432
  database: test
`)

	writeConfigFile(t, tempDir, "common", `
apigw:
  address: https://apigw.example.com
  app_code: test
  app_secret: test

test: test
`)

	writeConfigFile(t, tempDir, "redis", "")
	writeConfigFile(t, tempDir, "extra", "")

	return tempDir
}

func writeConfigFile(t *testing.T, dir, conf string, data string) {
	file := filepath.Join(dir, conf+".yaml")
	if err := os.WriteFile(file, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
}

type pgsqlConf struct {
	Host     string `json:"host"`
	Port     int64  `json:"port"`
	Database string `json:"database"`
}

type apigwConf struct {
	Address   string `json:"address"`
	AppCode   string `json:"app_code"`
	AppSecret string `json:"app_secret"`
}

func generateRegDisc(t *testing.T) (cc.Registry, cc.Discovery) {
	// create a new etcd client
	// NOTE: etcdConf is the etcd config for test, change it to your own etcd config.
	etcdConf := &etcdcli.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		Username:  "",
		Password:  "",
		TLS:       nil,
	}

	etcdCli, err := etcdcli.New(etcdConf)
	if err != nil {
		t.Fatal(err)
	}

	registry, err := etcd.NewRegistry(etcdCli)
	if err != nil {
		t.Fatal(err)
	}

	discovery, err := etcd.NewDiscovery(etcdCli)
	if err != nil {
		t.Fatal(err)
	}

	return registry, discovery
}

func testRegisterEventHandler[T any](t *testing.T, key string) *T {
	conf := new(T)
	err := cc.RegisterEventHandler(key, func(event *cc.Event) error {
		switch event.Type {
		case cc.UpsertEvent:
			upsertData, err := cc.Convert[T](event.Data)
			if err != nil {
				t.Fatal(err)
			}
			*conf = *upsertData
		case cc.DeleteEvent:
			var empty T
			*conf = empty
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return conf
}
