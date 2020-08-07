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

package esb

import (
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/thirdpartyclient/esbserver"
	"configcenter/src/thirdpartyclient/esbserver/esbutil"
)

var (
	esbClient esbserver.EsbClientInterface
	cfgChan   chan esbutil.EsbConfig = make(chan esbutil.EsbConfig, 10)

	lastInitErr   errors.CCErrorCoder
	lastConfigErr errors.CCErrorCoder
	tlsConfig     util.TLSClientConfig
)

func EsbClient() esbserver.EsbClientInterface {
	return esbClient
}

func ParseEsbConfig(config map[string]string) (*esbutil.EsbConfig, errors.CCErrorCoder) {
	esbAddr, addrOk := config["esb.addr"]
	if addrOk == false {
		blog.Infof("esb addr not found, unable to call esb service")
		lastConfigErr = errors.NewCCError(common.CCErrCommConfMissItem, "Configuration file missing [esb.addr] configuration item")
		return nil, lastConfigErr
	}
	esbAppCode, appCodeOk := config["esb.appCode"]
	if appCodeOk == false {
		blog.Errorf("esb appCode not found, unable to call esb service")
		lastConfigErr = errors.NewCCError(common.CCErrCommConfMissItem, "Configuration file missing [esb.esbAppCode] configuration item")
		return nil, lastConfigErr
	}
	esbAppSecret, appSecretOk := config["esb.appSecret"]
	if appSecretOk == false {
		blog.Errorf("esb appSecretOk not found,unable to call esb service")
		lastConfigErr = errors.NewCCError(common.CCErrCommConfMissItem, "Configuration file missing [esb.appSecret] configuration item")
		return nil, lastConfigErr
	}
	// 不支持热更新
	var err error
	tlsConfig, err = util.NewTLSClientConfigFromConfig("esb", config)
	if err != nil {
		lastInitErr = errors.NewCCError(common.CCErrCommResourceInitFailed, "'esb' initialization failed")
		return nil, lastInitErr
	}

	esbConfig := &esbutil.EsbConfig{
		Addrs:     esbAddr,
		AppCode:   esbAppCode,
		AppSecret: esbAppSecret,
	}
	return esbConfig, nil

}

func InitEsbClient(defaultCfg *esbutil.EsbConfig) errors.CCErrorCoder {

	apiMachineryConfig := &util.APIMachineryConfig{
		QPS:       1000,
		Burst:     1000,
		TLSConfig: &tlsConfig,
	}
	esbSrv, err := esbserver.NewEsb(apiMachineryConfig, cfgChan, defaultCfg, nil)
	if err != nil {
		blog.Errorf(" esbserve initialization error. err:%s", err.Error())
		lastInitErr = errors.NewCCError(common.CCErrCommResourceInitFailed, "'esb' initialization failed")
		return lastInitErr
	}
	esbClient = esbSrv
	return nil
}

func Validate() errors.CCErrorCoder {
	return nil
}

func UpdateEsbConfig(config esbutil.EsbConfig) {
	go func() {
		cfgChan <- config
	}()
}
