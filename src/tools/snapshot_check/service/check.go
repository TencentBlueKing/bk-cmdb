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

package service

import (
	"configcenter/src/common"
	"fmt"
	"time"

	"configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/backbone/service_mange/zk"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	ccRedis "configcenter/src/storage/dal/redis"

	"gopkg.in/redis.v5"
)

type Service struct {
	regdiscv        string
	strDefaultAppID string
	config          map[string]string
}

func NewService(regdiscv string, defaultAppID int) *Service {
	return &Service{
		regdiscv:        regdiscv,
		strDefaultAppID: fmt.Sprintf("%d", defaultAppID),
	}
}

func (s *Service) TriggerTicker(interval int) {
	s.triggerTicker(interval)
	timer := time.NewTicker(time.Duration(interval) * time.Minute)
	for range timer.C {
		s.triggerTicker(interval)
	}
}

func (s *Service) triggerTicker(interval int) {

	prefix := "\n\n\n\n\n\n"
	blog.Infof("%s=====================\nstart check  ", prefix)

	blog.Infof("start checkRegdiscv")
	if err := s.checkRegdiscv(); err != nil {
		blog.Errorf(err.Error())
	}

	blog.Infof("start checkConf")
	if err := s.checkConf(); err != nil {
		blog.Errorf(err.Error())
	}

	blog.Infof("start checkCCHostSnaphot")
	if err := s.checkCCHostSnaphot(); err != nil {
		blog.Errorf(err.Error())
	}

	blog.Infof("start checkHostSnapshot")
	if err := s.checkHostSnapshot(); err != nil {
		blog.Errorf(err.Error())
	}

	blog.Infof("\nend check  ")

}

func (s *Service) checkRegdiscv() error {
	client := zk.NewZkClient(s.regdiscv, 5*time.Second)
	if err := client.Start(); err != nil {
		return fmt.Errorf("connect regdiscv [%s] failed: %v", s.regdiscv, err)
	}
	return nil
}

func (s *Service) checkConf() error {
	client := zk.NewZkClient(s.regdiscv, 5*time.Second)
	if err := client.Start(); err != nil {
		return fmt.Errorf("connect regdiscv [%s] failed: %v", s.regdiscv, err)
	}
	if err := client.Ping(); err != nil {
		return fmt.Errorf("ping regdiscv [%s] failed: %v", s.regdiscv, err)
	}

	path := fmt.Sprintf("%s/%s", types.CC_SERVCONF_BASEPATH, types.CC_MODULE_DATACOLLECTION)
	strConf, err := client.Client().Get(path)
	if err != nil {
		return fmt.Errorf("get path [%s] from regdiscv [%s] failed: %v", path, s.regdiscv, err)
	}

	procConf, err := configcenter.ParseConfigWithData([]byte(strConf))
	if err != nil {
		return fmt.Errorf("get path [%s] from regdiscv [%s]  parse config failed: %v", path, s.regdiscv, err)
	}

	s.config = procConf.ConfigMap
	if len(s.config) == 0 {
		return fmt.Errorf("get path [%s] from regdiscv [%s]  parse config empty", path, s.regdiscv)
	}

	return nil
}

func (s *Service) checkCCHostSnaphot() error {

	redisConfig := ccRedis.ParseConfigFromKV("redis", s.config)
	client, err := ccRedis.NewFromConfig(redisConfig)
	if err != nil {
		return fmt.Errorf("connect redis [%s] failed: %s", redisConfig.Address, err.Error())
	}

	keys, err := client.Keys(common.RedisSnapKeyPrefix + "*").Result()
	if err != nil {
		return fmt.Errorf("execute keys command in redis [%s] failed: %s", redisConfig.Address, err.Error())
	}

	blog.Infof("checkCCHostSnaphost  has keys count:%s", len(keys))

	return nil
}

func (s *Service) checkHostSnapshot() error {

	redisConfig := ccRedis.ParseConfigFromKV("snap-redis", s.config)
	client, err := ccRedis.NewFromConfig(redisConfig)
	if err != nil {
		return fmt.Errorf("connect redis [%s] failed: %s", redisConfig.Address, err.Error())
	}

	channelArr := getSnapshotName(s.strDefaultAppID)
	sub, err := client.PSubscribe(channelArr...)
	if err != nil {
		return fmt.Errorf("subscribe channel [%#v] from redis [%s] failed: %s", channelArr, redisConfig.Address, err.Error())
	}

	stopChn := make(chan bool, 2)
	receiveMsgCount := 0
	timer := time.NewTimer(time.Second * 70)
	var receiveMsgErr error
	go func() {
		for len(stopChn) == 0 {
			received, err := sub.ReceiveTimeout(time.Second * 10)
			if err != nil {
				receiveMsgErr = fmt.Errorf("receive message from channel [%#v] in redis [%s] failed: %s", channelArr, redisConfig.Address, err.Error())
				return
			}
			msg, ok := received.(*redis.Message)
			if !ok {
				continue
			}
			if msg.Payload != "" {
				continue
			}

			receiveMsgCount++
		}
	}()

	<-timer.C
	if receiveMsgErr != nil {
		return receiveMsgErr
	}
	stopChn <- true

	blog.Infof("receive message from channel [%#v] of redis [%s] count: %d(1 minute total)", channelArr, redisConfig.Address, receiveMsgCount)
	if receiveMsgCount == 0 {
		return fmt.Errorf("not receive message from channel [%#v] of redis [%s]", channelArr, redisConfig.Address)
	}

	return nil
}

func getSnapshotName(strDefaultAppID string) []string {
	return []string{
		// 瘦身后的通道名
		"snapshot%d" + strDefaultAppID,
		// 瘦身前的通道名，为增加向前兼容的而订阅这个老通道
		strDefaultAppID, "_snapshot",
	}
}
