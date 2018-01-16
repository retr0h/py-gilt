// Copyright (c) 2018 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package util_test

import (
	"os/user"
	"path/filepath"
	"testing"

	capturer "github.com/kami-zh/go-capturer"
	"github.com/retr0h/go-gilt/util"
	"github.com/stretchr/testify/assert"
)

func TestPrintError(t *testing.T) {
	got := capturer.CaptureStderr(func() {
		util.PrintError("foo")
	})
	want := "\x1b[31mERROR\x1b[0m: foo\n"

	assert.Equal(t, want, got)
}

func TestPrintErrorAndExit(t *testing.T) {
	// Located in utils_test for mocking.
}

func TestExpandUserReturnsError(t *testing.T) {
	// Located in utils_test for mocking.
}

func TestExpandUser(t *testing.T) {
	got, err := util.ExpandUser("~/foo/bar")
	usr, _ := user.Current()
	want := filepath.Join(usr.HomeDir, "foo", "bar")

	assert.Equal(t, want, got)
	assert.NoError(t, err)
}

func TestExpandUserWithFullPath(t *testing.T) {
	got, err := util.ExpandUser("/foo/bar")
	want := "/foo/bar"

	assert.Equal(t, want, got)
	assert.NoError(t, err)
}

func TestRunCommandReturnsError(t *testing.T) {
	err := util.RunCmd(false, "false")

	assert.Error(t, err)
}

func TestRunCommandPrintsStreamingStdout(t *testing.T) {
	got := capturer.CaptureStdout(func() {
		err := util.RunCmd(true, "echo", "-n", "foo")
		assert.NoError(t, err)
	})
	want := "COMMAND: \x1b[30;41mecho -n foo\x1b[0m\nfoo"

	assert.Equal(t, want, got)
}

func TestRunCommandPrintsStreamingStderr(t *testing.T) {
	got := capturer.CaptureStderr(func() {
		err := util.RunCmd(true, "cat", "foo")
		assert.Error(t, err)
	})
	want := "cat: foo: No such file or directory\n"

	assert.Equal(t, want, got)
}
