package resolver

import (
	"github.com/pb33f/libopenapi/index"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"testing"
)

func TestNewResolver(t *testing.T) {
	assert.Nil(t, NewResolver(nil))
}

func Benchmark_ResolveDocumentStripe(b *testing.B) {
	stripe, _ := ioutil.ReadFile("../test_specs/stripe.yaml")
	for n := 0; n < b.N; n++ {
		var rootNode yaml.Node
		yaml.Unmarshal(stripe, &rootNode)
		index := index.NewSpecIndex(&rootNode)
		resolver := NewResolver(index)
		resolver.Resolve()
	}
}

func TestResolver_ResolveComponents_CircularSpec(t *testing.T) {

	circular, _ := ioutil.ReadFile("../test_specs/circular-tests.yaml")
	var rootNode yaml.Node
	yaml.Unmarshal(circular, &rootNode)

	index := index.NewSpecIndex(&rootNode)

	resolver := NewResolver(index)
	assert.NotNil(t, resolver)

	circ := resolver.Resolve()
	assert.Len(t, circ, 3)

	_, err := yaml.Marshal(resolver.resolvedRoot)
	assert.NoError(t, err)
}

func TestResolver_CheckForCircularReferences(t *testing.T) {

	circular, _ := ioutil.ReadFile("../test_specs/circular-tests.yaml")
	var rootNode yaml.Node
	yaml.Unmarshal(circular, &rootNode)

	index := index.NewSpecIndex(&rootNode)

	resolver := NewResolver(index)
	assert.NotNil(t, resolver)

	circ := resolver.CheckForCircularReferences()
	assert.Len(t, circ, 3)
	assert.Len(t, resolver.GetResolvingErrors(), 3)
	assert.Len(t, resolver.GetCircularErrors(), 3)

	_, err := yaml.Marshal(resolver.resolvedRoot)
	assert.NoError(t, err)
}

func TestResolver_DeepJourney(t *testing.T) {

	var journey []*index.Reference
	for f := 0; f < 200; f++ {
		journey = append(journey, nil)
	}
	index := index.NewSpecIndex(nil)
	resolver := NewResolver(index)
	assert.Nil(t, resolver.extractRelatives(nil, nil, journey, false))

}

func TestResolver_ResolveComponents_Stripe(t *testing.T) {

	stripe, _ := ioutil.ReadFile("../test_specs/stripe.yaml")
	var rootNode yaml.Node
	yaml.Unmarshal(stripe, &rootNode)

	index := index.NewSpecIndex(&rootNode)

	resolver := NewResolver(index)
	assert.NotNil(t, resolver)

	circ := resolver.Resolve()
	assert.Len(t, circ, 21)

	assert.Len(t, resolver.GetNonPolymorphicCircularErrors(), 2)
	assert.Len(t, resolver.GetPolymorphicCircularErrors(), 19)

}

func TestResolver_ResolveComponents_BurgerShop(t *testing.T) {

	mixedref, _ := ioutil.ReadFile("../test_specs/burgershop.openapi.yaml")
	var rootNode yaml.Node
	yaml.Unmarshal(mixedref, &rootNode)

	index := index.NewSpecIndex(&rootNode)

	resolver := NewResolver(index)
	assert.NotNil(t, resolver)

	circ := resolver.Resolve()
	assert.Len(t, circ, 0)

}

func TestResolver_ResolveComponents_PolyNonCircRef(t *testing.T) {

	yml := `paths:
  /hey:
    get:
      responses:
        "200":
          $ref: '#/components/schemas/crackers'
components:
  schemas:
    cheese:
      description: cheese
      anyOf:
        items:
          $ref: '#/components/schemas/crackers' 
    crackers:
      description: crackers
      allOf:
       - $ref: '#/components/schemas/tea'
    tea:
      description: tea`

	var rootNode yaml.Node
	yaml.Unmarshal([]byte(yml), &rootNode)

	index := index.NewSpecIndex(&rootNode)

	resolver := NewResolver(index)
	assert.NotNil(t, resolver)

	circ := resolver.CheckForCircularReferences()
	assert.Len(t, circ, 0)

}

func TestResolver_ResolveComponents_MixedRef(t *testing.T) {

	mixedref, _ := ioutil.ReadFile("../test_specs/mixedref-burgershop.openapi.yaml")
	var rootNode yaml.Node
	yaml.Unmarshal(mixedref, &rootNode)

	index := index.NewSpecIndex(&rootNode)

	resolver := NewResolver(index)
	assert.NotNil(t, resolver)

	circ := resolver.Resolve()
	assert.Len(t, circ, 2)

}

func TestResolver_ResolveComponents_k8s(t *testing.T) {

	k8s, _ := ioutil.ReadFile("../test_specs/k8s.json")
	var rootNode yaml.Node
	yaml.Unmarshal(k8s, &rootNode)

	index := index.NewSpecIndex(&rootNode)

	resolver := NewResolver(index)
	assert.NotNil(t, resolver)

	circ := resolver.Resolve()
	assert.Len(t, circ, 1)
}
