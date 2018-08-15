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

package repositories_test

import (
	"fmt"

	"github.com/retr0h/go-gilt/test/testutil"
	"github.com/stretchr/testify/assert"
)

func (suite *RepositoriesTestSuite) TestOverlayRemovesSrcDirPriorToCheckoutIndex() {
	tempDir := testutil.CreateTempDirectory()
	data := fmt.Sprintf(`
---
- git: https://github.com/retr0h/ansible-etcd.git
  version: 77a95b7
  dstDir: %s/retr0h.ansible-etcd
`, tempDir)
	suite.r.UnmarshalYAML([]byte(data))
	suite.r.Overlay()
	err := suite.r.Overlay()

	assert.NoError(suite.T(), err)
}

func (suite *RepositoriesTestSuite) TestOverlayFailsCopySourcesReturnsError() {
	data := `
---
- git: https://github.com/lorin/openstack-ansible-modules.git
  version: 2677cc3
  sources:
    - src: "*_manage"
      dstDir: /super/invalid/path/to/write/to
`
	suite.r.UnmarshalYAML([]byte(data))
	err := suite.r.Overlay()

	assert.Error(suite.T(), err)
}
