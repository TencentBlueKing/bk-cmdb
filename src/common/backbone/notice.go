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
	"github.com/samuel/go-zookeeper/zk"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/zkclient"
)

type noticeHandler struct {
	client *zkclient.ZkClient
	addrport string
}

func handleNotice(ctx context.Context, client *zkclient.ZkClient, addrport string) error {
	handler := &noticeHandler{
		client: client,
		addrport: addrport,
	}
	err := handler.handleLogNotice(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (handler *noticeHandler)handleLogNotice(ctx context.Context) error {
	logVPath := fmt.Sprintf("%s/%s/%s/v", types.CC_SERVNOTICE_BASEPATH, "log", handler.addrport)
	data := map[string]int32 {
		"defaultV": blog.GetV(),
		"v": blog.GetV(),
	}
	logVData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = handler.client.CreateDeepNode(logVPath, logVData)
	if err != nil {
		return err
	}
	go func() {
		var ch <-chan zk.Event
		for {
				_, _, ch, err = handler.client.GetW(logVPath)
				if err != nil {
					blog.Errorf("log watch failed, will watch after 60s, path: %s, err: %s", logVPath, err.Error())
					time.Sleep(10 * time.Second)
					continue
				}
				break
		}
		for {
			select {
			case event := <-ch:
				logVData, _, ch, err = handler.client.GetW(logVPath)
				if err != nil {
					blog.Errorf("log watch failed, will watch after 60s, path: %s, err: %s", logVPath, err.Error())
					time.Sleep(10 * time.Second)
					continue
				}
				if event.Type != zk.EventNodeDataChanged {
					continue
				}
				err = json.Unmarshal(logVData, &data)
				if err == nil {
					blog.SetV(data["v"])
				}
			case <-ctx.Done():
				blog.Warnf("log watch stopped because of context done.")
				return
			}
		}
	}()
	return nil
}