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
	"time"

	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/types"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/tools/cmdb_ctl/app/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewRedisOperationCommand())
}

const (
	redisDelBatchDefaulNum = 10
)

type redisOperation struct {
	cursor  uint64
	match   string
	count   int64
	service *config.Service
	zkAddr  string
}

func NewRedisOperationCommand() *cobra.Command {

	conf := new(redisOperation)
	cmd := &cobra.Command{
		Use:   "redis",
		Short: "redis  operations",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	scanCmds := make([]*cobra.Command, 0)
	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "scan the redis",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRedisScan(conf)
		},
	}
	scanCmd.Flags().Uint64Var(&conf.cursor, "cursor", 0, "redis scan cursor default is 0")
	scanCmd.Flags().StringVar(&conf.match, "match", "", "redis scan  match pattern  default is null")
	scanCmd.Flags().Int64Var(&conf.count, "count", 10, "redis scan count default value is 10")
	scanCmd.Flags().StringVar(&conf.zkAddr, "zk-addr", "127.0.0.1:2181", "zk address where the redis configuration file is stored")

	scanCmds = append(scanCmds, scanCmd)
	for _, fCmd := range scanCmds {
		cmd.AddCommand(fCmd)
	}

	delCmds := make([]*cobra.Command, 0)
	delCmd := &cobra.Command{
		Use:   "scan-del",
		Short: "del scanned keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRedisScanDel(conf)
		},
	}

	delCmd.Flags().Uint64Var(&conf.cursor, "cursor", 0, "redis scan cursor default is 0")
	delCmd.Flags().StringVar(&conf.match, "match", "", "redis scan  match pattern  default is null")
	delCmd.Flags().Int64Var(&conf.count, "count", 10, "redis scan count default value is 10")
	delCmd.Flags().StringVar(&conf.zkAddr, "zk-addr", "127.0.0.1:2181", "zk address where the redis configuration file is stored")

	delCmds = append(delCmds, delCmd)
	for _, fCmd := range delCmds {
		cmd.AddCommand(fCmd)
	}

	return cmd
}

func (s *redisOperation) setRedisConf() error {

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

func newService(zkaddr string) (*redisOperation, error) {
	service, err := config.NewZkService(zkaddr)
	if err != nil {
		return nil, err
	}
	return &redisOperation{
		service: service,
	}, nil
}

func runRedisScan(conf *redisOperation) error {

	s, err := newService(conf.zkAddr)
	if err != nil {
		fmt.Printf("connect zk fail, err :%v .\n", err)
		return err
	}
	err = s.setRedisConf()
	if err != nil {
		fmt.Printf("set Redis config  fail, err :%v .\n", err)
		return err
	}
	redisConfig, err := cc.Redis("redis")
	if err != nil {
		fmt.Printf("connect redis fail, err :%v .\n", err)
		return err
	}

	redisCli, err := redis.NewFromConfig(redisConfig)
	if err != nil {
		fmt.Printf("read redis config fail  err :%v.\n", err)
		return err
	}

	res := redisCli.Scan(context.Background(), conf.cursor, conf.match, conf.count)
	keys, cursor := res.Val()

	fmt.Printf("keys is begin :\n")
	for _, v := range keys {
		fmt.Printf("%v \n", v)
	}
	fmt.Printf("keys is end\n")

	fmt.Printf("cursor is %d \n", cursor)

	return nil
}

func runRedisScanDel(conf *redisOperation) error {

	var (
		start   int
		keysTmp []string
		bFlag   bool
	)

	s, err := newService(conf.zkAddr)
	if err != nil {
		fmt.Printf("connect zk fail, err :%v .\n", err)
		return err
	}
	err = s.setRedisConf()
	if err != nil {
		fmt.Printf("set Redis config  fail, err :%v .\n", err)
		return err
	}

	redisConfig, err := cc.Redis("redis")
	if err != nil {
		fmt.Printf("connect redis fail, err :%v .\n", err)
		return err
	}

	redisCli, err := redis.NewFromConfig(redisConfig)
	if err != nil {
		fmt.Printf("read redis config fail, err :%v.\n", err)
		return err
	}
	ctx := context.Background()
	res := redisCli.Scan(ctx, conf.cursor, conf.match, conf.count)
	keys, cursor := res.Val()

	keysNum := len(keys)
	for {

		if start >= keysNum {
			break
		}
		if start+redisDelBatchDefaulNum > keysNum {
			keysTmp = keys[start:]
		} else {
			keysTmp = keys[start : start+redisDelBatchDefaulNum]
		}

		err = redisCli.Del(ctx, keysTmp...).Err()
		if err != nil {
			fmt.Printf("del the keys fail err:%v \n", err)
			bFlag = true
			break
		}

		time.Sleep(100 * time.Millisecond)
		start = start + redisDelBatchDefaulNum
	}

	if bFlag {
		fmt.Printf("del the keys fail err:%v \n", err)
		return err
	} else {

		fmt.Printf("keys is begin :\n")
		for _, v := range keys {
			fmt.Printf("%v \n", v)
		}
		fmt.Printf("keys is end\n")

		fmt.Printf("del keys success ,total num is %v\n", len(keys))
		fmt.Printf("cursor is %d \n", cursor)
	}

	return nil
}
