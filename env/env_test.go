// Copyright 2015 Péter Surányi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env_test

import (
	"testing"

	"gopkg.in/pathlist.v0"
	"gopkg.in/pathlist.v0/env"
)

var sliceTests = [...]struct {
	env      []string
	key      string
	get, set pathlist.List
}{
	{env: nil, key: "PATH", get: "", set: "/bin"},
	{env: []string{}, key: "PATH", get: "", set: "/bin"},
	{env: []string{"PATH="}, key: "PATH", get: "", set: "/bin"},
	{env: []string{"PATH=/sbin"}, key: "PATH", get: "/sbin", set: "/bin"},
	{env: []string{"VAR=val"}, key: "PATH", get: "", set: "/bin"},
	{env: []string{"VAR=val", "PATH="}, key: "PATH", get: "", set: "/bin"},
	{env: []string{"VAR=val", "PATH=/sbin"}, key: "PATH", get: "/sbin", set: "/bin"},
	{env: []string{"PATH=", "VAR=val"}, key: "PATH", get: "", set: "/bin"},
	{env: []string{"PATH=/sbin", "VAR=val"}, key: "PATH", get: "/sbin", set: "/bin"},
}

func TestSlice(t *testing.T) {
	for _, tt := range sliceTests {
		got1 := env.Slice(tt.env, tt.key)
		if got1 != tt.get {
			t.Errorf("Slice(%q, %q) = %q; want %q", tt.env, tt.key, got1, tt.get)
		} else {
			t.Logf("Slice(%q, %q) = %q", tt.env, tt.key, got1)
		}
		env2 := env.SetSlice(tt.env, tt.key, tt.set)
		got2 := env.Slice(env2, tt.key)
		if got2 != tt.set {
			t.Errorf("env := SetSlice(%q, %q, %q); Slice(env, %q) = %q (env == %q); want %q",
				tt.env, tt.key, tt.set, tt.key, got2, env2, tt.set)
		} else {
			t.Logf("env := SetSlice(%q, %q, %q); Slice(env, %q) = %q (env == %q)",
				tt.env, tt.key, tt.set, tt.key, got2, env2)
		}
	}
}
