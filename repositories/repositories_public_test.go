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

package repositories_test

import (
	"errors"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/retr0h/go-gilt/git"
	"github.com/retr0h/go-gilt/repositories"
	"github.com/retr0h/go-gilt/test/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RepositoriesTestSuite struct {
	suite.Suite
	r repositories.Repositories
}

func (suite *RepositoriesTestSuite) SetupTest() {
	suite.r = repositories.Repositories{}
	repositories.GiltDir = testutil.CreateTempDirectory()
}

func (suite *RepositoriesTestSuite) TearDownTest() {
	testutil.RemoveTempDirectory(repositories.GiltDir)
}

func (suite *RepositoriesTestSuite) TestUnmarshalYAMLDoesNotParseYAMLAndReturnsError() {
	data := `
---
%foo:
`
	err := suite.r.UnmarshalYAML([]byte(data))
	want := "yaml: line 2: found unexpected non-alphabetical character"

	assert.Equal(suite.T(), want, err.Error())
	assert.Error(suite.T(), err)
}

func (suite *RepositoriesTestSuite) TestUnmarshalYAMLDoesNotValidateYAMLAndReturnsError() {
	data := `
---
foo: bar
`
	err := suite.r.UnmarshalYAML([]byte(data))
	want := errors.New("(root): Invalid type. Expected: array, given: object")

	assert.Equal(suite.T(), want, err)
}

func (suite *RepositoriesTestSuite) TestUnmarshalYAML() {
	data := `
---
- url: https://example.com/user/repo.git
  version: abc1234
  dst: path/user.repo
`
	err := suite.r.UnmarshalYAML([]byte(data))

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "https://example.com/user/repo.git", suite.r.Items[0].URL)
	assert.Equal(suite.T(), "abc1234", suite.r.Items[0].Version)
	assert.Equal(suite.T(), "path/user.repo", suite.r.Items[0].Dst)
}

func (suite *RepositoriesTestSuite) TestUnmarshalYAMLFileReturnsErrorWithMissingFile() {
	suite.r.Filename = "missing.yml"
	err := suite.r.UnmarshalYAMLFile()
	want := "open missing.yml: no such file or directory"

	assert.Equal(suite.T(), want, err.Error())
	assert.Error(suite.T(), err)
}

func (suite *RepositoriesTestSuite) TestUnmarshalYAMLFile() {
	suite.r.Filename = path.Join("..", "test", "gilt.yml")
	suite.r.UnmarshalYAMLFile()

	assert.NotNil(suite.T(), suite.r.Items[0].URL)
	assert.NotNil(suite.T(), suite.r.Items[0].Version)
	assert.NotNil(suite.T(), suite.r.Items[0].Dst)
}

func (suite *RepositoriesTestSuite) TestOverlayFailsCloneReturnsError() {
	data := `
---
- url: invalid.
  version: abc1234
  dst: path/user.repo
`
	suite.r.UnmarshalYAML([]byte(data))
	anon := func() error {
		err := suite.r.Overlay()
		assert.Error(suite.T(), err)

		return err
	}

	git.MockRunCommandErrorsOn("git", anon)
}

func (suite *RepositoriesTestSuite) TestOverlayFailsCheckoutIndexReturnsError() {
	data := `
---
- url: https://example.com/user/repo.git
  version: abc1234
  dst: /invalid/directory
`
	suite.r.UnmarshalYAML([]byte(data))
	anon := func() error {
		err := suite.r.Overlay()
		assert.Error(suite.T(), err)

		return err
	}

	git.MockRunCommandErrorsOn("checkout-index", anon)
}

func (suite *RepositoriesTestSuite) TestOverlay() {
	data := `
---
- url: https://example.com/user/repo.git
  version: abc1234
  dst: path/user.repo
`
	suite.r.UnmarshalYAML([]byte(data))
	anon := func() error {
		err := suite.r.Overlay()
		assert.NoError(suite.T(), err)

		return err
	}

	dstDir, _ := git.FilePathAbs(suite.r.Items[0].Dst)
	got := git.MockRunCommand(anon)
	want := []string{
		fmt.Sprintf("git clone https://example.com/user/repo.git %s/https---example.com-user-repo.git-abc1234",
			repositories.GiltDir),
		fmt.Sprintf("git -C %s/https---example.com-user-repo.git-abc1234 reset --hard abc1234",
			repositories.GiltDir),
		fmt.Sprintf("git -C %s/https---example.com-user-repo.git-abc1234 checkout-index --force --all --prefix %s",
			repositories.GiltDir, (dstDir + string(os.PathSeparator))),
	}

	assert.Equal(suite.T(), want, got)
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestRepositoriesTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoriesTestSuite))
}
