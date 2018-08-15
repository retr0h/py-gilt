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

package repositories

import (
	"errors"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/retr0h/go-gilt/test/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/xeipuuv/gojsonschema"
)

type RepositoriesTestSuite struct {
	suite.Suite
	r Repositories
}

func (suite *RepositoriesTestSuite) SetupTest() {
	suite.r = Repositories{}
	GiltDir = testutil.CreateTempDirectory()
}

func (suite *RepositoriesTestSuite) TearDownTest() {
	testutil.RemoveTempDirectory(GiltDir)
}

func (suite *RepositoriesTestSuite) TestValidateSchemaHasErrorReturnsError() {
	originalJSONSchemaValidator := jsonSchemaValidator
	jsonSchemaValidator = func(gojsonschema.JSONLoader, gojsonschema.JSONLoader) (*gojsonschema.Result, error) {
		return nil, errors.New("Failed to load schema")
	}
	defer func() { jsonSchemaValidator = originalJSONSchemaValidator }()

	err := suite.r.validate([]byte(""))
	want := "Failed to load schema"

	assert.Equal(suite.T(), want, err.Error())
}

func (suite *RepositoriesTestSuite) TestValidateWithoutRootArrayReturnsError() {
	data := `
---
foo:
`
	jsonData, _ := yaml.YAMLToJSON([]byte(data))
	err := suite.r.validate([]byte(jsonData))
	want := errors.New("(root): Invalid type. Expected: array, given: object")

	assert.Equal(suite.T(), want, err)
}

func (suite *RepositoriesTestSuite) TestValidateRequiredTopLevelKeysReturnsError() {
	data := `
---
- missing:
  required:
  keys:
`
	jsonData, _ := yaml.YAMLToJSON([]byte(data))
	err := suite.r.validate([]byte(jsonData))

	messages := []string{
		"git: git is required",
		"version: version is required",
		"dstDir: dstDir is required",
	}

	for _, want := range messages {
		assert.Contains(suite.T(), err.Error(), want)
	}
}

func (suite *RepositoriesTestSuite) TestValidateRequiredSourcesKeysReturnsError() {
	data := `
---
- git: https://example.com/user/repo.git
  version: abc1234
  sources:
    - missing:
      required:
      keys:
`
	jsonData, _ := yaml.YAMLToJSON([]byte(data))
	err := suite.r.validate([]byte(jsonData))

	messages := []string{
		"src: src is required",
		"dstFile: dstFile is required",
	}

	for _, want := range messages {
		assert.Contains(suite.T(), err.Error(), want)
	}
}

func (suite *RepositoriesTestSuite) TestValidateNoAdditionalTopLevelKeysReturnsError() {
	data := `
---
- git: https://example.com/user/repo.git
  version: abc1234
  dstDir: path/user.repo
  extra:
`
	jsonData, _ := yaml.YAMLToJSON([]byte(data))
	err := suite.r.validate([]byte(jsonData))
	want := "extra: Additional property extra is not allowed"

	assert.Equal(suite.T(), want, err.Error())
}

func (suite *RepositoriesTestSuite) TestValidateNoAdditionalSourcesKeysReturnsError() {
	data := `
---
- git: https://example.com/user/repo.git
  version: abc1234
  sources:
    - src: foo
      dstFile: bar
      extra:
`
	jsonData, _ := yaml.YAMLToJSON([]byte(data))
	err := suite.r.validate([]byte(jsonData))
	want := "extra: Additional property extra is not allowed"

	assert.Equal(suite.T(), want, err.Error())
}

func (suite *RepositoriesTestSuite) TestValidateMutuallyExclusiveSourcesKeysReturnsError() {
	data := `
---
- git: https://example.com/user/repo.git
  version: abc1234
  sources:
    - src: foo
      dstFile: bar
      dstDir: bar
`
	jsonData, _ := yaml.YAMLToJSON([]byte(data))
	err := suite.r.validate([]byte(jsonData))
	want := "0.sources.0: Must validate one and only one schema (oneOf)"

	assert.Equal(suite.T(), want, err.Error())
}

func (suite *RepositoriesTestSuite) TestValidateWithoutValueReturnsError() {
	data := `
---
- git:
  version:
  dstDir:

- git:
  version:
  sources:
    - src:
      dstFile:

- git:
  version:
  sources:
    - src:
      dstDir:
`
	jsonData, _ := yaml.YAMLToJSON([]byte(data))
	err := suite.r.validate([]byte(jsonData))
	messages := []string{
		"0.git: Invalid type. Expected: string, given: null",
		"0.version: Invalid type. Expected: string, given: null",
		"0.dstDir: Invalid type. Expected: string, given: null",
		"1.git: Invalid type. Expected: string, given: null",
		"1.version: Invalid type. Expected: string, given: null",
		"1.sources.0.src: Invalid type. Expected: string, given: null",
		"1.sources.0.dstFile: Invalid type. Expected: string, given: null",
		"2.git: Invalid type. Expected: string, given: null",
		"2.version: Invalid type. Expected: string, given: null",
		"2.sources.0.src: Invalid type. Expected: string, given: null",
		"2.sources.0.dstDir: Invalid type. Expected: string, given: null",
	}

	for _, want := range messages {
		assert.Contains(suite.T(), err.Error(), want)
	}
}

func (suite *RepositoriesTestSuite) TestValidateMutuallyExclusiveReturnsError() {
	data := `
---
- git: https://example.com/user/repo.git
  version: abc1234
  dstDir: path/user.repo
  sources:
    - src: foo
      dstFile: bar
`
	jsonData, _ := yaml.YAMLToJSON([]byte(data))
	err := suite.r.validate([]byte(jsonData))
	want := "0: Must validate one and only one schema (oneOf)"

	assert.Equal(suite.T(), want, err.Error())
}

func (suite *RepositoriesTestSuite) TestValidate() {
	data := `
---
- git: https://example.com/user/repo.git
  version: abc1234
  dstDir: path/user.repo
`
	jsonData, _ := yaml.YAMLToJSON([]byte(data))
	err := suite.r.validate([]byte(jsonData))

	assert.NoError(suite.T(), err)
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestRepositoriesTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoriesTestSuite))
}
