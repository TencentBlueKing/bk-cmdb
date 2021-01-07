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

// apiserver对于接收到的请求可以配置限流策略
// 该命令行工具可对这些限流策略（规则）进行增删查的操作
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"configcenter/src/common/metadata"
	"configcenter/src/common/types"
	"configcenter/src/tools/cmdb_ctl/app/config"

	"github.com/spf13/cobra"
)

const (
	// limiter introduction
	limiterIntro = `
********************************************************
示例:
以下命令是在配置了ZK_ADDR环境变量的情况下使用，没有配置时也可以通过命令行参数--zk-addr指定
# 列出所有策略
./tool_ctl limiter ls
# 配置策略，对url限制请求次数
./tool_ctl limiter set --rule='{"rulename":"rule1","appcode":"gse","user":"admin","ip":"","method":"POST","url":"^/api/v3/module/search/[^\\s/]+/[0-9]+/[0-9]+/?$","limit":1000,"ttl":60,"denyall":false}'
# 配置策略，将url直接禁掉
./tool_ctl limiter set --rule='{"rulename":"rule1","appcode":"gse","user":"admin","url":"^/api/v3/module/search/[^\\s/]+/[0-9]+/[0-9]+/?$","denyall":true}'
# 获取某些策略详情
./tool_ctl limiter get --rulenames=test1,test2
# 删除某些策略
./tool_ctl limiter del --rulenames=test1,test2
********************************************************
		`
	// ruleIntro rule introduction
	ruleIntro = `
********************************************************
- rule策略字段说明

| 字段     | 类型   | 必选 | 描述                                                         |
|----------|--------|------|--------------------------------------------------------------|
| rulename | string | 是   | 策略名                                                       |
| appcode  | string | 否   | 应用ID                                                       |
| user     | string | 否   | 请求发起的用户名                                             |
| ip       | string | 否   | api的来源ip                                                  |
| method   | string | 否   | 请求的类型，配置的情况下只能为POST、GET、PUT、DELETE中的一种 |
| url      | string | 否   | api的url正则表达式                                           |
| limit    | int64  | 否   | api请求限制总次数                                            |
| ttl      | int64  | 否   | 策略存活时间，单位为秒                                       |
| denyall  | bool   | 否   | 是否直接禁掉请求，默认为false，为true时忽略limit和ttl参数    |
 
appcode、user、ip、method、url需要至少配置一项  
denyall配置为false的情况下，limit和ttl配置才能生效
********************************************************
		`
)

func init() {
	rootCmd.AddCommand(NewLimiterCommand())
}

type limiterConf struct {
	rule      string
	rulenames string
}

func (c *limiterConf) addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&c.rule, "rule", "",
		`the api limiter rule to set, a json like '{"rulename":"rule1","appcode":"gse","user":"","ip":"","method":"POST","url":"^/api/v3/module/search/[^\\s/]+/[0-9]+/[0-9]+/?$","limit":1000,"ttl":60,"denyall":false}'`)
	cmd.PersistentFlags().StringVar(&c.rulenames, "rulenames", "", `the api limiter rule names to get or del, multiple names is separated with ',',like 'name1,name2'`)
}

func NewLimiterCommand() *cobra.Command {
	conf := new(limiterConf)

	cmd := &cobra.Command{
		Use:   "limiter",
		Short: "api limiter operations",
		Long:  limiterIntro,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "set",
		Short: "set api limiter rule, use with flag --rule",
		Long:  ruleIntro,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetRule(conf)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "get api limiter rules according rule names,use with flag --rulenames",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetRules(conf)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "del",
		Short: "del api limiter rules, use with flag --rulenames",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelRules(conf)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "ls",
		Short: "list all api limiter rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListRules(conf)
		},
	})

	conf.addFlags(cmd)

	return cmd
}

func runSetRule(c *limiterConf) error {
	rule := new(metadata.LimiterRule)
	err := json.Unmarshal([]byte(c.rule), rule)
	if err != nil {
		return err
	}

	err = rule.Verify()
	if err != nil {
		return err
	}

	zk, err := config.NewZkService(config.Conf.ZkAddr)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s/%s", types.CC_SERVLIMITER_BASEPATH, rule.RuleName)
	exist, err := zk.ZkCli.Exist(path)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("the rule %s has already existed", rule.RuleName)
	}

	data, err := json.Marshal(rule)
	if err != nil {
		return err
	}

	err = zk.ZkCli.CreateDeepNode(path, data)
	if err != nil {
		return err
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, data, "", "\t")
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(os.Stdout, "%s\nset rule successfully!\n", pretty.String())

	return nil
}

func runGetRules(c *limiterConf) error {
	if c.rulenames == "" {
		return fmt.Errorf("rulenames must be set")
	}
	zk, err := config.NewZkService(config.Conf.ZkAddr)
	if err != nil {
		return err
	}
	names := strings.Split(c.rulenames, ",")
	for _, name := range names {
		path := fmt.Sprintf("%s/%s", types.CC_SERVLIMITER_BASEPATH, name)
		data, err := zk.ZkCli.Get(path)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "get rule %s err:%s\n", name, err)
			continue
		}
		var pretty bytes.Buffer
		err = json.Indent(&pretty, []byte(data), "", "\t")
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "get rule %s Indent err:%s\n", name, err)
			continue
		}
		_, _ = fmt.Fprintf(os.Stdout, "%s\n%s\n\n", path, pretty.String())
	}
	return nil
}

func runDelRules(c *limiterConf) error {
	if c.rulenames == "" {
		return fmt.Errorf("rulenames must be set")
	}
	zk, err := config.NewZkService(config.Conf.ZkAddr)
	if err != nil {
		return err
	}
	names := strings.Split(c.rulenames, ",")
	for _, name := range names {
		path := fmt.Sprintf("%s/%s", types.CC_SERVLIMITER_BASEPATH, name)
		err := zk.ZkCli.Del(path, -1)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "del rule %s err:%s\n", name, err)
			continue
		}
		_, _ = fmt.Fprintf(os.Stdout, "del rule %v successfully!\n", name)
	}
	return nil
}

func runListRules(c *limiterConf) error {
	zk, err := config.NewZkService(config.Conf.ZkAddr)
	if err != nil {
		return err
	}
	path := types.CC_SERVLIMITER_BASEPATH
	children, err := zk.ZkCli.GetChildren(path)
	if err != nil {
		return err
	}
	for _, child := range children {
		data, err := zk.ZkCli.Get(path + "/" + child)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "list rule %s Get err:%s\n", child, err)
			continue
		}
		var pretty bytes.Buffer
		err = json.Indent(&pretty, []byte(data), "", "\t")
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "list rule %s Indent err:%s\n", child, err)
			continue
		}
		_, _ = fmt.Fprintf(os.Stdout, "%s\n%s\n\n", path+"/"+child, pretty.String())
	}
	return nil
}
