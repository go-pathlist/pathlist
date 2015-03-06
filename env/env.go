// Copyright 2015 Péter Surányi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package env provides helper functions for manipulating environment variables
// holding filepath lists.
package env

import (
	"os"
	"strings"

	"gopkg.in/pathlist.v0"
)

// VarGopath is the Go workspace path variable name.
const VarGopath = "GOPATH"

// Path gets the OS (shell) specific executable search path.
func Path() pathlist.List {
	return pathlist.List(os.Getenv(VarPath))
}

// SetPath sets the OS (shell) specific executable search path.
func SetPath(l pathlist.List) error {
	return os.Setenv(VarPath, string(l))
}

// Gopath gets the Go workspace path.
func Gopath() pathlist.List {
	return pathlist.List(os.Getenv(VarGopath))
}

// SetGopath sets the Go workspace path.
func SetGopath(l pathlist.List) error {
	return os.Setenv(VarGopath, string(l))
}

// Slice gets the value for key as a pathlist.List from a slice of environment
// variables (as used with os.Environ and os/exec.Cmd.Env).
func Slice(env []string, key string) pathlist.List {
	keyEq := key + "="
	for _, kv := range env {
		if strings.HasPrefix(kv, keyEq) {
			return pathlist.List(kv[len(keyEq):])
		}
	}
	return ""
}

// SetSlice takes a slice of environment variables (as used with os.Environ and
// os/exec.Cmd.Env), and returns a copy of env with key set to list.
func SetSlice(env []string, key string, list pathlist.List) []string {
	keyEq := key + "="
	keyEqVal := keyEq + string(list)
	for i, kv := range env {
		if strings.HasPrefix(kv, keyEq) {
			env := append([]string(nil), env...)
			env[i] = keyEqVal
			return env
		}
	}
	return append(append(make([]string, 0, len(env)+1), env...), keyEqVal)
}
