// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package base

import (
	"sort"

	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/datamodel/low/base"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/pb33f/libopenapi/utils"
	"gopkg.in/yaml.v3"
)

// SecurityRequirement is a high-level representation of a Swagger / OpenAPI 3 SecurityRequirement object.
//
// SecurityRequirement lists the required security schemes to execute this operation. The object can have multiple
// security schemes declared in it which are all required (that is, there is a logical AND between the schemes).
//
// The name used for each property MUST correspond to a security scheme declared in the Security Definitions
//   - https://swagger.io/specification/v2/#securityDefinitionsObject
type SecurityRequirement struct {
	Requirements *orderedmap.Map[string, []string] `json:"-" yaml:"-"`
	low          *base.SecurityRequirement
}

// NewSecurityRequirement creates a new high-level SecurityRequirement from a low-level one.
func NewSecurityRequirement(req *base.SecurityRequirement) *SecurityRequirement {
	r := new(SecurityRequirement)
	r.low = req
	values := orderedmap.New[string, []string]()
	// to keep things fast, avoiding copying anything - makes it a little hard to read.
	for pair := orderedmap.First(req.Requirements.Value); pair != nil; pair = pair.Next() {
		var vals []string
		for valK := range pair.Value().Value {
			vals = append(vals, pair.Value().Value[valK].Value)
		}
		values.Set(pair.Key().Value, vals)
	}
	r.Requirements = values
	return r
}

// GoLow returns the low-level SecurityRequirement used to create the high-level one.
func (s *SecurityRequirement) GoLow() *base.SecurityRequirement {
	return s.low
}

// GoLowUntyped will return the low-level Discriminator instance that was used to create the high-level one, with no type
func (s *SecurityRequirement) GoLowUntyped() any {
	return s.low
}

// Render will return a YAML representation of the SecurityRequirement object as a byte slice.
func (s *SecurityRequirement) Render() ([]byte, error) {
	return yaml.Marshal(s)
}

// MarshalYAML will create a ready to render YAML representation of the SecurityRequirement object.
func (s *SecurityRequirement) MarshalYAML() (interface{}, error) {
	type req struct {
		line   int
		key    string
		val    []string
		lowKey *low.KeyReference[string]
		lowVal *low.ValueReference[[]low.ValueReference[string]]
	}

	m := utils.CreateEmptyMapNode()
	keys := make([]*req, orderedmap.Len(s.Requirements))

	i := 0

	for pair := orderedmap.First(s.Requirements); pair != nil; pair = pair.Next() {
		keys[i] = &req{key: pair.Key(), val: pair.Value()}
		i++
	}
	i = 0
	if s.low != nil {
		for o := range keys {
			kv := keys[o].key
			for pair := orderedmap.First(s.low.Requirements.Value); pair != nil; pair = pair.Next() {
				if pair.Key().Value == kv {
					keys[o].line = pair.Key().KeyNode.Line
					keys[o].lowKey = pair.KeyPtr()
					keys[o].lowVal = pair.ValuePtr()
				}
				i++
			}
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].line < keys[j].line
	})

	for k := range keys {
		l := utils.CreateStringNode(keys[k].key)
		l.Line = keys[k].line

		// for each key, extract all the values and order them.
		type req struct {
			line int
			val  string
		}

		reqs := make([]*req, len(keys[k].val))
		for t := range keys[k].val {
			reqs[t] = &req{val: keys[k].val[t], line: 9999 + t}
			if keys[k].lowVal != nil {
				for range keys[k].lowVal.Value[t].Value {
					fh := keys[k].val[t]
					df := keys[k].lowVal.Value[t].Value
					if fh == df {
						reqs[t].line = keys[k].lowVal.Value[t].ValueNode.Line
						break
					}
				}
			}
		}
		sort.Slice(reqs, func(i, j int) bool {
			return reqs[i].line < reqs[j].line
		})
		sn := utils.CreateEmptySequenceNode()
		sn.Line = keys[k].line + 1
		for z := range reqs {
			n := utils.CreateStringNode(reqs[z].val)
			n.Line = reqs[z].line + 1
			sn.Content = append(sn.Content, n)
		}

		m.Content = append(m.Content, l, sn)
	}
	return m, nil
}
