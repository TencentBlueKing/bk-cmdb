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

package log

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLogger(t *testing.T) {
	ctx := context.Background()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	logger = logger.With(RidAttr("8e388ed2-ba59-48b7-b213-cf1afd6ac1e9"))
	logger.InfoContext(ctx, "msg0")
	logger.InfoContext(ctx, "msg")

	logger = logger.With(slog.GroupAttrs("g1", slog.String("w1", "k1")))
	logger.With("k", "o").WithGroup("m").InfoContext(ctx, "msg1", "w4", "k4")
	logger.With("k", "o").WithGroup("m").With("w2", "k2", "w3", "k3").WithGroup("g2").InfoContext(ctx, "msg1")
}

func TestLoggerHandler(t *testing.T) {
	ctx := context.Background()

	opts := NewHandlerOptions()
	opts.Format = "json"
	logger := slog.New(NewContextualHandler(opts))

	logger = logger.With(RidAttr("8e388ed2-ba59-48b7-b213-cf1afd6ac1e9"))
	logger.InfoContext(ctx, "msg0")
	logger.InfoContext(ctx, "msg")

	logger = logger.With(slog.GroupAttrs("g1", slog.String("w1", "k1")))
	logger.With("k", "o").WithGroup("m").InfoContext(ctx, "msg1", "w4", "k4")
	logger.With("k", "o").WithGroup("m").With("w2", "k2", "w3", "k3").WithGroup("g2").InfoContext(ctx, "msg1")
}

func TestLogger(t *testing.T) {
	ctx := context.Background()
	ctx = WithAttr(ctx, RidAttr("8e388ed2-ba59-48b7-b213-cf1afd6ac1e9"))

	opts := &HandlerOptions{Level: getLevelName(LevelTrace), Format: "json", Stdout: true}
	opts.LogDir = t.TempDir()
	h := NewContextualHandler(opts)

	SetDefault(h)

	Trace(ctx, "msg0")
	Depth(0).Trace(ctx, "msg")

	// gorm example
	l2 := slog.New(With("ddd", "ko").Handler())
	l2.InfoContext(ctx, "abc")

	ctx = WithAttr(ctx, slog.GroupAttrs("g1", slog.String("w1", "k1")))
	Depth(0).Trace(ctx, "msg1")
	Default().Trace(ctx, "msgddd")
	Depth(0).With("k", "o").WithGroup("m").Trace(ctx, "msg1", "w4", "k4")
	Depth(1).With("k", "o").WithGroup("m").With("w2", "k2", "w3", "k3").With(slog.Any("", nil)).With().Trace(ctx, "msg1", slog.Any("", "v"))
}

func TestLoggerAttr(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	baseH := slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue("test")
			}
			return a
		},
	})

	h := &contextualHandler{Handler: baseH}
	logger := slog.New(h)
	SetDefault(h)

	ctx := context.Background()
	ctx = WithAttr(ctx, RidAttr("8e388ed2-ba59-48b7-b213-cf1afd6ac1e9"))

	logger.InfoContext(ctx, "test attr")
	output := buf.String()
	buf.Reset()

	Info(ctx, "test attr")
	assert.Equal(t, buf.String(), output)
	buf.Reset()

	logger = logger.With(slog.GroupAttrs("g1", slog.String("w1", "k1")))
	logger.With("k", "o").WithGroup("m").InfoContext(ctx, "msg1", "w4", "k4")
	output = buf.String()
	buf.Reset()

	With(slog.GroupAttrs("g1", slog.String("w1", "k1"))).With("k", "o").WithGroup("m").Info(ctx, "msg1", "w4", "k4")
	assert.Equal(t, buf.String(), output)
}

func TestLoggerDepth(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	baseH := slog.NewTextHandler(buf, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelInfo,
		ReplaceAttr: replaceAttr,
	})

	h := &contextualHandler{Handler: baseH}
	logger := slog.New(h)

	ctx := context.Background()
	// 场景1
	logger.InfoContext(ctx, "test depth")
	assert.True(t, strings.Contains(buf.String(), "source=log/logger_test.go:"))
	buf.Reset()

	// 场景2
	slog.SetDefault(logger)
	slog.InfoContext(ctx, "test depth")
	assert.True(t, strings.Contains(buf.String(), "source=log/logger_test.go:"))
	buf.Reset()

	// 场景3
	SetDefault(h)
	Info(ctx, "test depth")
	assert.True(t, strings.Contains(buf.String(), "source=log/logger_test.go:"))
	buf.Reset()

	// 场景5
	Default().Info(ctx, "test depth")
	assert.True(t, strings.Contains(buf.String(), "source=log/logger_test.go:"))
	buf.Reset()

	Default().Info(ctx, "test depth")
	assert.True(t, strings.Contains(buf.String(), "source=log/logger_test.go:"))
	buf.Reset()

	// 场景6
	Depth(0).Info(ctx, "test depth")
	assert.True(t, strings.Contains(buf.String(), "source=log/logger_test.go:"))
	buf.Reset()

	// 场景7
	Depth(1).Info(ctx, "test depth")
	assert.True(t, strings.Contains(buf.String(), "source=log/logger.go:"))
	buf.Reset()
}

func BenchmarkLogger(b *testing.B) {
	textHandler := NewContextualHandler(&HandlerOptions{Level: getLevelName(LevelTrace)})
	SetDefault(textHandler)

	ctx := context.Background()
	ctx = WithAttr(ctx, RidAttr("8e388ed2-ba59-48b7-b213-cf1afd6ac1e9"))

	for b.Loop() {
		Debug(ctx, "msg")
	}
}

func BenchmarkStdLogger(b *testing.B) {
	textHandler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: replaceAttr,
	})

	logger := slog.New(textHandler)
	logger = logger.With(RidAttr("8e388ed2-ba59-48b7-b213-cf1afd6ac1e9"))

	for b.Loop() {
		logger.Debug("msg")
	}
}

func BenchmarkLoggerParallel(b *testing.B) {
	textHandler := NewContextualHandler(&HandlerOptions{Level: getLevelName(LevelTrace)})
	SetDefault(textHandler)

	ctx := context.Background()
	ctx = WithAttr(ctx, slog.String("rid", "8e388ed2-ba59-48b7-b213-cf1afd6ac1e9"))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Debug(ctx, "msg")
		}
	})
}

func BenchmarkStdLoggerParallel(b *testing.B) {
	textHandler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: replaceAttr,
	})

	logger := slog.New(textHandler)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debug("msg")
		}
	})
}
