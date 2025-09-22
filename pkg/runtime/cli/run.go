/*
 * Tencent is pleased to support the open source community by making
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

// Package cli provider app Run entry
package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/TencentBlueKing/bk-cmdb/pkg/version"
)

var (
	// SignalChan ...
	SignalChan = make(chan os.Signal, 1)
	// ErrSignal ...
	ErrSignal = errors.New("signal")
)

// Run provides the common boilerplate code around executing a cobra command.
func Run(cmd *cobra.Command) int {
	// 不开启 自动排序
	cobra.EnableCommandSorting = false
	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	// 添加版本
	cmd.SetVersionTemplate(`{{print .Version}}`)
	cmd.Version = version.GetVersion()

	if err := execute(cmd); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	return 0
}

func execute(cmd *cobra.Command) error {
	ctx := context.Background()
	// graceful shutdown signal
	signal.Notify(SignalChan, syscall.SIGINT, syscall.SIGTERM)

	err := cmd.ExecuteContext(ctx)
	if err == nil {
		return nil
	}

	if errors.Is(err, ErrSignal) {
		return nil
	}
	return err
}
