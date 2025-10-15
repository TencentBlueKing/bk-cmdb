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

package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"testing"
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

	logger := slog.New(NewContextualHandler(WithJsonFormat()))
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

	h := NewContextualHandler(
		WithLevel(LevelTrace),
		WithJsonFormat(),
		// WithFileOption(FileOption{Filename: "./example.log"}),
	)
	SetDefault(h)

	Trace(ctx, "msg0")
	Depth(0).Trace(ctx, "msg")

	// gorm example
	l2 := slog.New(With("ddd", "ko").Handler())
	l2.InfoContext(ctx, "abc")

	ctx = WithAttr(ctx, slog.GroupAttrs("g1", slog.String("w1", "k1")))
	Depth(0).Trace(ctx, "msg1")
	Depth(0).With("k", "o").WithGroup("m").Trace(ctx, "msg1", "w4", "k4")
	Depth(1).With("k", "o").WithGroup("m").With("w2", "k2", "w3", "k3").With(slog.Any("", nil)).With().Trace(ctx, "msg1", slog.Any("", "v"))
}

func BenchmarkLogger(b *testing.B) {
	textHandler := NewContextualHandler(WithLevel(slog.LevelDebug), WithWriter(io.Discard))
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

	_logger := slog.New(textHandler)
	_logger = _logger.With(RidAttr("8e388ed2-ba59-48b7-b213-cf1afd6ac1e9"))

	for b.Loop() {
		_logger.Debug("msg")
	}
}

func BenchmarkLoggerParallel(b *testing.B) {
	textHandler := NewContextualHandler(WithLevel(LevelTrace), WithWriter(io.Discard))
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

	_logger := slog.New(textHandler)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_logger.Debug("msg")
		}
	})
}
