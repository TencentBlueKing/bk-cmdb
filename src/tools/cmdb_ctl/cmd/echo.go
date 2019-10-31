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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewEchoCommand())
}

type echoConf struct {
	url string
}

func NewEchoCommand() *cobra.Command {
	conf := new(echoConf)

	cmd := &cobra.Command{
		Use:   "echo",
		Short: "echo server for http callback",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEchoServer(conf)
		},
	}

	conf.addFlags(cmd)

	return cmd
}

func (c *echoConf) addFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.url, "url", "", "the url of the echo server")
}

func runEchoServer(c *echoConf) error {
	u, err := url.Parse(c.url)
	if err != nil {
		return err
	}
	http.HandleFunc(u.Path, echo)
	if err := http.ListenAndServe(u.Host, nil); err != nil {
		return err
	}
	return nil
}

func echo(w http.ResponseWriter, r *http.Request) {
	s, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read http request failed. error: %s", err.Error())
	}
	fmt.Fprintf(os.Stdout, "%s\n", s)
}
