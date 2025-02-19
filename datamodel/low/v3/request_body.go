// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package v3

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/pb33f/libopenapi/utils"
	"gopkg.in/yaml.v3"
)

// RequestBody represents a low-level OpenAPI 3+ RequestBody object.
//   - https://spec.openapis.org/oas/v3.1.0#request-body-object
type RequestBody struct {
	Description low.NodeReference[string]
	Content     low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*MediaType]]]
	Required    low.NodeReference[bool]
	Extensions  *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	*low.Reference
}

// FindExtension attempts to locate an extension using the provided name.
func (rb *RequestBody) FindExtension(ext string) *low.ValueReference[*yaml.Node] {
	return low.FindItemInOrderedMap(ext, rb.Extensions)
}

// GetExtensions returns all RequestBody extensions and satisfies the low.HasExtensions interface.
func (rb *RequestBody) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return rb.Extensions
}

// FindContent attempts to find content/MediaType defined using a specified name.
func (rb *RequestBody) FindContent(cType string) *low.ValueReference[*MediaType] {
	return low.FindItemInOrderedMap[*MediaType](cType, rb.Content.Value)
}

// Build will extract extensions and MediaType objects from the node.
func (rb *RequestBody) Build(ctx context.Context, _, root *yaml.Node, idx *index.SpecIndex) error {
	root = utils.NodeAlias(root)
	utils.CheckForMergeNodes(root)
	rb.Reference = new(low.Reference)
	rb.Extensions = low.ExtractExtensions(root)

	// handle content, if set.
	con, cL, cN, cErr := low.ExtractMap[*MediaType](ctx, ContentLabel, root, idx)
	if cErr != nil {
		return cErr
	}
	if con != nil {
		rb.Content = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*MediaType]]]{
			Value:     con,
			KeyNode:   cL,
			ValueNode: cN,
		}
	}
	return nil
}

// Hash will return a consistent SHA256 Hash of the RequestBody object
func (rb *RequestBody) Hash() [32]byte {
	var f []string
	if rb.Description.Value != "" {
		f = append(f, rb.Description.Value)
	}
	if !rb.Required.IsEmpty() {
		f = append(f, fmt.Sprint(rb.Required.Value))
	}
	for pair := orderedmap.First(orderedmap.SortAlpha(rb.Content.Value)); pair != nil; pair = pair.Next() {
		f = append(f, low.GenerateHashString(pair.Value().Value))
	}
	f = append(f, low.HashExtensions(rb.Extensions)...)
	return sha256.Sum256([]byte(strings.Join(f, "|")))
}
