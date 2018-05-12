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

var constTests = []struct {
	exprStr, expr, want string
}{
	{exprStr: "ErrQuote", expr: ErrQuote, want: internal.ErrQuote},
	{exprStr: "ErrSep", expr: ErrSep, want: internal.ErrSep},
}

// Ensure that error constants don't diverge.
func TestConst(t *testing.T) {
	for _, tt := range constTests {
		if tt.expr != tt.want {
			t.Errorf("%s = %q; want %q", tt.exprStr, tt.expr, tt.want)
		} else {
			t.Logf("%s = %q", tt.exprStr, tt.expr)
		}
	}
}

func colonToSep(list List) List {
	return List(strings.Replace(string(list), ":", string(ListSeparator), -1))
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
	tests := newTests
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
			t.Errorf("New(%q) = %v, %v; want %#q, nil",
				tt.filepaths, l, err, exp)
		case !tt.ok && err == nil:
			t.Errorf("New(%q) = %#q, %v; want error", tt.filepaths, l, err)
		case l != exp:
			t.Errorf("New(%q) = %#q, %v; want %#q, nil",
				tt.filepaths, l, err, exp)
		default:
			t.Logf("New(%q) = %q, %v", tt.filepaths, l, err)
		}
	}
}

type appendToTest struct {
	list      List
	filepaths []string
	appended  List
	prepended List
}

var appendToTests = []appendToTest{
	{"", []string{}, "", ""},
	{"", []string{""}, ":", ":"},
	{"", []string{"c"}, "c", "c"},
	{"", []string{"c", "d"}, "c:d", "c:d"},
	{"a", []string{}, "a", "a"},
	{"a", []string{""}, "a:", ":a"},
	{"a", []string{"c"}, "a:c", "c:a"},
	{"a", []string{"c", "d"}, "a:c:d", "c:d:a"},
	{":", []string{""}, ":", ":"},
	{":", []string{"c"}, ":c", "c:"},
	{":", []string{"c", "d"}, ":c:d", "c:d:"},
	{"a:b", []string{}, "a:b", "a:b"},
	{"a:b", []string{""}, "a:b:", ":a:b"},
	{"a:b", []string{"c"}, "a:b:c", "c:a:b"},
	{"a:b", []string{"c", "d"}, "a:b:c:d", "c:d:a:b"},
	{"a:", []string{}, "a:", "a:"},
	{"a:", []string{""}, "a:", ":a:"},
	{"a:", []string{"c"}, "a::c", "c:a:"},
	{"a:", []string{"c", "d"}, "a::c:d", "c:d:a:"},
	{":b", []string{}, ":b", ":b"},
	{":b", []string{""}, ":b:", ":b"},
	{":b", []string{"c"}, ":b:c", "c::b"},
	{":b", []string{"c", "d"}, ":b:c:d", "c:d::b"},
	{"::", []string{}, "::", "::"},
	{"::", []string{""}, "::", "::"},
	{"::", []string{"c"}, "::c", "c::"},
	{"::", []string{"c", "d"}, "::c:d", "c:d::"},
	{"a::", []string{}, "a::", "a::"},
	{"a::", []string{""}, "a::", ":a::"},
	{"a::", []string{"c"}, "a::c", "c:a::"},
	{"a::", []string{"c", "d"}, "a::c:d", "c:d:a::"},
	{"::b", []string{}, "::b", "::b"},
	{"::b", []string{""}, "::b:", "::b"},
	{"::b", []string{"c"}, "::b:c", "c::b"},
	{"::b", []string{"c", "d"}, "::b:c:d", "c:d::b"},
}

func TestAppendTo(t *testing.T) {
	for _, tt := range appendToTests {
		testAppendToCase(t, tt)
	}
}

func testAppendToCase(t *testing.T, tt appendToTest) {
	appended, aerr := AppendTo(colonToSep(tt.list), tt.filepaths...)
	switch {
	case aerr != nil:
		t.Errorf("AppendTo(%q, %q) = %#q, %v; want equivalent to %q, nil",
			tt.list, tt.filepaths, appended, aerr, tt.appended)
	case !equiv(Split(colonToSep(tt.appended)), Split(appended)):
		t.Errorf("AppendTo(%q, %q) = %#q, %v; want equivalent to %q, nil",
			tt.list, tt.filepaths, appended, aerr, tt.appended)
	default:
		t.Logf("AppendTo(%q, %q) = %#q, %v",
			tt.list, tt.filepaths, appended, aerr)
	}
	prepended, perr := PrependTo(colonToSep(tt.list), tt.filepaths...)
	switch {
	case perr != nil:
		t.Errorf("PrependTo(%q, %q) = %#q, %v; want equivalent to %q, nil",
			tt.list, tt.filepaths, prepended, perr, tt.prepended)
	case !equiv(Split(colonToSep(tt.prepended)), Split(prepended)):
		t.Errorf("PrependTo(%q, %q) = %#q, %v; want equivalent to %q, nil",
			tt.list, tt.filepaths, prepended, perr, tt.prepended)
	default:
		t.Logf("PrependTo(%q, %q) = %#q, %v",
			tt.list, tt.filepaths, prepended, perr)
	}
}

func TestAppendToCloseQuote(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows-only test")
	}
	l := List(`a"a`)
	fp := "b"
	{ // AppendTo
		want := List(`a"a";b`)
		got, err := AppendTo(l, fp)
		if err != nil || got != want {
			t.Errorf("AppendTo(%#q, %q) = %#q, %v; want %#q, nil", l, fp, got, err, want)
		} else {
			t.Logf("AppendTo(%#q, %q) = %#q, %v", l, fp, got, err)
		}
	}
	{ // PrependTo
		want := List(`b;a"a"`)
		got, err := PrependTo(l, fp)
		if err != nil || got != want {
			t.Errorf("PrependTo(%#q, %q) = %#q, %v; want %#q, nil", l, fp, got, err, want)
		} else {
			t.Logf("PrependTo(%#q, %q) = %#q, %v", l, fp, got, err)
		}
	}
}

func TestAppendToInvalidFilepath(t *testing.T) {
	if invalidFilepath == "" {
		t.Skip("no invalid filepath on this OS")
	}
	{
		got, err := AppendTo(colonToSep("a:b"), invalidFilepath)
		if err == nil {
			t.Errorf("AppendTo(%#q, %q) = %#q, %v, want error", colonToSep("a:b"),
				invalidFilepath, got, err)
		} else {
			t.Logf("AppendTo(%#q, %q) = %#q, %v", colonToSep("a:b"),
				invalidFilepath, got, err)
		}
	}
	{
		got, err := PrependTo(colonToSep("a:b"), invalidFilepath)
		if err == nil {
			t.Errorf("PrependTo(%#q, %q) = %#q, %v, want error", colonToSep("a:b"),
				invalidFilepath, got, err)
		} else {
			t.Logf("PrependTo(%#q, %q) = %#q, %v", colonToSep("a:b"),
				invalidFilepath, got, err)
		}
	}
}

func TestMustOK(t *testing.T) {
	want := colonToSep("a:b:c")
	got := Must(AppendTo(colonToSep("a:b"), "c"))
	if equiv(Split(got), Split(want)) {
		t.Logf("Must(AppendTo(%#q, %q)) = %#q", colonToSep("a:b"), "c", got)
	} else {
		t.Logf("Must(AppendTo(%#q, %q)) = %#q, want equivalent to %#q",
			colonToSep("a:b"), "c", got, want)
	}
}

func TestMustPanic(t *testing.T) {
	if invalidFilepath == "" {
		t.Skip("no invalid filepath on this OS")
	}
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Must(AppendTo(%#q, %q)): panic %v", colonToSep("a:b"),
				invalidFilepath, r)
		}
	}()
	got := Must(AppendTo(colonToSep("a:b"), invalidFilepath))
	t.Logf("Must(AppendTo(%#q, %q)) = %v, want panic", colonToSep("a:b"),
		invalidFilepath, got)
}
