// Copyright 2013 Péter Surányi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package internal

import (
	"strings"
)

func NewElem(filepath string) (string, error) {
	if strings.ContainsRune(filepath, ListSeparator) {
		return "", Error{Cause_: ErrSep, Filepath_: filepath}
	}
	return filepath, nil
}
