// +build integration
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

package repository_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/retr0h/go-gilt/git"
	"github.com/retr0h/go-gilt/repositories"
	"github.com/retr0h/go-gilt/repository"
	"github.com/retr0h/go-gilt/test/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RepositoryIntegrationTestSuite struct {
	suite.Suite
	r  repository.Repository
	rr repositories.Repositories
	g  *git.Git
}

func (suite *RepositoryIntegrationTestSuite) SetupTest() {
	suite.rr = repositories.Repositories{}
	suite.g = git.NewGit(suite.rr.Debug)
}

func (suite *RepositoryIntegrationTestSuite) TearDownTest() {
	testutil.RemoveTempDirectory(suite.r.GiltDir)
}

func (suite *RepositoryIntegrationTestSuite) TestCopySourcesHasErrorWhenDstDirDoesNotExist() {
	data := `
- git: https://github.com/lorin/openstack-ansible-modules.git
  version: 2677cc3
  sources:
    - src: "*_manage"
      dstDir: invalid/path
`
	suite.rr.UnmarshalYAML([]byte(data))
	r := suite.rr.Items[0]
	r.GiltDir = testutil.CreateTempDirectory()
	suite.g.Clone(r)
	err := r.CopySources()

	assert.Error(suite.T(), err)
}

func (suite *RepositoryIntegrationTestSuite) TestCopySourcesHasErrorWhenFileCopyFails() {
	tempDir := testutil.CreateTempDirectory()
	dstDir := filepath.Join(tempDir, "library")
	data := fmt.Sprintf(`
- git: https://github.com/lorin/openstack-ansible-modules.git
  version: 2677cc3
  sources:
    - src: cinder_manage
      dstDir: %s
`, dstDir)
	suite.rr.UnmarshalYAML([]byte(data))
	r := suite.rr.Items[0]
	r.GiltDir = tempDir
	suite.g.Clone(r)
	os.Mkdir(dstDir, 0755)

	originalCopyFile := repository.CopyFile
	repository.CopyFile = func(src string, dst string) error {
		return errors.New("Failed to copy file")
	}
	defer func() { repository.CopyFile = originalCopyFile }()

	err := r.CopySources()

	assert.Error(suite.T(), err)
}

func (suite *RepositoryIntegrationTestSuite) TestCopySourcesHasErrorWhenDstFileDoesNotExist() {
	data := `
- git: https://github.com/lorin/openstack-ansible-modules.git
  version: 2677cc3
  sources:
    - src: cinder_manage
      dstFile: invalid/path
`
	suite.rr.UnmarshalYAML([]byte(data))
	r := suite.rr.Items[0]
	r.GiltDir = testutil.CreateTempDirectory()
	suite.g.Clone(r)
	err := r.CopySources()

	assert.Error(suite.T(), err)
}

func (suite *RepositoryIntegrationTestSuite) TestCopySourcesCopiesFile() {
	tempDir := testutil.CreateTempDirectory()
	dstFile := filepath.Join(tempDir, "cinder_manage")
	dstDir := filepath.Join(tempDir, "library")
	dstDirFile := filepath.Join(tempDir, "library", "glance_manage")
	data := fmt.Sprintf(`
- git: https://github.com/lorin/openstack-ansible-modules.git
  version: 2677cc3
  sources:
    - src: cinder_manage
      dstFile: %s
    - src: glance_manage
      dstDir: %s
`, dstFile, dstDir)
	suite.rr.UnmarshalYAML([]byte(data))
	r := suite.rr.Items[0]
	r.GiltDir = tempDir
	suite.g.Clone(r)
	os.Mkdir(dstDir, 0755)
	err := r.CopySources()

	assert.NoError(suite.T(), err)
	assert.FileExistsf(suite.T(), dstFile, "File does not exist")
	assert.FileExistsf(suite.T(), dstDirFile, "File does not exist")
}

func (suite *RepositoryIntegrationTestSuite) TestCopySourcesHasErrorWhenDirExistsAndDirCopyFails() {
	tempDir := testutil.CreateTempDirectory()
	dstDir := filepath.Join(tempDir, "tests")
	data := fmt.Sprintf(`
- git: https://github.com/lorin/openstack-ansible-modules.git
  version: 2677cc3
  sources:
    - src: tests
      dstDir: %s
`, dstDir)
	suite.rr.UnmarshalYAML([]byte(data))
	r := suite.rr.Items[0]
	r.GiltDir = tempDir
	suite.g.Clone(r)

	originalCopyDir := repository.CopyDir
	repository.CopyDir = func(src string, dst string) error {
		return errors.New("Failed to copy dir")
	}
	defer func() { repository.CopyDir = originalCopyDir }()

	err := r.CopySources()

	assert.Error(suite.T(), err)
}

func (suite *RepositoryIntegrationTestSuite) TestCopySourcesCopiesDir() {
	tempDir := testutil.CreateTempDirectory()
	dstDir := filepath.Join(tempDir, "tests")
	data := fmt.Sprintf(`
- git: https://github.com/lorin/openstack-ansible-modules.git
  version: 2677cc3
  sources:
    - src: tests
      dstDir: %s
`, dstDir)
	suite.rr.UnmarshalYAML([]byte(data))
	r := suite.rr.Items[0]
	r.GiltDir = tempDir
	os.Mkdir(dstDir, 0755) // execute the dstDir cleanup code prior to copy.
	suite.g.Clone(r)
	err := r.CopySources()

	assert.NoError(suite.T(), err)
	assert.DirExistsf(suite.T(), dstDir, "Dir does not exist")
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestRepositoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryIntegrationTestSuite))
}
