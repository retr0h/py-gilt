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

func (suite *RepositoriesTestSuite) TestValidateWithoutStringReturnsError() {
	data := `
---
- url:
  version:
  dst:
`
	jsonData, _ := yaml.YAMLToJSON([]byte(data))
	err := suite.r.validate([]byte(jsonData))

	assert.Error(suite.T(), err)

	messages := []string{
		"0.dst: Invalid type. Expected: string, given: null",
		"0.url: Invalid type. Expected: string, given: null",
		"0.version: Invalid type. Expected: string, given: null",
	}

	for _, want := range messages {
		assert.Contains(suite.T(), err.Error(), want)
	}
}

func (suite *RepositoriesTestSuite) TestValidate() {
	data := `
---
- url: https://example.com/user/repo.git
  version: abc1234
  dst: path/user.repo
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
