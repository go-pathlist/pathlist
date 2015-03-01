// Copyright 2013 Péter Surányi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pathlist

import (
	"runtime"
	"strings"
	"testing"

	"gopkg.in/pathlist.v0/internal"
)

var _ Error = internal.Error{}

func assertString(t *testing.T, expr, got, want string) {
	if want != got {
		t.Errorf("fail: %s: got %q, want %q", expr, got, want)
	} else {
		t.Logf("pass: %s: got %q", expr, got)
	}
}

func TestConst(t *testing.T) {
	assertString(t, "ErrQuote", ErrQuote, internal.ErrQuote)
	assertString(t, "ErrSep", ErrSep, internal.ErrSep)
}

func colonToSep(list List) List {
	return List(strings.Replace(string(list), ":", string(ListSeparator), -1))
}

func colonToSepSlice(lists []List) []List {
	ll := make([]List, len(lists))
	for i, l := range lists {
		ll[i] = colonToSep(l)
	}
	return ll
}

func equiv(l1, l2 []string) bool {
	i1, i2 := 0, 0
	for ; i1 < len(l1) && i2 < len(l2); i1, i2 = i1+1, i2+1 {
		// treat sequences of empty strings as equivalent.
		for i1 < len(l1)-1 && l1[i1] == "" && l1[i1+1] == "" {
			i1++
		}
		for i2 < len(l2)-1 && l2[i2] == "" && l2[i2+1] == "" {
			i2++
		}
		if l1[i1] != l2[i2] {
			return false
		}
	}
	return i1 == len(l1) && i2 == len(l2)
}

func nop(b []byte) {
}

func reverse(b []byte) {
	l := len(b)
	for i := 0; i < l/2; i++ {
		b[i], b[l-1-i] = b[l-1-i], b[i]
	}
}

type newTest struct {
	filepaths []string
	ok        bool
	list      List
}

var newTests = []newTest{
	{[]string{}, true, ""},
	{[]string{""}, true, ":"},
	{[]string{"a"}, true, "a"},
	{[]string{"", ""}, true, ":"},
	{[]string{"a", ""}, true, "a:"},
	{[]string{"", "b"}, true, ":b"},
	{[]string{"a", "b"}, true, "a:b"},
}

var newTestsUnixPlan9 = []newTest{
	{[]string{":"}, false, ""},
}

var newTestsWindows = []newTest{
	{[]string{`a;b`}, true, `"a;b"`},
	{[]string{`"a"`}, false, ""},
}

func TestNew(t *testing.T) {
	tests := append([]newTest{}, newTests...)
	switch runtime.GOOS {
	case "darwin", "dragonfly", "freebsd", "linux", "netbsd", "openbsd",
		"plan9", "solaris":
		tests = append(tests, newTestsUnixPlan9...)
	case "windows":
		tests = append(tests, newTestsWindows...)
	}
	for _, tt := range tests {
		exp := colonToSep(tt.list)
		l, err := New(tt.filepaths...)
		switch {
		case tt.ok && err != nil:
			t.Errorf("fail: New(%q): got error %v, want %#q",
				tt.filepaths, err, exp)
		case !tt.ok && err == nil:
			t.Errorf("fail: New(%q): got ok, want error", tt.filepaths)
		case l != exp:
			t.Errorf("fail: New(%q): got %#q, want %#q",
				tt.filepaths, l, exp)
		default:
			t.Logf("pass: New(%q): got %q, %v", tt.filepaths, l, err)
		}
	}
}

type appendToTestTarget struct {
	name    string
	prepare func([]byte)
	appnd   func(List, string) (List, error)
}

var appendToTestTargets = []appendToTestTarget{
	{"AppendTo", nop, AppendTo},
	{"PrependTo", reverse, PrependTo},
}

type appendToTest struct {
	list     List
	filepath string
	appended List
}

var appendToTests = []appendToTest{
	{"", "", ":"},
	{"", "c", "c"},
	{"a", "", "a:"},
	{"a", "c", "a:c"},
	{":", "", ":"},
	{":", "c", ":c"},
	{"a:b", "", "a:b:"},
	{"a:b", "c", "a:b:c"},
	{"a:", "", "a:"},
	{"a:", "c", "a::c"},
	{":b", "", ":b:"},
	{":b", "c", ":b:c"},
	{"::", "", "::"},
	{"::", "c", "::c"},
	{"a::", "", "a::"},
	{"a::", "c", "a::c"},
}

func TestAppendTo(t *testing.T) {
	for _, target := range appendToTestTargets {
		for _, tt := range appendToTests {
			testAppendToCase(t, target, tt)
		}
	}
}

func testAppendToCase(t *testing.T, target appendToTestTarget,
	tt appendToTest) {

	listBytes := []byte(colonToSep(tt.list))
	filepathBytes := []byte(tt.filepath)
	expBytes := []byte(colonToSep(tt.appended))
	target.prepare(listBytes)
	target.prepare(filepathBytes)
	target.prepare(expBytes)
	list := List(listBytes)
	filepath := string(filepathBytes)
	exp := List(expBytes)
	appended, err := target.appnd(list, filepath)
	switch {
	case err != nil:
		t.Errorf("fail: %s(%q, %q), got error %v, want equivalent to %q",
			target.name, list, filepath, err, exp)
	case !equiv(Split(exp), Split(appended)):
		t.Errorf("fail: %s(%q, %q), got %q, want equivalent to %q",
			target.name, list, filepath, appended, exp)
	default:
		t.Logf("pass: %s(%q, %q), got %q",
			target.name, list, filepath, appended)
	}
}

func TestAppendToCloseQuote(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows-only test")
	}
	l := List(`a"a`)
	fp := "b"
	{ // AppendTo
		want := List(`a"a":b`)
		got, err := AppendTo(l, fp)
		if err != nil || got != want {
			t.Errorf("fail: AppendTo(%#q, %q) = %#q, %v, want %#q, nil", l, fp, got, err, want)
		} else {
			t.Logf("pass: AppendTo(%#q, %q) = %#q, %v", l, fp, got, err)
		}
	}
	{ // PrependTo
		want := List(`b:a"a"`)
		got, err := PrependTo(l, fp)
		if err != nil || got != want {
			t.Errorf("fail: PrependTo(%#q, %q) = %#q, %v, want %#q, nil", l, fp, got, err, want)
		} else {
			t.Logf("pass: PrependTo(%#q, %q) = %#q, %v", l, fp, got, err)
		}
	}
}
