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
	"fmt"
	"testing"

	"github.com/retr0h/go-gilt/repository"
	"github.com/retr0h/go-gilt/test/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GitTestSuite struct {
	suite.Suite
	g *Git
	r repository.Repository
}

func (suite *GitTestSuite) SetupTest() {
	suite.g = NewGit(false)
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

func (suite *GitTestSuite) TestCloneReturnsError() {
	anon := func() error {
		err := suite.g.clone(suite.r)
		assert.Error(suite.T(), err)

		return err
	}

	MockRunCommandErrorsOn("git", anon)
}

func (suite *GitTestSuite) TestClone() {
	anon := func() error {
		err := suite.g.clone(suite.r)
		assert.NoError(suite.T(), err)

		return err
	}

	got := MockRunCommand(anon)
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

	MockRunCommandErrorsOn("git", anon)
}

func (suite *GitTestSuite) TestReset() {
	anon := func() error {
		err := suite.g.reset(suite.r)
		assert.NoError(suite.T(), err)

		return err
	}

	got := MockRunCommand(anon)
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
