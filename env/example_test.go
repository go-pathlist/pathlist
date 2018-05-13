// Copyright 2015 Péter Surányi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/pathlist.v0"
	"gopkg.in/pathlist.v0/env"
)

const helloworld = `
package main

import "fmt"

func main() {
	fmt.Println("Hello World!")
}
`

// This example demonstrates setting PATH and GOPATH, both for the current
// process (SetPath) and a child process (SetSlice).
// It first creates a new Go workspace for building a Hello World executable,
// adding the workspace to GOPATH when invoking the go tool.
// Then it invokes the executable built, adding the workspace bin directory to
// PATH.
func Example() {
	// use a separate function as deferreds are not invoked with log.Fatal
	out, err := doExample()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))

	// Output: Hello World!
}

func doExample() ([]byte, error) {
	// create workspace
	wkspc, err := ioutil.TempDir(os.TempDir(), "env_test.Example_")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(wkspc)

	// place source file
	hellodir := filepath.Join(wkspc, "src", "hello")
	if err := os.MkdirAll(hellodir, 0777); err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(filepath.Join(hellodir, "hello.go"),
		[]byte(helloworld), 0666); err != nil {

		return nil, err
	}

	// invoke go install
	newgopath, err := pathlist.PrependTo(env.Gopath(), wkspc)
	if err != nil {
		return nil, err
	}
	newenv := env.SetSlice(os.Environ(), env.VarGopath, newgopath)
	goinst := exec.Command("go", "install", "hello")
	goinst.Env = newenv
	if stdouterr, err := goinst.CombinedOutput(); err != nil {
		log.Print(string(stdouterr))
		return nil, err
	}

	// invoke executable
	bindir := filepath.Join(wkspc, "bin") // should be created by go install
	oldpath := env.Path()
	newpath, err := pathlist.PrependTo(oldpath, bindir)
	if err != nil {
		return nil, err
	}
	env.SetPath(newpath)
	defer env.SetPath(oldpath)
	hello := exec.Command("hello")
	stdouterr, err := hello.CombinedOutput()
	if err != nil {
		log.Print(string(stdouterr))
		return nil, err
	}

	return stdouterr, nil
}
