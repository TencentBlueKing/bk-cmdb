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
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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
	cmd.Flags().StringP("dir","d", "", "the directory path where the configuration file is located")
	cmd.Flags().StringP("file","f", "", "the path of the configuration file")
	return cmd
}
// isExists checks target dir/file exist or not.
func isExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func checkConf(cmd *cobra.Command) error {
	// check the configuration of a single file
	file, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}
	if len(file) != 0 {
		if !isExists(file) {
			fmt.Println("this file does not exist")
			return nil
		}
		if err := checkFileType(file); err != nil {
			fmt.Println("file format error,the file is " + file)
			return nil
		}
	}

	// check the configuration files in the directory
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		return err
	}
	if len(dir) != 0 {
		if !isExists(dir) {
			fmt.Println("this directory does not exist")
			return nil
		}
		// query config file from configPath
		filePaths, err := filepath.Glob(filepath.Join(dir, "*"))
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