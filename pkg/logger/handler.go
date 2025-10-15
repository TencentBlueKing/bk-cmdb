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
	"path/filepath"
	"runtime"
	"strconv"

	"gopkg.in/natefinch/lumberjack.v2"
)

type ctxKey string

const (
	// LevelTrace ...
	LevelTrace = slog.Level(-8)
)

var (
	attrCtxKey    = ctxKey("logger.attr")
	depthCtxKey   = ctxKey("logger.depth")
	defaultLogger = new(contextualLogger)
)

// FileOption 输出到文件的配置
type FileOption struct {
	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.
	Filename string `json:"filename" yaml:"filename"`

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `json:"maxsize" yaml:"maxsize"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `json:"maxage" yaml:"maxage"`

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `json:"maxbackups" yaml:"maxbackups"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool `json:"localtime" yaml:"localtime"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `json:"compress" yaml:"compress"`
}

// contextualHandler 支持结构化/上下文的Handler
type contextualHandler struct {
	slog.Handler
	w       io.Writer
	level   slog.Leveler
	attr    *groupOrAttrs // 自定义attr
	fileOpt *FileOption   // 同时输出到文件
	format  string        // 日志格式text/json
}

// Option setter for contextualHandler
type Option func(*contextualHandler)

// WithWriter set handler base writer
func WithWriter(w io.Writer) Option {
	return func(h *contextualHandler) {
		h.w = w
	}
}

// WithLevel set handler level
func WithLevel(level slog.Leveler) Option {
	return func(h *contextualHandler) {
		h.level = level
	}
}

// WithFileOption set handler FileOption
func WithFileOption(opt FileOption) Option {
	return func(h *contextualHandler) {
		h.fileOpt = &opt
	}
}

// WithJsonFormat set output log to json format
func WithJsonFormat() Option {
	return func(h *contextualHandler) {
		h.format = "json"
	}
}

// NewContextualHandler make a new handler
func NewContextualHandler(opts ...Option) slog.Handler {
	h := &contextualHandler{
		format: "text",
		w:      os.Stdout,
	}

	for _, opt := range opts {
		opt(h)
	}

	ioWriter := h.w
	if h.fileOpt != nil && h.fileOpt.Filename != "" {
		fileWriter := &lumberjack.Logger{
			Filename:   h.fileOpt.Filename,
			MaxSize:    h.fileOpt.MaxSize,
			MaxBackups: h.fileOpt.MaxBackups,
			MaxAge:     h.fileOpt.MaxAge,
			LocalTime:  h.fileOpt.LocalTime,
			Compress:   h.fileOpt.Compress,
		}
		ioWriter = io.MultiWriter(h.w, fileWriter)
	}

	slogOpt := &slog.HandlerOptions{
		AddSource:   true,
		Level:       h.level,
		ReplaceAttr: replaceAttr,
	}

	if h.format == "json" {
		h.Handler = slog.NewJSONHandler(ioWriter, slogOpt)
	} else {
		h.Handler = slog.NewTextHandler(ioWriter, slogOpt)
	}

	return h
}

// makeRecord slog默认只能把attr放到末尾, WithGroup时, rid等会有嵌套, 不符合预期, 这里自定义控制展示顺序
func (h *contextualHandler) makeRecord(ctx context.Context, r slog.Record) slog.Record {
	newR := slog.Record{
		Time:    r.Time,
		Level:   r.Level,
		Message: r.Message,
		PC:      r.PC,
	}

	// content中的attr
	if attrs, ok := ctx.Value(attrCtxKey).([]slog.Attr); ok {
		newR.AddAttrs(attrs...)
	}

	finalAttrs := make([]slog.Attr, 0, r.NumAttrs())

	// record中的attr
	r.Attrs(func(a slog.Attr) bool {
		finalAttrs = append(finalAttrs, a)
		return true
	})

	// handler中的attr
	for g := h.attr; g != nil; g = g.next {
		if g.group != "" {
			finalAttrs = []slog.Attr{{
				Key:   g.group,
				Value: slog.GroupValue(finalAttrs...),
			}}
		} else {
			finalAttrs = append(g.attrs, finalAttrs...)
		}
	}

	newR.AddAttrs(finalAttrs...)
	return newR
}

// Handle 支持自定义属性等
func (h *contextualHandler) Handle(ctx context.Context, r slog.Record) error {
	r = h.makeRecord(ctx, r)

	// slog的入口
	d, ok := ctx.Value(depthCtxKey).(*int)
	if !ok {
		return h.Handler.Handle(ctx, r)
	}

	// 自定义Logger入口
	depth := 0
	if d != nil {
		// 自定义depth, 需要减去全局函数的一次调用
		depth = *d - 1
	}

	// skip ref https://github.com/golang/go/blob/master/src/log/slog/logger.go#L94
	var pcs [1]uintptr
	runtime.Callers(7+depth, pcs[:])
	(&r).PC = pcs[0]

	return h.Handler.Handle(ctx, r)
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments
func (h *contextualHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newH := *h
	newH.attr = newH.attr.WithAttrs(attrs)
	return &newH
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
func (h *contextualHandler) WithGroup(name string) slog.Handler {
	newH := *h
	newH.attr = newH.attr.WithGroup(name)
	return &newH
}

// replaceAttr source 格式化为 file:line 格式
func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	// 自定义Source
	if a.Key == slog.SourceKey {
		src, ok := a.Value.Any().(*slog.Source)
		if !ok {
			return a
		}

		a.Value = slog.StringValue(filepath.Base(src.File) + ":" + strconv.Itoa(src.Line))
		return a
	}

	// 自定义日志等级
	if a.Key == slog.LevelKey {
		level, ok := a.Value.Any().(slog.Level)
		if !ok {
			return a
		}

		switch {
		case level < slog.LevelDebug:
			a.Value = slog.StringValue("TRACE")
		case level < slog.LevelInfo:
			a.Value = slog.StringValue("DEBUG")
		case level < slog.LevelWarn:
			a.Value = slog.StringValue("INFO")
		case level < slog.LevelError:
			a.Value = slog.StringValue("WARNING")
		default:
			a.Value = slog.StringValue("ERROR")
		}
	}

	return a
}

// groupOrAttrs holds either a group name or a list of slog.Attrs
type groupOrAttrs struct {
	group string        // group name if non-empty
	attrs []slog.Attr   // attrs if non-empty
	next  *groupOrAttrs // parent
}

// WithGroup ...
func (g *groupOrAttrs) WithGroup(name string) *groupOrAttrs {
	// Empty-name groups are inlined as if they didn't exist
	if name == "" {
		return g
	}
	return &groupOrAttrs{
		group: name,
		next:  g,
	}
}

// WithAttrs ...
func (g *groupOrAttrs) WithAttrs(attrs []slog.Attr) *groupOrAttrs {
	if len(attrs) == 0 {
		return g
	}
	return &groupOrAttrs{
		attrs: attrs,
		next:  g,
	}
}
