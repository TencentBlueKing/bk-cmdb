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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewEchoCommand())
}

type echo struct {
	url        string
	jsonPretty bool
}

func NewEchoCommand() *cobra.Command {
	conf := new(echo)

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

func (c *echo) addFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.url, "url", "", "the url of the echo server, eg: http://127.0.0.1:80/echo")
	cmd.Flags().BoolVar(&c.jsonPretty, "pretty", false, "json indent the received data if it's json format.")
}

func runEchoServer(c *echo) error {
	u, err := url.Parse(c.url)
	if err != nil {
		return err
	}
	http.HandleFunc(u.Path, c.echoServer)
	if err := http.ListenAndServe(u.Host, nil); err != nil {
		return err
	}
	return nil
}

func (c *echo) echoServer(w http.ResponseWriter, r *http.Request) {
	s, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read http request failed. error: %s", err.Error())
		return
	}
	fmt.Fprintf(os.Stdout, "%c[1;40;31m>> received new data, time: %s %c[0m\n", 0x1B, time.Now().Format(time.RFC3339), 0x1B)
	if c.jsonPretty {
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, s, "", "    "); err != nil {
			fmt.Fprintf(os.Stderr, "json indent the body failed. error: %v", err)
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n\n", prettyJSON.Bytes())
		return
	}

	fmt.Fprintf(os.Stdout, "%s\n\n", s)
}
