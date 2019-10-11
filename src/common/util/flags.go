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

package util

import (
	goflag "flag"
	"os"
	"strings"

	"configcenter/src/common/version"
	"configcenter/src/common/zkclient"

	"github.com/spf13/pflag"
)

// WordSepNormalizeFunc changes all flags that contain "_" separators
func WordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	}
	return pflag.NormalizedName(name)
}

type CommonFlagConfig struct {
	version bool
}

// AddCommonFlags add common flags that is needed by all modules
func AddCommonFlags(cmdline *pflag.FlagSet, zkConf *zkclient.ZkConf, flagConf *CommonFlagConfig) {
	cmdline.BoolVar(&flagConf.version, "version", false, "show version information")
	cmdline.StringVar(&zkConf.ZkAddr, "zkaddr", "", "The zookeeper server address, e.g: 127.0.0.1:2181")
	cmdline.StringVar(&zkConf.ZkUser, "zkuser", "", "The zookeeper auth user")
	cmdline.StringVar(&zkConf.ZkPwd, "zkpwd", "", "The zookeeper auth password")
}

// InitFlags normalizes and parses the command line flags
func InitFlags(zkConf *zkclient.ZkConf) {
	pflag.CommandLine.SetNormalizeFunc(WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flagConf := new(CommonFlagConfig)
	AddCommonFlags(pflag.CommandLine, zkConf, flagConf)
	pflag.Parse()

	// add handler if flag include --version/-v
	if flagConf.version {
		version.ShowVersion()
		os.Exit(0)
	}
}
