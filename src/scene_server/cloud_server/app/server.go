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

package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/cryptor"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/cloud_server/app/options"
	"configcenter/src/scene_server/cloud_server/cloudsync"
	"configcenter/src/scene_server/cloud_server/logics"
	svc "configcenter/src/scene_server/cloud_server/service"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdpartyclient/secrets"
)

func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	service := svc.NewService(ctx)

	process := new(CloudServer)
	input := &backbone.BackboneParameter{
		ConfigUpdate: process.onCloudConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     op.ServConf.RegDiscover,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	service.Engine = engine
	process.Core = engine
	process.Service = service

	configReady := false
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		if nil != process.Config {
			configReady = true
			break
		}
		blog.Infof("waiting for config ready ...")
		time.Sleep(time.Second)
	}
	if false == configReady {
		blog.Infof("waiting config timeout.")
		return errors.New("configuration item not found")
	}

	mongoConfig, err := engine.WithMongo()
	if nil != err {
		blog.Errorf("get mongo conf failed, err: %s", err.Error())
		return err
	}

	mongoConf := mongoConfig.GetMongoConf()
	db, err := local.NewMgo(mongoConf, time.Minute)
	if err != nil {
		return fmt.Errorf("connect mongo server failed, err: %s", err.Error())
	}
	process.Service.SetDB(db)

	redisConf, err := engine.WithRedis()
	if nil != err {
		blog.Errorf("get redis conf failed, err: %s", err.Error())
		return err
	}

	cache, err := redis.NewFromConfig(redisConf)
	if err != nil {
		return fmt.Errorf("connect redis server failed, err: %s", err.Error())
	}
	process.Service.SetCache(cache)

	blog.Infof("enable auth center: %v", auth.IsAuthed())

	var accountCryptor cryptor.Cryptor
	blog.Infof("enable cryptor: %v", op.EnableCryptor)
	if op.EnableCryptor == true {
		secretKey, err := process.getSecretKey()
		if err != nil {
			blog.Errorf("getSecretKey failed, err: %s", err.Error())
			return err
		}
		accountCryptor = cryptor.NewAesEncrpytor(secretKey)
	}

	process.Service.SetEncryptor(accountCryptor)
	process.Service.EnableTxn = op.EnableTxn

	process.Service.Logics = logics.NewLogics(service.Engine, db, cache, accountCryptor)

	syncConf := cloudsync.SyncConf{
		ZKClient:  service.Engine.ServiceManageClient().Client(),
		Logics:    process.Service.Logics,
		AddrPort:  input.SrvInfo.Instance(),
		MongoConf: mongoConf,
	}
	err = cloudsync.CloudSync(&syncConf)
	if err != nil {
		return fmt.Errorf("ProcessTask failed: %v", err)
	}

	err = backbone.StartServer(ctx, cancel, engine, service.WebService(), true)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		blog.Infof("process will exit!")
	}

	return nil
}

type CloudServer struct {
	Core    *backbone.Engine
	Config  *options.Config
	Service *svc.Service
}

func (c *CloudServer) onCloudConfigUpdate(previous, current cc.ProcessConfig) {
	if c.Config == nil {
		c.Config = new(options.Config)
	}
	c.Config.SecretKeyUrl = current.ConfigMap["cryptor.secret_key_url"]
	c.Config.SecretsAddrs = current.ConfigMap["cryptor.secrets_addrs"]
	c.Config.SecretsToken = current.ConfigMap["cryptor.secrets_token"]
	c.Config.SecretsProject = current.ConfigMap["cryptor.secrets_project"]
	c.Config.SecretsEnv = current.ConfigMap["cryptor.secrets_env"]
}

// getSecretKey get the secret key from bk-secrets service
func (c *CloudServer) getSecretKey() (string, error) {
	if c.Config.SecretKeyUrl == "" {
		return "", errors.New("config cryptor.secret_key_url is not set")
	}

	if c.Config.SecretsAddrs == "" {
		return "", errors.New("config cryptor.secrets_addrs is not set")
	}

	if c.Config.SecretsToken == "" {
		return "", errors.New("config cryptor.secrets_token is not set")
	}

	if c.Config.SecretsProject == "" {
		return "", errors.New("config cryptor.secrets_project is not set")
	}

	if c.Config.SecretsEnv == "" {
		return "", errors.New("config cryptor.secrets_env is not set")
	}

	secretsConfig := secrets.SecretsConfig{
		SecretKeyUrl:   c.Config.SecretKeyUrl,
		SecretsAddrs:   c.Config.SecretsAddrs,
		SecretsToken:   c.Config.SecretsToken,
		SecretsProject: c.Config.SecretsProject,
		SecretsEnv:     c.Config.SecretsEnv,
	}

	secretsClient, err := secrets.NewSecretsClient(nil, secretsConfig, nil)
	if err != nil {
		blog.Errorf("NewSecretsClient err:%s", err.Error())
		return "", err
	}

	return secretsClient.GetCloudAccountSecretKey(context.Background(), http.Header{})
}
