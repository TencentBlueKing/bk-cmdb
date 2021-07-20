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

package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/types"
	ccRedis "configcenter/src/storage/dal/redis"
	"configcenter/src/tools/cmdb_ctl/app/config"

	rawRedis "github.com/go-redis/redis/v7"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewSnapshotCheckCommand())
}

type snapshotCheckConf struct {
	bizID int
}

func NewSnapshotCheckCommand() *cobra.Command {
	conf := new(snapshotCheckConf)

	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "check host snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSnapshotCheck(conf)
		},
	}

	conf.addFlags(cmd)

	return cmd
}

func (c *snapshotCheckConf) addFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&c.bizID, "bizId", 2, "blueking business id. e.g: 2")
}

type snapshotCheckService struct {
	service *config.Service
	bizID   string
	config  map[string]string
}

func newSnapshotCheckService(zkaddr string, bizID int) (*snapshotCheckService, error) {
	service, err := config.NewZkService(zkaddr)
	if err != nil {
		return nil, err
	}
	return &snapshotCheckService{
		service: service,
		bizID:   strconv.Itoa(bizID),
	}, nil
}

func runSnapshotCheck(c *snapshotCheckConf) error {
	srv, err := newSnapshotCheckService(config.Conf.ZkAddr, c.bizID)
	if err != nil {
		return err
	}
	return srv.snapshotCheck()
}

func (s *snapshotCheckService) snapshotCheck() error {
	fmt.Println("=====================\nstart check")
	fmt.Println("start checkConf")
	if err := s.checkConf(); err != nil {
		return err
	}

	fmt.Println("start checkCCHostSnaphot")
	if err := s.checkCCHostSnaphot(); err != nil {
		return err
	}

	fmt.Println("start checkHostSnapshot")
	if err := s.checkHostSnapshot(); err != nil {
		return err
	}

	fmt.Println("end check")
	return nil
}

func (s *snapshotCheckService) checkConf() error {
	if err := s.service.ZkCli.Ping(); err != nil {
		if err = s.service.ZkCli.Connect(); err != nil {
			return err
		}
	}

	path := fmt.Sprintf("%s/%s", types.CC_SERVCONF_BASEPATH, types.CCConfigureRedis)
	strConf, err := s.service.ZkCli.Get(path)
	if err != nil {
		return fmt.Errorf("get path [%s] from zk [%v] failed: %v", path, s.service.ZkCli.ZkHost, err)
	}

	if err := cc.SetRedisFromByte([]byte(strConf)); err != nil {
		return fmt.Errorf("get path [%s] from regdiscv [%v]  parse config failed: %v", path, s.service.ZkCli.ZkHost, err)
	}

	if len(strConf) == 0 {
		return fmt.Errorf("get path [%s] from regdiscv [%v]  parse config empty", path, s.service.ZkCli.ZkHost)
	}

	return nil
}

func (s *snapshotCheckService) checkCCHostSnaphot() error {

	redisConfig, err := cc.Redis("redis")
	if err != nil {
		return err
	}
	client, err := ccRedis.NewFromConfig(redisConfig)
	if err != nil {
		return fmt.Errorf("connect redis [%s] failed: %s", redisConfig.Address, err.Error())
	}

	keys, err := client.Keys(context.Background(), common.RedisSnapKeyPrefix+"*").Result()
	if err != nil {
		return fmt.Errorf("execute keys command in redis [%s] failed: %s", redisConfig.Address, err.Error())
	}

	fmt.Printf("checkCCHostSnaphost has keys count: %d\n", len(keys))

	return nil
}

func (s *snapshotCheckService) checkHostSnapshot() error {

	redisConfig, err := cc.Redis("redis.snap")
	if err != nil {
		return err
	}
	client, err := ccRedis.NewFromConfig(redisConfig)
	if err != nil {
		return fmt.Errorf("connect redis [%s] failed: %s", redisConfig.Address, err.Error())
	}

	channelArr := getSnapshotName(s.bizID)
	sub := client.PSubscribe(context.Background(), channelArr...)

	stopChn := make(chan bool, 2)
	receiveMsgCount := 0
	timer := time.NewTimer(time.Minute * 2)
	var receiveMsgErr error
	go func() {
		for len(stopChn) == 0 {
			received, err := sub.ReceiveTimeout(time.Minute * 1)
			if err != nil {
				receiveMsgErr = fmt.Errorf("receive message from channel [%#v] in redis [%s] failed: %s", channelArr, redisConfig.Address, err.Error())
				return
			}
			msg, ok := received.(*rawRedis.Message)
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

	fmt.Printf("receive message from channel [%#v] of redis [%s] count: %d(1 minute total)\n", channelArr, redisConfig.Address, receiveMsgCount)
	if receiveMsgCount == 0 {
		return fmt.Errorf("not receive message from channel [%#v] of redis [%s]", channelArr, redisConfig.Address)
	}

	return nil
}

func getSnapshotName(strDefaultAppID string) []string {
	return []string{
		// 瘦身后的通道名
		"snapshot" + strDefaultAppID,
		// 瘦身前的通道名，为增加向前兼容的而订阅这个老通道
		strDefaultAppID + "_snapshot",
	}
}
