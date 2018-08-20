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

package rpc

import (
	"compress/flate"
	"io"
)

type compressor interface {
	flushWriter
	io.Reader
}

type flushWriter interface {
	io.Writer
	Flush() error
}

type Compressor struct {
	zr    io.Reader
	zw    flushWriter
	flush func() error
}

func newCompressor(r io.Reader, w flushWriter, compress string) (*Compressor, error) {
	var zr io.Reader
	var zw flushWriter
	var err error
	switch compress {
	case "deflate":
		zr = flate.NewReader(r)
		zw, err = flate.NewWriter(w, flate.DefaultCompression)
		if err != nil {
			return nil, err
		}
	default:
		zr = r
		zw = w
	}

	return &Compressor{
		zr:    zr,
		zw:    zw,
		flush: w.Flush,
	}, nil
}

func (c *Compressor) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}
func (c *Compressor) Write(p []byte) (n int, err error) {
	return c.zw.Write(p)
}
func (c *Compressor) Flush() (err error) {
	if err = c.zw.Flush(); err != nil {
		return err
	}
	return c.flush()
}
