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

	"configcenter/src/common/blog"
	"configcenter/src/common/registerdiscover"
	"configcenter/src/common/types"
)

type noticeHandler struct {
	regdiscv *registerdiscover.RegDiscv
	addrport string
}

func handleNotice(ctx context.Context, rd *registerdiscover.RegDiscv, addrport string) error {
	handler := &noticeHandler{
		regdiscv: rd,
		addrport: addrport,
	}
	return handler.handleLogNotice(ctx)
}

func (handler *noticeHandler) handleLogNotice(ctx context.Context) error {
	logVPath := fmt.Sprintf("%s/%s/%s/v", types.CC_SERVNOTICE_BASEPATH, "log", handler.addrport)
	ch, err := handler.regdiscv.Watch(ctx, logVPath)
	if err != nil {
		return err
	}

	go func () {
		handler.waitLogNoticeEvent(ctx, ch)
		// delete logv key when exit
		handler.regdiscv.Delete(logVPath)
	}()

	return nil
}

func (handler *noticeHandler) waitLogNoticeEvent(ctx context.Context, ch <-chan *registerdiscover.DiscoverEvent) {
	data := map[string]int32{
		"defaultV": blog.GetV(),
		"v":        blog.GetV(),
	}

	for {
		select {
		case event := <-ch:
			if event.Type == registerdiscover.EVENT_PUT {
				if err := json.Unmarshal([]byte(event.Value), &data); err != nil {
					blog.Errorf("fail to unmarshal data(%s), err: %s", event.Value, err.Error())
					continue
				}
				blog.SetV(data["v"])
			}
		case <-ctx.Done():
			blog.Infof("watch stopped because of context done")
			return
		}
	}
}