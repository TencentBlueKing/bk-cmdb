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

package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"configcenter/src/storage/dal/redis"

	rawRedis "github.com/go-redis/redis/v7"
	"github.com/spf13/cobra"
)

var (
	hostname   string
	topic      string
	message    string
	file       string
	requestNum int
	clientNum  int
)

func genBenchmarkMQCmd() *cobra.Command {
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.IntVar(&requestNum, "num", 1, "Num of requests.")
	f.IntVar(&clientNum, "conn", 1, "Num of clients.")
	f.StringVar(&hostname, "host", "localhost:6379", "Host of redis server.")
	f.StringVar(&topic, "topic", "benchmark", "Topic of redis SubPub MQ.")
	f.StringVar(&message, "data", "{}", "Message for benchmark.")
	f.StringVar(&file, "file", "", "File path of meessage for benchmark, eg ./data.json.")

	mqCmd := &cobra.Command{
		Use:   "mq",
		Short: "mq benchmark command",
		Run: func(cmd *cobra.Command, args []string) {
			wg := &sync.WaitGroup{}
			wg.Add(clientNum)

			if len(file) != 0 {
				// read benchmark message data from local file.
				data, err := ioutil.ReadFile(file)
				if err != nil {
					panic(err)
				}
				message = string(data)
			}

			for i := 0; i < clientNum; i++ {
				go func() {
					client := redis.NewClient(&rawRedis.Options{Addr: hostname})
					if _, err := client.Ping(context.Background()).Result(); err != nil {
						panic(err)
					}
					for i := 0; i < requestNum; i++ {
						if _, err := client.Publish(context.Background(), topic, message).Result(); err != nil {
							panic(err)
						}
					}

					wg.Done()
				}()
			}

			// wait for cmd.
			wg.Wait()

			// done.
			fmt.Println("Benchmark Done!")
		},
	}

	// command flags.
	mqCmd.Flags().AddGoFlag(f.Lookup("host"))
	mqCmd.Flags().AddGoFlag(f.Lookup("num"))
	mqCmd.Flags().AddGoFlag(f.Lookup("conn"))
	mqCmd.Flags().AddGoFlag(f.Lookup("topic"))
	mqCmd.Flags().AddGoFlag(f.Lookup("data"))
	mqCmd.Flags().AddGoFlag(f.Lookup("file"))

	return mqCmd
}

func main() {
	// root command.
	rootCmd := &cobra.Command{Use: "benchmark"}

	// benchmark command in MQ mode.
	benchmark := genBenchmarkMQCmd()

	// add sub cmds to root cmd.
	rootCmd.AddCommand(benchmark)

	// run root command.
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
