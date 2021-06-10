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

	"configcenter/src/storage/dal/redis"
	"configcenter/src/tools/cmdb_ctl/app/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewRedisOperationCommand())
}

type redisOperation struct {
	cursor uint64
	match  string
	count  int64
	keys   []string
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
		Short: "redis scan operation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRedisScan(conf)
		},
	}
	scanCmd.Flags().Uint64Var(&conf.cursor, "cursor", 0, "redis scan cursor default is 0")
	scanCmd.Flags().StringVar(&conf.match, "match", "", "redis scan  match pattern, default is null")
	scanCmd.Flags().Int64Var(&conf.count, "count", 10, "redis scan count, default value is 10")

	scanCmds = append(scanCmds, scanCmd)
	for _, fCmd := range scanCmds {
		cmd.AddCommand(fCmd)
	}

	delCmds := make([]*cobra.Command, 0)
	delCmd := &cobra.Command{
		Use:   "del",
		Short: "redis del operation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRedisDel(conf)
		},
	}

	delCmd.Flags().StringSliceVar(&conf.keys, "keys", []string{}, "del keys,the parameter must be assigned ")

	delCmds = append(delCmds, delCmd)
	for _, fCmd := range delCmds {
		cmd.AddCommand(fCmd)
	}

	return cmd
}

func runRedisScan(conf *redisOperation) error {

	redisCfg := redis.Config{
		Address:  config.Conf.RedisAddr + ":" + config.Conf.RedisPort,
		Password: config.Conf.RedisPwd,
		Database: config.Conf.RedisDatabase,
	}

	redisCli, err := redis.NewFromConfig(redisCfg)
	if err != nil {
		fmt.Printf("connect redis fail  err :%v.\n", err)
		return err
	}

	ctx := context.Background()
	res := redisCli.Scan(ctx, conf.cursor, conf.match, conf.count)
	keys, cursor := res.Val()

	fmt.Printf("keys is begin:\n")
	for _, v := range keys {
		fmt.Printf("%v \n", v)
	}
	fmt.Printf("keys is end \n")

	fmt.Printf("cursor is %d \n", cursor)

	return nil
}
func runRedisDel(conf *redisOperation) error {

	redisCfg := redis.Config{
		Address:  config.Conf.RedisAddr + ":" + config.Conf.RedisPort,
		Password: config.Conf.RedisPwd,
		Database: config.Conf.RedisDatabase,
	}

	redisCli, err := redis.NewFromConfig(redisCfg)
	if err != nil {
		fmt.Printf("connect redis fail err :%v.\n", err)
		return err
	}

	err = redisCli.Del(context.Background(), conf.keys...).Err()
	if err != nil {
		fmt.Printf("del the keys fail err:%v \n", err)
		return err
	}

	fmt.Printf("del keys success !\n")
	return nil
}
