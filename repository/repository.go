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

package repository

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Repository containing the repository's details.
type Repository struct {
	URL     string `yaml:"url"`
	Version string `yaml:"version"`
	Dst     string `yaml:"dst"`
	GiltDir string // GiltDir option set from CLI.
}

// // NewRepository factory to create a new Repository instance.
// func NewRepository(debug bool) *Git {
//     return &Repository{}
// }

// GetCloneDir returns the path to the Repository's clone directory.
func (r *Repository) GetCloneDir() string {
	return filepath.Join(r.GiltDir, r.getCloneHash())
}

func (r *Repository) getCloneHash() string {
	replacer := strings.NewReplacer(
		"/", "-",
		":", "-",
	)
	replacedGitURL := replacer.Replace(r.URL)

	return fmt.Sprintf("%s-%s", replacedGitURL, r.Version)
}
