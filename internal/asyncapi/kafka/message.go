package kafka

import (
	"encoding/json"

	"github.com/xcnt/go-asyncapi/internal/common"
	"github.com/xcnt/go-asyncapi/internal/render"
	"github.com/xcnt/go-asyncapi/internal/types"
	"gopkg.in/yaml.v3"
)

type messageBindings struct {
	Key                     any    `json:"key" yaml:"key"` // jsonschema object
	SchemaIDLocation        string `json:"schemaIdLocation" yaml:"schemaIdLocation"`
	SchemaIDPayloadEncoding string `json:"schemaIdPayloadEncoding" yaml:"schemaIdPayloadEncoding"`
	SchemaLookupStrategy    string `json:"schemaLookupStrategy" yaml:"schemaLookupStrategy"`
}

func (pb ProtoBuilder) BuildMessageBindings(ctx *common.CompileContext, rawData types.Union2[json.RawMessage, yaml.Node]) (vals *render.GoValue, jsonVals types.OrderedMap[string, string], err error) {
	var bindings messageBindings
	if err = types.UnmarshalRawsUnion2(rawData, &bindings); err != nil {
		err = types.CompileError{Err: err, Path: ctx.PathStackRef(), Proto: pb.ProtoName}
		return
	}

	vals = render.ConstructGoValue(
		bindings, []string{"Key"}, &render.GoSimple{Name: "MessageBindings", Import: ctx.RuntimeModule(pb.ProtoName)},
	)
	if bindings.Key != nil {
		v, err2 := json.Marshal(bindings.Key)
		if err2 != nil {
			err = types.CompileError{Err: err2, Path: ctx.PathStackRef(), Proto: pb.ProtoName}
			return
		}
		jsonVals.Set("Key", string(v))
	}

	return
}
