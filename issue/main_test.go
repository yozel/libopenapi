package main

import (
	"context"
	"os"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/datamodel/low/base"
	"github.com/pb33f/libopenapi/index"
)

func TestThroughComponents(t *testing.T) {
	model, err := buildModel("spec.yaml")
	require.NoError(t, err)

	sp := model.Components.Schemas.Value("ExternalRef")
	schema, err := sp.BuildSchema()
	require.NoError(t, err)
	require.NotNil(t, schema)
}

func TestThroughPaths(t *testing.T) {
	model, err := buildModel("spec.yaml")
	require.NoError(t, err)

	sp := model.Paths.PathItems.Value("/hello").Post.Responses.Codes.Value("200").Content.Value("application/json").Schema
	schema, err := sp.BuildSchema()
	require.NoError(t, err)
	require.NotNil(t, schema)
}

func TestThroughComponentsCtx(t *testing.T) {
	model, err := buildModel("spec.yaml")
	require.NoError(t, err)

	modelLow := model.GoLow()
	lowSchemaProxy := modelLow.Components.Value.Schemas.Value.First().Value().Value

	ctx := unsafeExtractPrivateCtx(lowSchemaProxy)
	val := ctx.Value(index.CurrentPathKey)
	require.NotNil(t, val)
	require.IsType(t, "", val)
	require.Contains(t, val, "spec-ext.yaml")
}

func TestThroughPathsCtx(t *testing.T) {
	model, err := buildModel("spec.yaml")
	require.NoError(t, err)

	modelLow := model.GoLow()
	lowSchemaProxy := modelLow.Paths.Value.PathItems.First().Value().Value.Post.Value.Responses.Value.Codes.First().Value().Value.Content.Value.First().Value().Value.Schema.Value

	ctx := unsafeExtractPrivateCtx(lowSchemaProxy)
	val := ctx.Value(index.CurrentPathKey)
	require.NotNil(t, val)
	require.IsType(t, "", val)
	require.Contains(t, val, "spec-ext.yaml")
}

func buildModel(p string) (*v3.Document, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	document, err := libopenapi.NewDocumentWithConfiguration(data, &datamodel.DocumentConfiguration{
		AllowFileReferences: true,
		BasePath:            ".",
	})
	if err != nil {
		return nil, err
	}

	modelv3, errs := document.BuildV3Model()
	for _, err := range errs {
		return nil, err
	}
	return &modelv3.Model, nil
}

func unsafeExtractPrivateCtx(lowSchemaProxy *base.SchemaProxy) context.Context {
	ptr := unsafe.Pointer(lowSchemaProxy)
	offset := uintptr(72) // offset := unsafe.Offsetof(base.SchemaProxy{}.ctx)
	ctx := (*context.Context)(unsafe.Pointer(uintptr(ptr) + offset))
	return *ctx
}
