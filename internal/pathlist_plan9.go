// Copyright 2013 Péter Surányi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

func NewElem(fp string) (string, error) {
	return fp, nil
}

func CloseQuote(el string) string {
	// no quoting on Plan9
	return el
}
