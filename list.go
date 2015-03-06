// Copyright 2013 Péter Surányi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pathlist manipulates lists of filepaths joined by the OS-specific
// ListSeparator (pathlists), usually found in PATH or GOPATH environment
// variables.
// See also package env ( https://godoc.org/gopkg.in/pathlist.v0/env ) for
// helper functions and examples for using this package.
//
// The package complements the functionality of the standard path/filepath with:
//  - New: create a pathlist from individual filepaths.
//  - AppendTo/PrependTo: extend a pathlist with an individual filepath.
//
// Pathlist handles the quoting/unquoting of filepaths as required on Windows.
// The package uses two separate types to make the context of string parameters
// (single/list, quoted/unquoted) explicit:
//  - string: individual OS-native raw (unquoted) filepath
//  - List:   path list in OS-specific format (potentially quoted on Windows)
//
// This also helps prevent some common mistakes with pathlist manipulation that
// can lead to unexpected behavior and security issues.
// Specifically, a naive implementation could have the following problems:
//  - NaiveJoin("", "/mybin") on Unix would return ":/mybin", which would search
//    the working directory before /mybin, even if "" was meant to represent the
//    empty pathlist.
//  - NaiveJoin("/mnt/C:/tmp/bin", os.Getenv("PATH")) on Unix would prepend two
//    directories "/mnt/C" and "/tmp/bin" to PATH, even if the argument was
//    meant to represent the single filepath "/mnt/C:/tmp/bin".
//
// In addition to using distinct types, filepath arguments are validated.
// On Unix, filepaths containing the separator (':') cannot be used as part of a
// List; these trigger an Error with Cause returning ErrSep.
// On Windows, raw (unquoted) filepaths must be used; an Error with Cause
// returning ErrQuote is issued otherwise.
// Calls returning (List, error) can be wrapped using Must when the validation
// is known to succeed in advance.
//
// Using this package, the two examples above can be written as:
//  pathlist.AppendTo("", "/mybin") // returns "/mybin", nil
//  pathlist.PrependTo(pathlist.OsPath(), "/mnt/C:/tmp/bin") // returns ErrSep
//
// Note that even though the error handling API is common across platforms, the
// behavior is somewhat asymmetric.
// On Unix, ':' is generally a valid (although uncommon) character in filepaths,
// so an error typically indicates input error and the user should be notified.
// In contrast, on Windows all filepaths can be included in a list; an error
// generally means that the caller is mixing quoted and unquoted paths, which
// likely indicates a bug and that the implementation should be fixed.
package pathlist // import "gopkg.in/pathlist.v0"

import (
	"os"
	"path/filepath"

	"gopkg.in/pathlist.v0/internal"
)

// ErrSep and ErrQuote are returned by Error.Cause.
const ( // replicated from internal to keep messages in godoc
	ErrSep   = "filepath must not contain ListSeparator" // Unix only
	ErrQuote = "filepath must not be quoted"             // Windows only
)

// ListSeparator is the OS-specific path list separator.
const ListSeparator = os.PathListSeparator

const (
	listsep  = List(ListSeparator)
	listsep2 = listsep + listsep
)

// Error holds a pathlist handling error.
// Functions in this package return error values implementing this interface.
type Error interface {
	error
	// Cause returns the cause of the error; either ErrSep or ErrQuote.
	Cause() string
	// Filepath returns the offending filepath.
	Filepath() string
}

// List represents a list of zero or more filepaths joined by the OS-specific
// ListSeparator, usually found in PATH or GOPATH environment variables.
//
// On Unix, each list element is the filepath verbatim, and filepaths containing
// the separator (':') cannot be part of a List.
// On Windows, parts or the whole of each list element may be double-quoted, and
// quoting is mandatory when the filepath contains the separator (';').
type List string

// New returns a new List consisting of the given filepaths, or an Error if
// there is an invalid filepath.
func New(filepaths ...string) (List, error) {
	elems := make([]string, len(filepaths))
	for i, fp := range filepaths {
		elem, err := internal.NewElem(fp)
		if err != nil {
			return "", err
		}
		elems[i] = elem
	}
	return List(internal.NewList(elems...)), nil
}

// Split returns the (raw/unquoted) filepaths contained in list.
// Unlike strings.Split, Split returns an empty slice when passed an empty
// string, and returns a single empty filepath when passed the sole
// ListSeparator.
func Split(list List) []string {
	if list == listsep {
		return []string{""}
	}
	return filepath.SplitList(string(list))
}

// AppendTo returns list with filepath appended if valid, or returns an Error.
func AppendTo(list List, filepath string) (List, error) {
	e, err := internal.NewElem(filepath)
	if err != nil {
		return "", err
	}
	return List(internal.Append(string(list), string(e))), nil
}

// PrependTo returns list with filepath prepended if valid, or returns an Error.
func PrependTo(list List, filepath string) (List, error) {
	e, err := internal.NewElem(filepath)
	if err != nil {
		return "", err
	}
	return List(internal.Prepend(string(list), string(e))), nil
}

// Must takes the results from New, AppendTo or PrependTo and returns the List
// on success or panics on error.
func Must(list List, err error) List {
	if err != nil {
		panic(err)
	}
	return list
}
