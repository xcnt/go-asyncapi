package compiler

import (
	"github.com/xcnt/go-asyncapi/internal/asyncapi"
)

type SpecKind string

const (
	SpecKindAsyncapi   SpecKind = "asyncapi"
	SpecKindJsonschema SpecKind = "jsonschema"
	SpecKindOpenapi    SpecKind = "openapi"
)

type specTypeTest struct {
	Asyncapi string `json:"asyncapi" yaml:"asyncapi"`
	Openapi  string `json:"openapi" yaml:"openapi"`
}

type anyDecoder interface {
	Decode(v any) error
}

func guessSpecKind(decoder anyDecoder) (SpecKind, compiledObject, error) {
	test := specTypeTest{}

	if err := decoder.Decode(&test); err != nil {
		return "", nil, err
	}
	switch {
	case test.Asyncapi != "":
		return SpecKindAsyncapi, &asyncapi.AsyncAPI{}, nil
	case test.Openapi != "":
		panic("openapi not implemented")
	}
	panic("jsonschema not implemented")
	// Assume that some data is jsonschema, TODO: maybe it's better to match more strict?
}
