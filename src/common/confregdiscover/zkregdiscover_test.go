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
 
package confregdiscover

import (
	"testing"
	"time"
)

const (
	confPath = "/cc/config/test"
)

func TestWriteAndDiscover(t *testing.T) {

	zkRegDcv := NewZkRegDiscover("127.0.0.1:2181", time.Second*5)
	zkRegDcv.Start()

	// discover
	env, err := zkRegDcv.Discover(confPath)
	if err != nil {
		t.Errorf("fail to discover config path. err:%s", err.Error())
		zkRegDcv.Stop()
		return
	}

	if err := zkRegDcv.Write(confPath, []byte("conf")); err != nil {
		t.Errorf("fail to write config to zkRegDcv. err:%s", err.Error())
		zkRegDcv.Stop()
		return
	}

	for {
		select {
		case confEnv := <-env:
			t.Logf("config has changed. key:%s, data:%s", confEnv.Key, string(confEnv.Data))
		}
	}
}
