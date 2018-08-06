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

package git

import (
	"errors"
	"os/exec"
	"strings"
)

// MockRunCommandImpl records and return the arguments passed to RunCommand as
// a string.  An errString may be passed to force the command to fail with
// an error, if the command executed contains the errString.  Useful, when
// mocking functions with multiple calls to RunCommand.
// This seems super wrong to have a test function living in the `git` package.
func MockRunCommandImpl(errString string, f func() error) []string {
	var got []string

	originalRunCommand := RunCommand
	RunCommand = func(debug bool, name string, args ...string) error {
		cmd := exec.Command(name, args...)
		cmdString := strings.Join(cmd.Args, " ")

		if errString == "" {
			got = append(got, cmdString)
			return nil
		}

		if strings.Contains(cmdString, errString) {
			return errors.New("RunCommand had an error")
		}

		// NOTE(retr0h): Never hit this path since this is only used by unit tests
		// and we keep an eye on our returns.  However, the lack of this path causes
		// our code coverage to drop.
		return nil
	}
	defer func() { RunCommand = originalRunCommand }()

	f()

	return got
}

// MockedRunCommand is sugar around MockRunCommandImpl, and returns
// a string with the arguments passed to RunCommand.
func MockRunCommand(f func() error) []string {
	return MockRunCommandImpl("", f)
}

// MockedRunCommandErrors is sugar around MockedRunCommandImpl and
// returns an error when invoked.
func MockRunCommandErrorsOn(errCmd string, f func() error) {
	MockRunCommandImpl(errCmd, f)
}
