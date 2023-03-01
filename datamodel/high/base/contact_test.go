// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package base

import (
	"fmt"
	lowmodel "github.com/pb33f/libopenapi/datamodel/low"
	lowbase "github.com/pb33f/libopenapi/datamodel/low/base"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestNewContact(t *testing.T) {

	var cNode yaml.Node

	yml := `name: pizza
url: https://pb33f.io
email: buckaroo@pb33f.io`

	_ = yaml.Unmarshal([]byte(yml), &cNode)

	// build low
	var lowContact lowbase.Contact
	_ = lowmodel.BuildModel(cNode.Content[0], &lowContact)

	// build high
	highContact := NewContact(&lowContact)

	assert.Equal(t, "pizza", highContact.Name)
	assert.Equal(t, "https://pb33f.io", highContact.URL)
	assert.Equal(t, "buckaroo@pb33f.io", highContact.Email)
	assert.Equal(t, 1, highContact.GoLow().Name.KeyNode.Line)

}

func ExampleNewContact() {

	// define a Contact using yaml (or JSON, it doesn't matter)
	yml := `name: Buckaroo
url: https://pb33f.io
email: buckaroo@pb33f.io`

	// unmarshal yaml into a *yaml.Node instance
	var cNode yaml.Node
	_ = yaml.Unmarshal([]byte(yml), &cNode)

	// build low
	var lowContact lowbase.Contact
	_ = lowmodel.BuildModel(cNode.Content[0], &lowContact)

	// build high
	highContact := NewContact(&lowContact)
	fmt.Print(highContact.Name)
	// Output: Buckaroo
}

func TestContact_MarshalYAML(t *testing.T) {

	highC := &Contact{Name: "dave", URL: "https://pb33f.io", Email: "dave@pb33f.io"}
	dat, _ := highC.Render()

	// unmarshal yaml into a *yaml.Node instance
	var cNode yaml.Node
	_ = yaml.Unmarshal(dat, &cNode)

	// build low
	var lowContact lowbase.Contact
	_ = lowmodel.BuildModel(cNode.Content[0], &lowContact)

	// build high
	highContact := NewContact(&lowContact)

	assert.Equal(t, "dave", highContact.Name)
	assert.Equal(t, "dave@pb33f.io", highContact.Email)
	assert.Equal(t, "https://pb33f.io", highContact.URL)

}
