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

// Package logger provider the contextual and structured logger, base the golang's slog
package logger

import (
	"context"
	"log/slog"
)

// Logger contextual and structured logger
type Logger interface {
	With(args ...any) Logger
	WithGroup(name string) Logger
	Handler() slog.Handler

	Trace(ctx context.Context, msg string, args ...any)
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
}

// SetDefault make a new Logger and set to default logger
func SetDefault(h slog.Handler) {
	defaultLogger = newLogger(h, 0)
}

// Default get default logger's
func Default() Logger {
	return defaultLogger
}

// New make a new Logger
func New(h slog.Handler) Logger {
	return newLogger(h, 0)
}

// With returns a Logger that includes the given attributes in each output operation.
func With(args ...any) Logger {
	return defaultLogger.With(args...)
}

// WithGroup returns a Logger that starts a group
func WithGroup(name string) Logger {
	return defaultLogger.WithGroup(name)
}

// Depth return a Logger with depth
func Depth(depth int) Logger {
	d := -depth
	return newLogger(defaultLogger.Handler(), d)
}

// Trace logs at LevelTrace with the given context.
func Trace(ctx context.Context, msg string, args ...any) {
	defaultLogger.Log(ctx, LevelTrace, msg, args...)
}

// Debug logs at LevelDebug with the given context.
func Debug(ctx context.Context, msg string, args ...any) {
	defaultLogger.Log(ctx, slog.LevelDebug, msg, args...)
}

// Info logs at LevelInfo with the given context.
func Info(ctx context.Context, msg string, args ...any) {
	defaultLogger.Log(ctx, slog.LevelInfo, msg, args...)
}

// Warn logs at LevelWarn with the given context.
func Warn(ctx context.Context, msg string, args ...any) {
	defaultLogger.Log(ctx, slog.LevelWarn, msg, args...)
}

// Error logs at LevelError with the given context.
func Error(ctx context.Context, msg string, args ...any) {
	defaultLogger.Log(ctx, slog.LevelError, msg, args...)
}

type contextualLogger struct {
	logger *slog.Logger
	depth  int
}

func newLogger(h slog.Handler, depth int) *contextualLogger {
	return &contextualLogger{logger: slog.New(h), depth: depth}
}

// With returns a Logger that includes the given attributes in each output operation.
func (l *contextualLogger) With(args ...any) Logger {
	return newLogger(l.logger.With(args...).Handler(), l.depth)
}

// WithGroup returns a Logger that starts a group
func (l *contextualLogger) WithGroup(name string) Logger {
	return newLogger(l.logger.WithGroup(name).Handler(), l.depth)
}

// Handler returns logger's Handler.
func (l *contextualLogger) Handler() slog.Handler {
	return l.logger.Handler()
}

// Trace logs at LevelInfo with the given context.
func (l *contextualLogger) Trace(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, LevelTrace, msg, args...)
}

// Debug logs at LevelInfo with the given context.
func (l *contextualLogger) Debug(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, slog.LevelDebug, msg, args...)
}

// Info logs at LevelInfo with the given context.
func (l *contextualLogger) Info(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, slog.LevelInfo, msg, args...)
}

// Warn logs at LevelInfo with the given context.
func (l *contextualLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, slog.LevelWarn, msg, args...)
}

// Error logs at LevelInfo with the given context.
func (l *contextualLogger) Error(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, slog.LevelError, msg, args...)
}

// Log logs at Level with the given context.
func (l *contextualLogger) Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	ctx = context.WithValue(ctx, depthCtxKey, l.depth)
	l.logger.Log(ctx, level, msg, args...)
}

// WithAttr 添加自定义属性
func WithAttr(ctx context.Context, attrs ...slog.Attr) context.Context {
	rawAttrs, ok := ctx.Value(attrCtxKey).([]slog.Attr)
	if !ok {
		return context.WithValue(ctx, attrCtxKey, attrs)
	}

	rawAttrs = append(rawAttrs, attrs...)
	ctx = context.WithValue(ctx, attrCtxKey, rawAttrs)
	return ctx
}

// RidAttr request_id类型Attr
func RidAttr(rid string) slog.Attr {
	return slog.String("rid", rid)
}

// E error类型Attr
func E(err error) slog.Attr {
	return slog.String("err", err.Error())
}

func init() {
	handler := NewContextualHandler(NewHandlerOptions())
	SetDefault(handler)
}
