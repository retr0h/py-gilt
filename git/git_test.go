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
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/retr0h/go-gilt/repository"
	"github.com/retr0h/go-gilt/test/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GitTestSuite struct {
	suite.Suite
	g *Git
	c *Cloneable
	r repository.Repository
}

func (suite *GitTestSuite) SetupTest() {
	g := NewGit(false)
	c := NewCloneable(false)
	c.GC = g // Real Git Implementation

	suite.g = g
	suite.c = c
	suite.r = repository.Repository{
		URL:     "https://example.com/user/repo.git",
		Version: "abc1234",
		Dst:     "path/user.repo",
		GiltDir: testutil.CreateTempDirectory(),
	}
}

func (suite *GitTestSuite) TearDownTest() {
	testutil.RemoveTempDirectory(suite.r.GiltDir)
}

func (suite *GitTestSuite) mockRunCommand(f func()) []string {
	var got []string

	originalRunCommand := RunCommand
	RunCommand = func(debug bool, name string, args ...string) error {
		cmd := exec.Command(name, args...)
		got = append(got, strings.Join(cmd.Args, " "))

		return nil
	}
	defer func() { RunCommand = originalRunCommand }()

	f()

	return got
}

func (suite *GitTestSuite) mockRunCommandErrors(f func() error) {
	originalRunCommand := RunCommand
	RunCommand = func(debug bool, name string, args ...string) error {
		return errors.New("RunCommand had an error")
	}
	defer func() { RunCommand = originalRunCommand }()

	f()
}

func (suite *GitTestSuite) TestCloneReturnsError() {
	anon := func() error {
		err := suite.g.clone(suite.r)
		assert.Error(suite.T(), err)

		return err
	}

	suite.mockRunCommandErrors(anon)
}

func (suite *GitTestSuite) TestClone() {
	anon := func() { suite.g.clone(suite.r) }

	got := suite.mockRunCommand(anon)
	want := []string{
		fmt.Sprintf("git clone https://example.com/user/repo.git %s/https---example.com-user-repo.git-abc1234",
			suite.r.GiltDir),
	}

	assert.Equal(suite.T(), want, got)
}

func (suite *GitTestSuite) TestResetReturnsError() {
	anon := func() error {
		err := suite.g.reset(suite.r)
		assert.Error(suite.T(), err)

		return err
	}

	suite.mockRunCommandErrors(anon)
}

func (suite *GitTestSuite) TestReset() {
	anon := func() { suite.g.reset(suite.r) }

	got := suite.mockRunCommand(anon)
	want := []string{
		fmt.Sprintf("git -C %s/https---example.com-user-repo.git-abc1234 reset --hard abc1234",
			suite.r.GiltDir),
	}

	assert.Equal(suite.T(), want, got)
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestGitTestSuite(t *testing.T) {
	suite.Run(t, new(GitTestSuite))
}
