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

package host

import (
	"context"
	"errors"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/cache/tools"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/reflector"
	"gopkg.in/redis.v5"
)

type Client struct {
	rds   *redis.Client
	event reflector.Interface
	db    dal.DB
	lock  tools.RefreshingLock
}

// GetHostWithID get host with host id.
// fields allows you can specify which fields you need only.
func (c *Client) GetHostWithID(ctx context.Context, opt *metadata.SearchHostWithIDOption) (string, error) {
	rid := ctx.Value(common.ContextRequestIDField)
	needRefresh := false
	data, err := c.rds.Get(hostKey.HostDetailKey(opt.HostID)).Result()
	if err != nil {
		if err != redis.Nil {
			// return directly to avoid cache penetration
			blog.Errorf("get host: %d from redis failed, err: %v, rid: %s", opt.HostID, err, rid)
			return "", err
		}
		// do not exist in cache, need to refresh from db.
		needRefresh = true
	}

	if !needRefresh {
		// already get the data, use cache
		if len(opt.Fields) == 0 {
			return data, nil
		}
		return *json.CutJsonDataWithFields(&data, opt.Fields), nil
	}

	// data has already expired, need to refresh from db
	ips, cloudID, detail, err := getHostDetailsFromMongoWithHostID(c.db, opt.HostID)
	if err != nil {
		blog.Errorf("get host with id: %d, and cache expired, but get from mongo failed, err: %v, rid: %s", opt.HostID, err, rid)
		return "", err
	}

	// try refresh cache
	c.tryRefreshHostDetail(opt.HostID, ips, cloudID, detail)

	if len(opt.Fields) == 0 {
		return string(detail), nil
	} else {
		h := string(detail)
		return *json.CutJsonDataWithFields(&h, opt.Fields), nil
	}
}

// GetHostWithInnerIP is to get host with the ip and cloud id it belongs.
// the ip must be a unique one, can not be a ip string with multiple ip separated with comma.
func (c *Client) GetHostWithInnerIP(ctx context.Context, opt *metadata.SearchHostWithInnerIPOption) (string, error) {
	rid := ctx.Value(common.ContextRequestIDField)
	if len(opt.InnerIP) == 0 || len(strings.Split(opt.InnerIP, ",")) > 1 {
		return "", errors.New("invalid ip address with multiple ip")
	}

	detail, err := c.getHostDetailWithIP(opt.InnerIP, opt.CloudID)
	if err != nil {
		blog.Errorf("get host with inner ip: %s failed, err：%v, rid: %s", opt.InnerIP, err, rid)
		return "", err
	}

	if len(opt.Fields) == 0 {
		return *detail, nil
	} else {
		return *json.CutJsonDataWithFields(detail, opt.Fields), nil
	}
}
