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
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"configcenter/src/common/metadata"
	"configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/common/watch"
	"configcenter/src/tools/cmdb_ctl/app/config"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

func init() {
	rootCmd.AddCommand(NewWatchCommand())
}

type watchConf struct {
	startFrom int64
	cursor    string
	resource  string
	fields    []string
	filter    string
}

func (w *watchConf) addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&w.cursor, "cursor", "", "the start cursor from where to watch")
	cmd.PersistentFlags().Int64Var(&w.startFrom, "start-from", 0, "unix time, where to start from, can be negative, which is means start from now-(start-from)")
	cmd.PersistentFlags().StringVar(&w.resource, "rsc", "host", "the resource to watch, can be host or host_relation")
	cmd.PersistentFlags().StringSliceVar(&w.fields, "fields", nil, "the resource fields to return")
	cmd.PersistentFlags().StringVar(&w.filter, "filter", "", "a k:v pair to filter events, k and v is separate with ':' , multiple kv is separated with ';', like k1:v1;k2:v2")
}

func NewWatchCommand() *cobra.Command {
	conf := new(watchConf)

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "watch resource related operation",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "decode",
		Short: "decode a cursor information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDecodeCursor(conf)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "start",
		Short: "decode a cursor information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStartFromWatch(conf)
		},
	})

	conf.addFlags(cmd)
	return cmd
}

func runDecodeCursor(c *watchConf) error {
	cursor := new(watch.Cursor)
	if err := cursor.Decode(c.cursor); err != nil {
		return err
	}
	fmt.Printf("\ndecode cursor: %s\n", c.cursor)
	fmt.Printf("----------------\n")
	fmt.Printf("         type: %s\n", cursor.Type)
	fmt.Printf("          oid: %s\n", cursor.Oid)
	fmt.Printf("     unixTime: %d:%d\n", cursor.ClusterTime.Sec, cursor.ClusterTime.Nano)
	fmt.Printf("  clusterTime: %s\n\n", time.Unix(int64(cursor.ClusterTime.Sec), int64(cursor.ClusterTime.Nano)).Format(time.RFC3339))
	return nil
}

func runStartFromWatch(c *watchConf) error {
	zk, err := config.NewZkService(config.Conf.ZkAddr)
	if err != nil {
		fmt.Printf("new zk client failed, err: %v\n", err)
		return err
	}

	path := types.CC_SERV_BASEPATH + "/" + types.CC_MODULE_EVENTSERVER
	children, err := zk.ZkCli.GetChildren(path)
	if err != nil {
		fmt.Printf("get event server failed, err: %v\n", err)
		return err
	}

	server := ""
	for _, child := range children {
		node, err := zk.ZkCli.Get(path + "/" + child)
		if err != nil {
			return err
		}
		svr := new(types.EventServInfo)
		if err := json.Unmarshal([]byte(node), svr); err != nil {
			return err
		}
		server = fmt.Sprintf("%s:%d", svr.RegisterIP, svr.Port)
		break
	}

	if server == "" {
		return fmt.Errorf("no event server")
	}
	fmt.Println("server: ", server)

	if c.startFrom < 0 {
		c.startFrom = time.Now().Unix() - c.startFrom
	}

	filter := make(map[string]string)
	for _, f := range strings.Split(c.filter, ";") {
		array := strings.Split(f, ":")
		if len(array) != 2 {
			continue
		}
		filter[array[0]] = array[1]
	}
	fmt.Println(">> watch with filter: ", filter)

	opt := watch.WatchEventOptions{
		Fields:    c.fields,
		StartFrom: c.startFrom,
		Cursor:    c.cursor,
		Resource:  watch.CursorType(c.resource),
	}

	optByte, _ := json.Marshal(opt)
	rid := util.GenerateRID()
	fmt.Printf(">> rid: %s\n>> request options: %s\n", rid, string(optByte))

	client := new(http.Client)
	url := fmt.Sprintf("http://%s/event/v3/watch/resource/%s", server, c.resource)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(optByte))
	if err != nil {
		return err
	}
	req.Header.Add("HTTP_BLUEKING_SUPPLIER_ID", "0")
	req.Header.Add("BK_User", "cmdb_tool")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cc_Request_Id", rid)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	type response struct {
		metadata.BaseResp `json:",inline"`
		Data              WatchResp `json:"data"`
	}

	event := new(response)
	if err := json.NewDecoder(resp.Body).Decode(event); err != nil {
		return err
	}

	for {

		if !event.Result {
			return fmt.Errorf("request failed, err: %s", event.ErrMsg)
		}

		if !event.Data.Watched {
			fmt.Printf(">>> watched 0 event, try next round...\n")
			opt = watch.WatchEventOptions{
				Fields:   c.fields,
				Cursor:   event.Data.Events[0].Cursor,
				Resource: watch.CursorType(c.resource),
			}
		} else {
			js, _ := json.MarshalIndent(event.Data.Events, "", "    ")

			if len(filter) != 0 {
				allEvent := gjson.ParseBytes(js)
				all := make([]string, 0)

				for _, event := range allEvent.Array() {
					hit := true
					for k, v := range filter {
						if event.Get("bk_detail").Get(k).String() != v {
							hit = false
							break
						}
					}
					if hit {
						all = append(all, event.Raw)
					}
				}
				if len(all) == 0 {
					fmt.Printf("\n>>>watched filtered 0 events, try next round... \n ")
				} else {
					fmt.Printf("\n>>>watched filtered %d events -> : \n %v \n", len(all), all)
				}
			} else {
				fmt.Printf("\n>>>watched %d events -> : \n %s \n", len(event.Data.Events), string(js))
			}

			if len(event.Data.Events) <= 0 {
				fmt.Println("** received unknown error **, exit")
				return errors.New("unknown error")
			}
			opt = watch.WatchEventOptions{
				Fields:   c.fields,
				Cursor:   event.Data.Events[len(event.Data.Events)-1].Cursor,
				Resource: watch.CursorType(c.resource),
			}

		}

		optByte, _ := json.Marshal(opt)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(optByte))
		if err != nil {
			return err
		}
		req.Header.Add("HTTP_BLUEKING_SUPPLIER_ID", "0")
		req.Header.Add("BK_User", "cmdb_tool")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Cc_Request_Id", "666666666666666666")
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		event = new(response)
		if err := json.NewDecoder(resp.Body).Decode(event); err != nil {
			return err
		}
	}

	return nil
}

type WatchResp struct {
	// watched events or not
	Watched bool      `json:"bk_watched"`
	Events  []*Detail `json:"bk_events"`
}

type Detail struct {
	Cursor    string `json:"bk_cursor"`
	Resource  string `json:"bk_resource"`
	EventType string `json:"bk_event_type"`
	// Default instance is JsonString type
	Detail map[string]interface{} `json:"bk_detail"`
}
