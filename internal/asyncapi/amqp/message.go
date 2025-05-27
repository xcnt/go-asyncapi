package amqp

import (
	"encoding/json"

	"github.com/xcnt/go-asyncapi/internal/common"
	"github.com/xcnt/go-asyncapi/internal/render"
	"github.com/xcnt/go-asyncapi/internal/types"
	"gopkg.in/yaml.v3"
)

type messageBindings struct {
	ContentEncoding string `json:"contentEncoding" yaml:"contentEncoding"`
	MessageType     string `json:"messageType" yaml:"messageType"`
}

func (pb ProtoBuilder) BuildMessageBindings(ctx *common.CompileContext, rawData types.Union2[json.RawMessage, yaml.Node]) (vals *render.GoValue, jsonVals types.OrderedMap[string, string], err error) {
	var bindings messageBindings
	if err = types.UnmarshalRawsUnion2(rawData, &bindings); err != nil {
		err = types.CompileError{Err: err, Path: ctx.PathStackRef(), Proto: pb.ProtoName}
		return
	}

	vals = render.ConstructGoValue(
		bindings, nil, &render.GoSimple{Name: "MessageBindings", Import: ctx.RuntimeModule(pb.ProtoName)},
	)
	return
}
