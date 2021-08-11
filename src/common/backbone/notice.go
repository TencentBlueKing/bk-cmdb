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

	"configcenter/src/common/backbone/service_mange"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
)

type noticeHandler struct {
	client   service_mange.ClientInterface
	addrport string
}

func handleNotice(ctx context.Context, client service_mange.ClientInterface, addrport string) error {
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
		defer handler.client.Delete(logVPath)
		for {
			val, err := handler.client.Get(logVPath)
			if err != nil {
				blog.Errorf("fail to get node(%s), err:%s\n", logVPath, err.Error())
				continue
			}
			if val == "" {
				data["v"] = blog.GetV()
				logVData, _ := json.Marshal(data)
				err = handler.client.Put(logVPath, string(logVData))
				if err != nil {
					blog.Errorf("fail to register node(%s), err:%s\n", logVPath, err.Error())
					continue
				}
			}
			err = json.Unmarshal([]byte(val), &data)
			if err != nil {
				blog.Errorf("fail to unmarshal data(%v), err:%s\n", val, err.Error())
				continue
			}
			blog.SetV(data["v"])
			select {
			case <-ctx.Done():
				blog.Warnf("log watch stopped because of context done.")
				_ = handler.client.Delete(logVPath)
				logPath := fmt.Sprintf("%s/%s/%s", types.CC_SERVNOTICE_BASEPATH, "log", handler.addrport)
				_ = handler.client.Delete(logPath)
				return
			default:
				blog.Infof("update log level, node(%s)\n", logVPath)
			}
		}
	}()
	return nil
}
