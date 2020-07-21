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
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)
const(
	DefaultDir = "../cmdb_adminserver/configures"
)

func init() {
	rootCmd.AddCommand(NewConfCommand())
}


func NewConfCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "checkconf",
		Short: "check if the yaml file format is correct",
		RunE: func(cmd *cobra.Command, args []string) error {
			return checkConf(cmd)
		},
	}
	cmd.Flags().StringP("path","p", "", "the path of config")
	return cmd
}

func checkConf(cmd *cobra.Command) error {
	path, err := cmd.Flags().GetString("path")
	if err != nil {
		return err
	}
	var configPath string
	if len(path) == 0 {
		configPath = DefaultDir
	} else {
		configPath = path
	}
	fmt.Println("the path is "+configPath)
	// query config file from configPath
	filePaths, err := filepath.Glob(filepath.Join(configPath, "*"))
	if err != nil {
		return err
	}
	// check config file
	for _, filepath := range filePaths {
		err := checkFileType(filepath)
		if err != nil {
			fmt.Println("file format error,the file is " + filepath)
			return nil
		}
	}
	fmt.Println("vaild successfuly!")
	return nil
}

func checkFileType(path string) error {
	s := strings.Split(path,".")
	fileType := s[len(s) - 1]
	if fileType != "yaml" && fileType != "yml" {
		return fmt.Errorf("file format error")
	}
	file, _ := ioutil.ReadFile(path)
	result := make(map[string]interface{})
	err := yaml.Unmarshal(file, &result)
	if err != nil {
		return err
	}
	return nil
}