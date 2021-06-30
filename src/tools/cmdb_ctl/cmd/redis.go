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

	"configcenter/src/storage/dal/redis"
	"configcenter/src/tools/cmdb_ctl/app/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewRedisOperationCommand())
}

const (
	redisDelBatchDefaulNum        = 100
	redisDefaultConnNum           = 10
	redisDefaultCursor     uint64 = 0
	redisDefaultCount             = 1000
	redisDefaultResultNum         = 5
)

type redisOperation struct {
	match   string
	service *config.Service
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
	scanCmd.Flags().StringVar(&conf.match, "match", "", "redis scan  match pattern  default is null")

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
	delCmd.Flags().StringVar(&conf.match, "match", "", "redis scan  match pattern  default is null")

	delCmds = append(delCmds, delCmd)
	for _, fCmd := range delCmds {
		cmd.AddCommand(fCmd)
	}

	return cmd
}

func runRedisScan(conf *redisOperation) error {

	var total int
	cursor := redisDefaultCursor
	printLen := redisDefaultResultNum

	//don't need to open too many conns
	config.Conf.RedisConf.MaxOpenConns = redisDefaultConnNum

	redisCli, err := redis.NewFromConfig(config.Conf.RedisConf)
	if err != nil {
		fmt.Printf("read redis config fail  err :%v.\n", err)
		return err
	}
	ctx := context.Background()
	fmt.Printf("show some results as an example :\n")
	for {
		res := redisCli.Scan(ctx, cursor, conf.match, redisDefaultCount)
		keys, cur := res.Val()

		if len(keys) >= printLen {
			for _, v := range keys[:printLen] {
				fmt.Printf("%v \n", v)
			}
			total = redisDefaultResultNum
			break

		} else if len(keys) > 0 && len(keys) < redisDefaultResultNum {
			for _, v := range keys {
				fmt.Printf("%v \n", v)
			}
			total += len(keys)
			printLen = printLen - len(keys)
		}

		if cur == 0 {
			break
		}
		cursor = cur
	}
	fmt.Printf("example end \n")

	return nil
}

func runRedisScanDel(conf *redisOperation) error {

	var (
		start   int
		keysTmp []string
		bFlag   bool
		total   uint64
	)
	cursor := redisDefaultCursor

	//don't need to open too many conns
	config.Conf.RedisConf.MaxOpenConns = redisDefaultConnNum

	redisCli, err := redis.NewFromConfig(config.Conf.RedisConf)
	if err != nil {
		fmt.Printf("read redis config fail, err :%v.\n", err)
		return err
	}

	ctx := context.Background()
	fmt.Printf("del keys start :\n")
	for {
		res := redisCli.Scan(ctx, cursor, conf.match, redisDefaultCount)
		keys, cur := res.Val()
		keysNum := len(keys)

		for keysNum > 0 {
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

			time.Sleep(10 * time.Millisecond)
			start = start + redisDelBatchDefaulNum
		}
		start = 0
		if bFlag {
			fmt.Printf("del the keys:(%v) fail,err: %v\n", keysTmp, err)
			bFlag = false
			continue
		} else {
			if len(keys) > 0 {
				fmt.Printf("del keys num: %d. \n", len(keys))
			}
		}

		total += uint64(len(keys))
		if cur == redisDefaultCursor {
			break
		}
		cursor = cur
	}

	fmt.Printf("del keys success ,total num is %v\n", total)

	return nil
}
