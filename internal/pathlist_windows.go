// Copyright 2013 Péter Surányi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"strings"
)

func NewElem(fp string) (string, error) {
	if strings.ContainsRune(filepath, '"') {
		return "", Error{Cause_: ErrQuote, Filepath_: filepath}
	}
	if strings.ContainsRune(filepath, ListSeparator) {
		return `"` + filepath + `"`, nil
	}
	return filepath, nil
}

func CloseQuote(el string) string {
	c := strings.Count(el, ListSeparator)
	if c%2 != 0 {
		return el + `"`
	}
	return el
}
