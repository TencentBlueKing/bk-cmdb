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

package tests

import (
	"encoding/json/v2"
	"fmt"
	"os"
	"testing"

	clientv3 "go.etcd.io/etcd/client/v3"

	etcdcli "github.com/TencentBlueKing/bk-cmdb/pkg/client/etcd"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

// GetTestEtcd create a new etcd client for test from env TEST_ETCD_CONF, should only be used in test code
func GetTestEtcd(t testing.TB) (*clientv3.Client, error) {
	var config = os.Getenv("TEST_ETCD_CONF")
	if config == "" {
		log.Warn(t.Context(), "TEST_ETCD_CONF is not set, skip", "test", t.Name())
		t.Skip("TEST_ETCD_CONF is not set")
		return nil, fmt.Errorf("TEST_ETCD_CONF is not set")
	}

	etcdConf := new(etcdcli.Config)
	err := json.Unmarshal([]byte(config), etcdConf)
	if err != nil {
		return nil, fmt.Errorf("unmarshal etcd config %s failed, %w", config, err)
	}

	etcdCli, err := etcdcli.New(etcdConf)
	if err != nil {
		return nil, fmt.Errorf("create etcd client failed, %w", err)
	}

	_, err = etcdCli.Delete(t.Context(), "/cc", clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("clean etcd failed, %w", err)
	}
	return etcdCli, nil
}
