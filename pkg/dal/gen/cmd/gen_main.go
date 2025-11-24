/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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

// Package main provides gorm/gen code generation
package main

import (
	"flag"
	"log"
	"os"

	"github.com/samber/lo"
	"gorm.io/gen"

	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/table"
)

//go:generate go run gen_main.go -cwd ..

var outPath = "."
var cwd = "."

func main() {
	flag.StringVar(&outPath, "out-path", outPath, "output path for generated code")
	flag.StringVar(&cwd, "cwd", cwd, "current working directory")
	flag.Parse()
	err := os.Chdir(cwd)
	if err != nil {
		log.Fatalf("fail to change working directory to %s, err: %v", cwd, err)
		return
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           outPath,
		FieldWithIndexTag: true,
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	// 加载全部静态表
	models := lo.Values(table.GetAllStaticTables())
	// id生成器不在静态模型中，单独处理
	models = append(models, table.IDGenerator{})
	// 这里添加需要生成的模型
	g.ApplyBasic(models...)
	g.Execute()
}
