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

package middleware

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/ac/extensions"
	bkc "configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/storage/dal/redis"
)

type Discover struct {
	ctx        context.Context
	httpHeader http.Header

	redisCli redis.Client
	*backbone.Engine
	authManager *extensions.AuthManager
}

var msgHandlerCnt = int64(0)

func NewDiscover(ctx context.Context, redisCli redis.Client, backbone *backbone.Engine, authManager *extensions.AuthManager) *Discover {
	header := http.Header{}
	header.Add(bkc.BKHTTPOwnerID, bkc.BKDefaultOwnerID)
	header.Add(bkc.BKHTTPHeaderUser, bkc.CCSystemCollectorUserName)

	discover := &Discover{
		redisCli:    redisCli,
		ctx:         ctx,
		httpHeader:  header,
		authManager: authManager,
	}
	discover.Engine = backbone
	return discover
}

// Hash returns hash value base on message.
func (d *Discover) Hash(cloudid, ip string) (string, error) {
	if len(cloudid) == 0 {
		return "", fmt.Errorf("can't make hash from invalid message format, cloudid empty")
	}
	if len(ip) == 0 {
		return "", fmt.Errorf("can't make hash from invalid message format, ip empty")
	}

	hash := fmt.Sprintf("%s:%s", cloudid, ip)

	return hash, nil
}

// Mock returns local mock message for testing.
func (d *Discover) Mock() string {
	return MockMessage
}

func (d *Discover) Analyze(msg *string) error {
	err := d.UpdateOrCreateInst(msg)
	if err != nil {
		return fmt.Errorf("create inst err: %v, raw: %s", err, msg)
	}
	return nil
}

var MockMessage = `{
    "meta": {
        "model": {
            "bk_obj_id": "bk_apache",
            "bk_supplier_account": "0"
        }
    },
    "data": {
        "bk_inst_name": "apache",
        "bk_ip": "192.168.0.1"
    }
}`
