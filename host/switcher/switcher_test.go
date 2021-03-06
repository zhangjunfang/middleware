// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package switcher

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/issue9/assert"
)

var (
	f1 = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(1)
	}

	f2 = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(2)
	}
)

var (
	h1 = http.HandlerFunc(f1)
	h2 = http.HandlerFunc(f2)
)

func TestNew(t *testing.T) {
	a := assert.New(t)

	switcher := New()
	a.NotNil(switcher)

	switcher.AddHost(h1, "caixw.io")
	switcher.AddHost(h2, "caixw.oi")

	// h1
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "https://caixw.io/test", nil)
	a.NotNil(w).NotNil(r)
	switcher.ServeHTTP(w, r)
	a.Equal(w.Code, 1)

	// h2
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "http://caixw.oi/test", nil)
	a.NotNil(w).NotNil(r)
	switcher.ServeHTTP(w, r)
	a.Equal(w.Code, 2)

	// 未带域名
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "http://not.exsits/test", nil)
	switcher.ServeHTTP(w, r)
	a.NotNil(w).NotNil(r)
	a.Equal(w.Code, http.StatusNotFound)
}
