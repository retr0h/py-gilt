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
	"os"
	"path/filepath"

	"github.com/logrusorgru/aurora"
	"github.com/retr0h/go-gilt/repository"
	"github.com/retr0h/go-gilt/util"
)

var (
	// RunCommand is mocked for tests.
	RunCommand = util.RunCmd
	// FilePathAbs is mocked for tests.
	FilePathAbs = filepath.Abs
)

// Git struct for adding methods.
type Git struct {
	Debug bool // Debug option set from CLI with debug state.
}

// NewGit factory to create a new Git instance.
func NewGit(debug bool) *Git {
	return &Git{
		Debug: debug,
	}
}

// Clone clone Repository.URL to Repository.getCloneDir, and hard checkout
// to Repository.Version.
func (g *Git) Clone(repository repository.Repository) error {
	cloneDir := repository.GetCloneDir()

	msg := fmt.Sprintf("[%s@%s]:", aurora.Magenta(repository.URL), aurora.Magenta(repository.Version))
	fmt.Println(msg)

	msg = fmt.Sprintf("%-2s - %s '%s'", "", aurora.Cyan("Cloning to"), aurora.Cyan(cloneDir))
	fmt.Println(msg)

	if _, err := os.Stat(cloneDir); os.IsNotExist(err) {
		if err := g.clone(repository); err != nil {
			return err
		}

		if err := g.reset(repository); err != nil {
			return err
		}
	} else {
		msg := fmt.Sprintf("%-4s * %s", "", aurora.Brown("Clone already exists"))
		fmt.Println(msg)
	}

	return nil
}

func (g *Git) clone(repository repository.Repository) error {
	cloneDir := repository.GetCloneDir()
	err := RunCommand(g.Debug, "git", "clone", repository.URL, cloneDir)

	return err
}

func (g *Git) reset(repository repository.Repository) error {
	cloneDir := repository.GetCloneDir()
	err := RunCommand(g.Debug, "git", "-C", cloneDir, "reset", "--hard", repository.Version)

	return err
}

// CheckoutIndex checkout Repository.Git to Repository.Dst.
func (g *Git) CheckoutIndex(repository repository.Repository) error {
	cloneDir := repository.GetCloneDir()
	dstDir, err := FilePathAbs(repository.Dst)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("%-2s - %s '%s'", "", aurora.Cyan("Extracting to"), aurora.Cyan(dstDir))
	fmt.Println(msg)

	cmdArgs := []string{
		"-C",
		cloneDir,
		"checkout-index",
		"--force",
		"--all",
		"--prefix",
		// Trailing separator needed by git checkout-index.
		dstDir + string(os.PathSeparator),
	}
	if err := RunCommand(g.Debug, "git", cmdArgs...); err != nil {
		return err
	}

	return nil
}
