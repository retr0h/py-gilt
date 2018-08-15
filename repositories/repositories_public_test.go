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
- git: https://example.com/user/repo.git
  version: abc1234
  dstDir: path/user.repo

- git: https://example.com/user/repo.git
  version: abc6789
  sources:
    - src: foo
      dstFile: bar
`
	err := suite.r.UnmarshalYAML([]byte(data))

	assert.NoError(suite.T(), err)

	firstItem := suite.r.Items[0]
	assert.Equal(suite.T(), "https://example.com/user/repo.git", firstItem.Git)
	assert.Equal(suite.T(), "abc1234", firstItem.Version)
	assert.Equal(suite.T(), "path/user.repo", firstItem.DstDir)
	assert.Empty(suite.T(), firstItem.Sources)

	secondItem := suite.r.Items[1]
	fmt.Println(secondItem)
	assert.Equal(suite.T(), "https://example.com/user/repo.git", secondItem.Git)
	assert.Equal(suite.T(), "abc6789", secondItem.Version)
	assert.Empty(suite.T(), secondItem.DstDir)
	assert.Equal(suite.T(), "foo", secondItem.Sources[0].Src)
	assert.Equal(suite.T(), "bar", secondItem.Sources[0].DstFile)
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

	firstItem := suite.r.Items[0]
	assert.NotNil(suite.T(), firstItem.Git)
	assert.NotNil(suite.T(), firstItem.Version)
	assert.NotNil(suite.T(), firstItem.DstDir)
	assert.Empty(suite.T(), firstItem.Sources)

	secondItem := suite.r.Items[1]
	assert.NotNil(suite.T(), secondItem.Git)
	assert.NotNil(suite.T(), secondItem.Version)
	assert.Empty(suite.T(), secondItem.DstDir)
	assert.NotNil(suite.T(), secondItem.Sources[0].Src)
	assert.NotNil(suite.T(), secondItem.Sources[0].DstFile)
	assert.NotNil(suite.T(), secondItem.Sources[1].Src)
	assert.NotNil(suite.T(), secondItem.Sources[1].DstFile)
	assert.NotNil(suite.T(), secondItem.Sources[2].Src)
	assert.NotNil(suite.T(), secondItem.Sources[2].DstFile)
}

func (suite *RepositoriesTestSuite) TestOverlayFailsCloneReturnsError() {
	data := `
---
- git: invalid.
  version: abc1234
  dstDir: path/user.repo
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
- git: https://example.com/user/repo.git
  version: abc1234
  dstDir: /invalid/directory
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
- git: https://example.com/user/repo1.git
  version: abc1234
  dstDir: path/user.repo
- git: https://example.com/user/repo2.git
  version: abc1234
  sources:
    - src: foo
      dstFile: bar
`
	suite.r.UnmarshalYAML([]byte(data))
	anon := func() error {
		err := suite.r.Overlay()
		assert.NoError(suite.T(), err)

		return err
	}

	dstDir, _ := git.FilePathAbs(suite.r.Items[0].DstDir)
	got := git.MockRunCommand(anon)
	want := []string{
		fmt.Sprintf("git clone https://example.com/user/repo1.git %s/https---example.com-user-repo1.git-abc1234",
			repositories.GiltDir),
		fmt.Sprintf("git -C %s/https---example.com-user-repo1.git-abc1234 reset --hard abc1234",
			repositories.GiltDir),
		fmt.Sprintf("git -C %s/https---example.com-user-repo1.git-abc1234 checkout-index --force --all --prefix %s",
			repositories.GiltDir, (dstDir + string(os.PathSeparator))),
		fmt.Sprintf("git clone https://example.com/user/repo2.git %s/https---example.com-user-repo2.git-abc1234",
			repositories.GiltDir),
		fmt.Sprintf("git -C %s/https---example.com-user-repo2.git-abc1234 reset --hard abc1234",
			repositories.GiltDir),
	}

	assert.Equal(suite.T(), want, got)
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestRepositoriesTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoriesTestSuite))
}
