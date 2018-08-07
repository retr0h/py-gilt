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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/retr0h/go-gilt/git"
	"github.com/retr0h/go-gilt/repository"
	"github.com/xeipuuv/gojsonschema"
)

var (
	// GiltDir option set from CLI.  // Should move to r.New
	GiltDir string
	// jsonSchemaValidator is mocked for tests.
	jsonSchemaValidator = gojsonschema.Validate
)

// Repositories is an object which implements the business logic interface.
type Repositories struct {
	Debug    bool                    // Debug option set from CLI with debug state.
	Filename string                  // Filename option set from CLI with path to gilt.yml.
	Items    []repository.Repository // Items slice containing Repository items.
}

// UnmarshalYAML decodes the first YAML document found within the data byte
// slice, passes the string through a generic YAML-to-JSON converter, performs
// validation, provides the resulting JSON to Unmarshaler, and assigns the
// decoded values to the Repositories struct.
func (r *Repositories) UnmarshalYAML(data []byte) error {
	jsonData, err := yaml.YAMLToJSON(data)
	if err != nil {
		return err
	}

	// Validate the jsonData against the schema.
	if err = r.validate(jsonData); err != nil {
		return err
	}

	// Unmarshal the jsonData to the Repositories struct.
	err = json.Unmarshal(jsonData, &r.Items)
	return err
}

// UnmarshalYAMLFile reads the file named by Filename and passes the source
// data byte slice to UnmarshalYAML for decoding.
func (r *Repositories) UnmarshalYAMLFile() error {
	// Open and Read the provided filename.
	source, err := ioutil.ReadFile(r.Filename)
	if err != nil {
		return err
	}

	// Unmarshal the file contents.
	err = r.UnmarshalYAML([]byte(source))
	return err
}

// Validate the the data byte slice against the JSON schema.
func (r *Repositories) validate(data []byte) error {
	schemaLoader := gojsonschema.NewStringLoader(configSchema)
	documentLoader := gojsonschema.NewBytesLoader(data)

	// Validate the document against the schema.
	result, err := jsonSchemaValidator(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	// Build schema validation failures.
	if !result.Valid() {
		var errstrings []string
		for _, desc := range result.Errors() {
			err := fmt.Errorf("%s", desc)
			errstrings = append(errstrings, err.Error())
		}

		return errors.New(strings.Join(errstrings, "\n"))
	}

	return nil
}

// Overlay clone and extract the Repository items.
func (r *Repositories) Overlay() error {
	g := git.NewGit(r.Debug)

	for _, repository := range r.Items {
		repository.GiltDir = GiltDir
		err := g.Clone(repository)
		if err != nil {
			return err
		}

		// Checkout into repository.Dst.
		if err := g.CheckoutIndex(repository); err != nil {
			return err
		}
	}

	return nil
}
