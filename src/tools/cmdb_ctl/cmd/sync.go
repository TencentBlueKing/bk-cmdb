/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package cmd

import (
	"encoding/json"
	"fmt"

	ftypes "configcenter/pkg/types/sync/full-text-search"
	"configcenter/src/common/metadata"
	"configcenter/src/common/types"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewSyncCommand())
}

// NewSyncCommand new sync command
func NewSyncCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "sync resource related operation",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmd.AddCommand(NewFullTextSearchCmd())

	return cmd
}

// NewFullTextSearchCmd new full-text-search sync command
func NewFullTextSearchCmd() *cobra.Command {
	conf := new(fullTextSearchConf)

	cmd := &cobra.Command{
		Use:   "full-text-search",
		Short: "sync full-text-search related info",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFullTextSearchSync(conf)
		},
	}

	conf.addFlags(cmd)
	return cmd
}

type fullTextSearchConf struct {
	isSyncData bool
	isMigrate  bool
	dataOpt    ftypes.SyncDataOption
}

func (c *fullTextSearchConf) addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&c.isMigrate, "is-migrate", false, "is migrate full-text-search data")
	cmd.PersistentFlags().BoolVar(&c.isSyncData, "is-sync-data", false, "is sync specified full-text-search data")
	cmd.PersistentFlags().BoolVar(&c.dataOpt.IsAll, "is-all", false, "is sync all full-text-search data")
	cmd.PersistentFlags().StringVar(&c.dataOpt.Index, "index", "", "need sync index")
	cmd.PersistentFlags().StringVar(&c.dataOpt.Collection, "collection", "", "need sync collection")
	cmd.PersistentFlags().StringSliceVar(&c.dataOpt.Oids, "oids", make([]string, 0), "need sync data ids")
}

func runFullTextSearchSync(c *fullTextSearchConf) error {
	if c.isMigrate {
		return runFullTextSearchMigrate()
	}

	if c.isSyncData {
		return runFullTextSearchDataSync(&c.dataOpt)
	}

	return fmt.Errorf("one of is-migrate and is-sync-data option must be set")
}

func runFullTextSearchMigrate() error {
	resp, err := doCmdbHttpRequest(types.CC_MODULE_SYNC, "/sync/v3/migrate/full/text/search", "{}")
	if err != nil {
		return err
	}

	res := new(migrateResp)
	if err = json.NewDecoder(resp.Body).Decode(res); err != nil {
		fmt.Printf("decode response body failed, err: %v\n", err)
		return err
	}

	if err = res.CCError(); err != nil {
		fmt.Printf("do full text search migration failed, err: %v\n", err)
		return err
	}

	resJs, err := json.Marshal(res.Data)
	if err != nil {
		fmt.Printf("marshal full text search migration result(%+v) failed, err: %v\n", res.Data, err)
		return err
	}

	fmt.Printf("do full text search migration success, result: %sv\n", string(resJs))
	return nil
}

type migrateResp struct {
	metadata.BaseResp `json:",inline"`
	Data              ftypes.MigrateResult `json:"data"`
}

func runFullTextSearchDataSync(opt *ftypes.SyncDataOption) error {
	resp, err := doCmdbHttpRequest(types.CC_MODULE_SYNC, "/sync/v3/sync/full/text/search/data", opt)
	if err != nil {
		return err
	}

	res := new(metadata.BaseResp)
	if err = json.NewDecoder(resp.Body).Decode(res); err != nil {
		fmt.Printf("decode response body failed, err: %v\n", err)
		return err
	}

	if err = res.CCError(); err != nil {
		fmt.Printf("do full text search migration failed, err: %v\n", err)
		return err
	}

	return nil
}
