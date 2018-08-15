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

package git_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/retr0h/go-gilt/git"
	"github.com/retr0h/go-gilt/repository"
	"github.com/retr0h/go-gilt/test/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GitTestSuite struct {
	suite.Suite
	g *git.Git
	r repository.Repository
}

func (suite *GitTestSuite) SetupTest() {
	suite.g = git.NewGit(false)
	suite.r = repository.Repository{
		Git:     "https://example.com/user/repo.git",
		Version: "abc1234",
		DstDir:  "path/user.repo",
		GiltDir: testutil.CreateTempDirectory(),
	}
}

func (suite *GitTestSuite) TearDownTest() {
	testutil.RemoveTempDirectory(suite.r.GiltDir)
}

func (suite *GitTestSuite) TestCloneAlreadyExists() {
	cloneDir := filepath.Join(suite.r.GiltDir, "https---example.com-user-repo.git-abc1234")
	if _, err := os.Stat(cloneDir); os.IsNotExist(err) {
		os.Mkdir(cloneDir, 0755)
	}

	suite.g.Clone(suite.r)

	defer os.RemoveAll(cloneDir)
}

func (suite *GitTestSuite) TestCloneErrorsOnCloneReturnsError() {
	anon := func() error {
		err := suite.g.Clone(suite.r)
		assert.Error(suite.T(), err)

		return err
	}

	git.MockRunCommandErrorsOn("clone", anon)
}

func (suite *GitTestSuite) TestCloneErrorsOnResetReturnsError() {
	anon := func() error {
		err := suite.g.Clone(suite.r)
		assert.Error(suite.T(), err)

		return err
	}

	git.MockRunCommandErrorsOn("reset", anon)
}

func (suite *GitTestSuite) TestClone() {
	anon := func() error {
		err := suite.g.Clone(suite.r)
		assert.NoError(suite.T(), err)

		return err
	}

	got := git.MockRunCommand(anon)
	want := []string{
		fmt.Sprintf("git clone https://example.com/user/repo.git %s/https---example.com-user-repo.git-abc1234",
			suite.r.GiltDir),
		fmt.Sprintf("git -C %s/https---example.com-user-repo.git-abc1234 reset --hard abc1234",
			suite.r.GiltDir),
	}

	assert.Equal(suite.T(), want, got)
}

func (suite *GitTestSuite) TestCheckoutIndexFailsFilepathAbsReturnsError() {
	originalFilepathAbs := git.FilePathAbs
	git.FilePathAbs = func(string) (string, error) {
		return "", errors.New("Failed filepath.Abs")
	}
	defer func() { git.FilePathAbs = originalFilepathAbs }()

	err := suite.g.CheckoutIndex(suite.r)
	assert.Error(suite.T(), err)
}

func (suite *GitTestSuite) TestCheckoutIndexFailsCheckoutIndexReturnsError() {
	anon := func() error {
		err := suite.g.CheckoutIndex(suite.r)
		assert.Error(suite.T(), err)

		return err
	}

	git.MockRunCommandErrorsOn("git", anon)
}

func (suite *GitTestSuite) TestCheckoutIndex() {
	anon := func() error {
		err := suite.g.CheckoutIndex(suite.r)
		assert.NoError(suite.T(), err)

		return err
	}

	dstDir, _ := git.FilePathAbs(suite.r.DstDir)
	got := git.MockRunCommand(anon)
	want := []string{
		fmt.Sprintf("git -C %s/https---example.com-user-repo.git-abc1234 checkout-index --force --all --prefix %s",
			suite.r.GiltDir, (dstDir + string(os.PathSeparator))),
	}

	assert.Equal(suite.T(), want, got)
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestGitTestSuite(t *testing.T) {
	suite.Run(t, new(GitTestSuite))
}
