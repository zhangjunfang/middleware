// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package compress 提供一个支持内容压缩的中间件。
package compress

import (
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"io"
	"log"
	"net/http"

	"github.com/issue9/middleware/compress/accept"
)

// BuildCompressWriter 定义了将一个 io.Writer 声明为具有压缩功能的 io.WriteCloser
type BuildCompressWriter func(w io.Writer) (io.WriteCloser, error)

// NewGzip 表示支持 gzip 格式的压缩
func NewGzip(w io.Writer) (io.WriteCloser, error) {
	return gzip.NewWriter(w), nil
}

// NewDeflate 表示支持 deflate 压缩
func NewDeflate(w io.Writer) (io.WriteCloser, error) {
	return flate.NewWriter(w, flate.DefaultCompression)
}

// NewCompress 用于支持 compress 压缩算法
func NewCompress(w io.Writer) (io.WriteCloser, error) {
	return lzw.NewWriter(w, lzw.LSB, 8), nil
}

type compress struct {
	h      http.Handler
	errlog *log.Logger
	funcs  map[string]BuildCompressWriter
}

// New 构建一个支持压缩的中间件。
// 支持 gzip 或是 deflate 功能的 handler。
// 根据客户端请求内容自动匹配相应的压缩算法，优先匹配 gzip。
//
// NOTE: 经过压缩的内容，可能需要重新指定 Content-Type，系统检测的类型未必正确。
//
// 注意 funcs 键名的大小写。
func New(next http.Handler, errlog *log.Logger, funcs map[string]BuildCompressWriter) http.Handler {
	return &compress{
		h:      next,
		errlog: errlog,
		funcs:  funcs,
	}
}

func (c *compress) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	accepts, err := accept.Parse(r.Header.Get("Accept-Encoding"))
	if err != nil {
		c.errlog.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	var gzw io.WriteCloser
	var accept *accept.Accept
	for _, accept = range accepts {
		// 不支持压缩
		if accept.Value == "identity" || accept.Value == "*" {
			break
		}

		f, found := c.funcs[accept.Value]
		if !found {
			continue
		}

		gzw, err = f(w)
		if err != nil { // 若出错，不压缩，直接返回
			c.errlog.Println(err)
			c.h.ServeHTTP(w, r)
			return
		}
	} // end for

	if gzw == nil { // 不支持的压缩格式
		return
	}

	w.Header().Set("Content-Encoding", accept.Value)
	w.Header().Add("Vary", "Accept-Encoding")
	resp := &response{
		gzw: gzw,
		rw:  w,
	}

	defer gzw.Close() // 只要 gzw!=nil 的，必须会执行到此处。

	// 此处可能 panic，所以得保证在 panic 之前，gzw 变量已经 Close
	c.h.ServeHTTP(resp, r)
}
