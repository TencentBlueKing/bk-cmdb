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
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/pflag"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ctxKey string

const (
	// LevelTrace ...
	LevelTrace = slog.Level(-8)
)

var (
	attrCtxKey    = ctxKey("log.attr")
	depthCtxKey   = ctxKey("log.depth")
	program       = filepath.Base(os.Args[0])
	defaultLogger *contextualLogger
	supportLevel  = []string{"trace", "debug", "info", "warn", "error"}
	supportFormat = []string{"text", "json"}
)

// HandlerOptions ...
type HandlerOptions struct {
	Level   string // 日志等级,可选 trace, debug, info, warn, error
	Format  string // 日志格式,可选 text, json
	Stdout  bool   // 是否输出到标准输出
	LogDir  string // 日志文件目录
	MaxSize int    // 日志文件大小
	MaxNum  int    // 日志文件保留个数
}

// AddFlags adds flags to fs and binds them to options.
func (o *HandlerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Level, "log-level", o.Level,
		fmt.Sprintf("Log filtering level. options: %s", strings.Join(supportLevel, "|")))
	fs.StringVar(&o.Format, "log-format", o.Format,
		fmt.Sprintf("Log format to use. options: %s", strings.Join(supportFormat, "|")))
	fs.BoolVar(&o.Stdout, "log-stdout", o.Stdout, "Log output to stdout")
	fs.StringVar(&o.LogDir, "log-dir", o.LogDir, "If non-empty, write log files in this directory")
	fs.IntVar(&o.MaxSize, "log-max-size", o.MaxSize, "Max size (MB) per file")
	fs.IntVar(&o.MaxNum, "log-max-num", o.MaxNum,
		"Max number of file. The oldest will be removed if there is a extra file created")
}

// Validate opt
func (o *HandlerOptions) Validate() error {
	if !lo.Contains(supportLevel, o.Level) {
		return fmt.Errorf("log-level not valid, options: %s", strings.Join(supportLevel, "|"))
	}
	if !lo.Contains(supportFormat, o.Format) {
		return fmt.Errorf("log-format not valid, options: %s", strings.Join(supportFormat, "|"))
	}
	if o.MaxSize < 0 {
		return fmt.Errorf("log-max-size must be greater than or equal 0")
	}
	if o.MaxNum < 0 {
		return fmt.Errorf("log-max-num must be greater than or equal 0")
	}
	return nil
}

// NewHandlerOptions returns initialized Options
func NewHandlerOptions() *HandlerOptions {
	return &HandlerOptions{
		Level:   getLevelName(slog.LevelInfo),
		Format:  "text",
		Stdout:  true,
		LogDir:  "",
		MaxSize: 500,
		MaxNum:  6,
	}
}

// contextualHandler 支持结构化/上下文的Handler
type contextualHandler struct {
	slog.Handler
	opts  *HandlerOptions // 自定义参数
	attrs *groupOrAttrs   // 自定义attr
}

// NewContextualHandler make a new handler
func NewContextualHandler(opts *HandlerOptions) slog.Handler {
	if opts == nil {
		opts = NewHandlerOptions()
	}

	ioWriter := io.Discard

	if opts.Stdout {
		ioWriter = os.Stdout
	}

	if opts.LogDir != "" {
		fileWriter := &lumberjack.Logger{
			Filename:   filepath.Join(opts.LogDir, fmt.Sprintf("%s.log", program)),
			MaxSize:    opts.MaxSize,
			MaxBackups: opts.MaxNum,
			LocalTime:  true,
			Compress:   false,
		}

		ioWriter = io.MultiWriter(ioWriter, fileWriter)
	}

	slogOpt := &slog.HandlerOptions{
		AddSource:   true,
		Level:       getLevelByName(opts.Level),
		ReplaceAttr: replaceAttr,
	}

	h := &contextualHandler{opts: opts}
	if opts.Format == "json" {
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
	for g := h.attrs; g != nil; g = g.next {
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
	depth, ok := ctx.Value(depthCtxKey).(int)
	if !ok {
		return h.Handler.Handle(ctx, r)
	}

	// 自定义Logger入口
	// skip ref https://github.com/golang/go/blob/master/src/log/slog/logger.go#L94
	var pcs [1]uintptr
	runtime.Callers(6+depth, pcs[:])
	(&r).PC = pcs[0]

	return h.Handler.Handle(ctx, r)
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments
func (h *contextualHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newH := *h
	newH.attrs = newH.attrs.WithAttrs(attrs)
	return &newH
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
func (h *contextualHandler) WithGroup(name string) slog.Handler {
	newH := *h
	newH.attrs = newH.attrs.WithGroup(name)
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

		if level < slog.LevelDebug {
			a.Value = slog.StringValue("TRACE")
		}
		return a
	}

	return a
}

// getLevelerByName human readable logger level
func getLevelByName(name string) slog.Level {
	switch name {
	case "error":
		return slog.LevelError
	case "warn":
		return slog.LevelWarn
	case "info":
		return slog.LevelInfo
	case "debug":
		return slog.LevelDebug
	case "trace":
		return LevelTrace
	default:
		return slog.LevelInfo
	}
}

func getLevelName(level slog.Level) string {
	if level == LevelTrace {
		return "trace"
	}
	return strings.ToLower(level.String())
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
