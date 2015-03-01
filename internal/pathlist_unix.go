// Copyright 2013 Péter Surányi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package internal

import (
	"strings"
)

func NewElem(fp string) (string, error) {
	if strings.ContainsRune(fp, ListSeparator) {
		return "", Error{Cause_: ErrSep, Filepath_: fp}
	}
	return fp, nil
}

func CloseQuote(el string) string {
	// no quoting on Unix
	return el
}
