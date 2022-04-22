package middleware

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"bytes"

	ttlcache "github.com/jellydator/ttlcache/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type ResponseCache struct {
	cache *ttlcache.Cache
}

func NewRespCache() *ResponseCache {
	cache := ttlcache.NewCache()

	cache.SetTTL(time.Duration(5 * time.Second))

	return &ResponseCache{
		cache: cache,
	}
}

// https://github.com/labstack/echo/blob/master/middleware/body_dump.go was used as a reference
// Thank you for your work!
/*
	The MIT License (MIT)

	Copyright (c) 2021 LabStack

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.
*/

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (j *ResponseCache) RespCache(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		id := c.Param(UserIDKey)

		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Error(err)
		}

		cacheKey := fmt.Sprintf("%d.%s?%s", idInt, c.Request().URL.Path, c.QueryString())

		if c.Request().Method == "GET" {
			resp, err := j.cache.Get(cacheKey)
			if err != ttlcache.ErrNotFound {
				resp, ok := resp.(*bytes.Buffer)
				if !ok {
					// dangerous FIXME
					log.Fatal("could not cast to echo response")
				}

				log.Infof("CACHE HIT: %v", resp)

				c.Response().Write(resp.Bytes())
				return nil
			} else {
				// cache miss
				resBody := new(bytes.Buffer)
				mw := io.MultiWriter(c.Response().Writer, resBody)
				writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
				c.Response().Writer = writer

				err := next(c)
				if err != nil {
					return err
				}

				j.cache.Set(cacheKey, resBody)

				log.Infof("CACHE MISS, SET: %v", resBody)
				return nil

			}
		}

		return next(c)
	}
}
