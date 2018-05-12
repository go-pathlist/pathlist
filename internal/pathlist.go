// Copyright 2013 Péter Surányi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// pathlist*.go contain implementations for the pathlist package.

// Package internal implements internal logic for the pathlist package.
//
// Parameter names refer to "type":
//  fp: filepath (single, raw/unquoted)
//  l:  list (0 or more elems, potentially quoted)
//  e:  elem (single, potentially quoted)
//  el: e or l
//
// TODO(speter): sane behavior on elems containing unclosed quotes
package internal

import (
	"fmt"
	"os"
	"strings"

	"path/filepath"
)

const (
	ErrSep   = "filepath must not contain ListSeparator" // Unix only
	ErrQuote = "filepath must not be quoted"             // Windows only
)

const ListSeparator = os.PathListSeparator

const (
	listsep  = string(ListSeparator)
	listsep2 = listsep + listsep
)

type Error struct {
	Cause_    string
	Filepath_ string
}

// Error implements the pathlist.Error interface.
func (e Error) Error() string {
	return fmt.Sprintf("pathlist: %s; filepath: %#q", e.Cause_, e.Filepath_)
}

func (e Error) Cause() string {
	return e.Cause_
}

func (e Error) Filepath() string {
	return e.Filepath_
}

func NewList(e ...string) string {
	if len(e) == 0 {
		return ""
	}
	if len(e) == 1 && e[0] == "" {
		return listsep
	}
	return strings.Join(e, listsep)
}

func Filepaths(l string) []string {
	if l == listsep {
		return []string{""}
	}
	return filepath.SplitList(l)
}

func Append(l, e string) string {
	l = CloseQuote(l)
	if (l == "" && e != "") || l == listsep ||
		strings.HasSuffix(l, listsep2) {
		return l + e
	}
	return l + listsep + e
}

func Prepend(l, e string) string {
	l = CloseQuote(l)
	if (l == "" && e != "") || l == listsep ||
		strings.HasPrefix(l, listsep2) {
		return e + l
	}
	return e + listsep + l
}
