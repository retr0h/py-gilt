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

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version   string
	buildHash string
	buildDate string
	debug     bool
	giltDir   string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-gilt",
	Short: "A GIT layering command line tool",
	Long: `

     ,o888888o.     8 8888 8 8888         8888888 8888888888
    8888     '88.   8 8888 8 8888               8 8888
 ,8 8888       '8.  8 8888 8 8888               8 8888
 88 8888            8 8888 8 8888               8 8888
 88 8888            8 8888 8 8888               8 8888
 88 8888            8 8888 8 8888               8 8888
 88 8888   8888888  8 8888 8 8888               8 8888
 '8 8888       .8'  8 8888 8 8888               8 8888
    8888     ,88'   8 8888 8 8888               8 8888
     '8888888P'     8 8888 8 888888888888       8 8888

A GIT layering command line tool.

https://github.com/retr0h/go-gilt
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string, bh string, bd string) {
	version = v
	buildHash = bh
	buildDate = bd

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable or disable debug mode")
	rootCmd.PersistentFlags().StringVarP(&giltDir, "giltdir", "c", "~/.gilt/clone", "Path to Gilt's clone dir")
}
