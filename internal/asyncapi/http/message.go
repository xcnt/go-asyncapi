package http

import (
	"encoding/json"

	"github.com/xcnt/go-asyncapi/internal/common"
	"github.com/xcnt/go-asyncapi/internal/render"
	"github.com/xcnt/go-asyncapi/internal/types"
	"gopkg.in/yaml.v3"
)

type messageBindings struct {
	Headers any `json:"headers" yaml:"headers"` // jsonschema object
}

func (pb ProtoBuilder) BuildMessageBindings(ctx *common.CompileContext, rawData types.Union2[json.RawMessage, yaml.Node]) (vals *render.GoValue, jsonVals types.OrderedMap[string, string], err error) {
	var bindings messageBindings
	if err = types.UnmarshalRawsUnion2(rawData, &bindings); err != nil {
		err = types.CompileError{Err: err, Path: ctx.PathStackRef(), Proto: pb.ProtoName}
		return
	}

	if bindings.Headers != nil {
		v, err2 := json.Marshal(bindings.Headers)
		if err2 != nil {
			err = types.CompileError{Err: err2, Path: ctx.PathStackRef(), Proto: pb.ProtoName}
			return
		}
		jsonVals.Set("Headers", string(v))
	}

	return
}
