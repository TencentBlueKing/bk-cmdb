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

package backbone

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/zkclient"

	"github.com/samuel/go-zookeeper/zk"
)

type noticeHandler struct {
	client   *zkclient.ZkClient
	addrport string
}

func handleNotice(ctx context.Context, client *zkclient.ZkClient, addrport string) error {
	handler := &noticeHandler{
		client:   client,
		addrport: addrport,
	}
	if err := handler.handleLogNotice(ctx); err != nil {
		return err
	}
	return nil
}

func (handler *noticeHandler) handleLogNotice(ctx context.Context) error {
	logVPath := fmt.Sprintf("%s/%s/%s/v", types.CC_SERVNOTICE_BASEPATH, "log", handler.addrport)
	data := map[string]int32{
		"defaultV": blog.GetV(),
		"v":        blog.GetV(),
	}
	go func() {
		defer handler.client.Del(logVPath, -1)
		var ch <-chan zk.Event
		var err error
		for {
			_, _, ch, err = handler.client.GetW(logVPath)
			if err != nil {
				blog.Errorf("log watch failed, will watch after 10s, path: %s, err: %s", logVPath, err.Error())
				switch err {
				case zk.ErrClosing, zk.ErrConnectionClosed:
					if conErr := handler.client.Connect(); conErr != nil {
						blog.Errorf("fail to watch register node(%s), reason: connect closed. retry connect err:%s\n", logVPath, conErr.Error())
						time.Sleep(10 * time.Second)
					}
				case zk.ErrNoNode:
					data["v"] = blog.GetV()
					logVData, _ := json.Marshal(data)
					err = handler.client.CreateDeepNode(logVPath, logVData)
					if err != nil {
						blog.Errorf("fail to register node(%s), err:%s\n", logVPath, err.Error())
					}
				}
				continue
			}
			select {
			case event := <-ch:
				if event.Type != zk.EventNodeDataChanged {
					continue
				}
				dat, err := handler.client.Get(logVPath)
				if err != nil {
					blog.Errorf("fail to get node(%s), err:%s\n", logVPath, err.Error())
					continue
				}
				err = json.Unmarshal([]byte(dat), &data)
				if err != nil {
					blog.Errorf("fail to unmarshal data(%v), err:%s\n", dat, err.Error())
					continue
				}
				blog.SetV(data["v"])
			case <-ctx.Done():
				blog.Warnf("log watch stopped because of context done.")
				_ = handler.client.Del(logVPath, -1)
				logPath := fmt.Sprintf("%s/%s/%s", types.CC_SERVNOTICE_BASEPATH, "log", handler.addrport)
				_ = handler.client.Del(logPath, -1)
				return
			}
		}
	}()
	return nil
}
