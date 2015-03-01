// Copyright 2013 Péter Surányi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"strings"
)

func NewElem(fp string) (string, error) {
	if strings.ContainsRune(fp, '"') {
		return "", Error{Cause_: ErrQuote, Filepath_: fp}
	}
	if strings.ContainsRune(fp, ListSeparator) {
		return `"` + fp + `"`, nil
	}
	return fp, nil
}

func CloseQuote(el string) string {
	c := strings.Count(el, string(ListSeparator))
	if c%2 != 0 {
		return el + `"`
	}
	return el
}
