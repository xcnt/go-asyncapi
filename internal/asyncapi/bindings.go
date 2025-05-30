package asyncapi

import (
	"encoding/json"
	"fmt"

	"github.com/xcnt/go-asyncapi/internal/common"
	"github.com/xcnt/go-asyncapi/internal/render"
	"github.com/xcnt/go-asyncapi/internal/types"
	"gopkg.in/yaml.v3"
)

type RawBindings struct {
	Ref            string                                                             `json:"$ref" yaml:"$ref"`
	ProtocolValues types.OrderedMap[string, types.Union2[json.RawMessage, yaml.Node]] `json:"-" yaml:"-"`
}

type Bindings struct {
	RawBindings
}

const (
	bindingsKindMessage = iota
	bindingsKindOperation
	bindingsKindChannel
	bindingsKindServer
)

func (b *Bindings) build(
	ctx *common.CompileContext,
	bindingsKind int,
	bindingsKey string,
) (common.Renderer, error) {
	if b.Ref != "" {
		ctx.Logger.Trace("Ref", "$ref", b.Ref)
		res := render.NewRendererPromise(b.Ref, common.PromiseOriginUser)
		ctx.PutPromise(res)
		return res, nil
	}

	res := render.Bindings{Name: bindingsKey}
	for _, e := range b.ProtocolValues.Entries() {
		ctx.Logger.Trace("Bindings", "proto", e.Key)
		builder, ok := ProtocolBuilders[e.Key]
		if !ok {
			ctx.Logger.Warn("Skip bindings protocol %q since it is not supported", "proto", e.Key)
			continue
		}

		var vals *render.GoValue
		var jsonVals types.OrderedMap[string, string]
		var err error

		switch bindingsKind {
		case bindingsKindMessage:
			ctx.Logger.Trace("Message bindings", "proto", builder.ProtocolName())
			vals, jsonVals, err = builder.BuildMessageBindings(ctx, e.Value)
		case bindingsKindOperation:
			ctx.Logger.Trace("Operation bindings", "proto", builder.ProtocolName())
			vals, jsonVals, err = builder.BuildOperationBindings(ctx, e.Value)
		case bindingsKindChannel:
			ctx.Logger.Trace("Channel bindings", "proto", builder.ProtocolName())
			vals, jsonVals, err = builder.BuildChannelBindings(ctx, e.Value)
		case bindingsKindServer:
			ctx.Logger.Trace("Server bindings", "proto", builder.ProtocolName())
			vals, jsonVals, err = builder.BuildServerBindings(ctx, e.Value)
		}
		if err != nil {
			return nil, types.CompileError{Err: fmt.Errorf("bindings build: %w", err), Path: ctx.PathStackRef(), Proto: e.Key}
		}
		if vals != nil {
			ctx.Logger.Trace("Have bindings values", "proto", builder.ProtocolName(), "value", vals)
			res.Values.Set(e.Key, vals)
		}
		if jsonVals.Len() > 0 {
			ctx.Logger.Trace("Have bindings jsonschema values", "proto", builder.ProtocolName(), "keys", jsonVals.Keys())
			res.JSONValues.Set(e.Key, jsonVals)
		}
	}

	return &res, nil
}

func (b *Bindings) UnmarshalJSON(value []byte) error {
	if err := json.Unmarshal(value, &b.RawBindings); err != nil {
		return err
	}
	if err := json.Unmarshal(value, &b.RawBindings.ProtocolValues); err != nil {
		return err
	}
	b.ProtocolValues.Delete("$ref")
	return nil
}

func (b *Bindings) UnmarshalYAML(value *yaml.Node) error {
	if err := value.Decode(&b.RawBindings); err != nil {
		return err
	}
	if err := value.Decode(&b.RawBindings.ProtocolValues); err != nil {
		return err
	}
	b.ProtocolValues.Delete("$ref")
	return nil
}

type MessageBindings struct {
	Bindings
}

func (b *MessageBindings) Compile(ctx *common.CompileContext) error {
	ctx.RegisterNameTop(ctx.Stack.Top().PathItem)
	obj, err := b.build(ctx, bindingsKindMessage, ctx.Stack.Top().PathItem)
	if err != nil {
		return err
	}
	ctx.PutObject(obj)
	return nil
}

type OperationBinding struct {
	Bindings
}

func (b *OperationBinding) Compile(ctx *common.CompileContext) error {
	ctx.RegisterNameTop(ctx.Stack.Top().PathItem)
	obj, err := b.build(ctx, bindingsKindOperation, ctx.Stack.Top().PathItem)
	if err != nil {
		return err
	}
	ctx.PutObject(obj)
	return nil
}

type ChannelBindings struct {
	Bindings
}

func (b *ChannelBindings) Compile(ctx *common.CompileContext) error {
	ctx.RegisterNameTop(ctx.Stack.Top().PathItem)
	obj, err := b.build(ctx, bindingsKindChannel, ctx.Stack.Top().PathItem)
	if err != nil {
		return err
	}
	ctx.PutObject(obj)
	return nil
}

type ServerBindings struct {
	Bindings
}

func (b *ServerBindings) Compile(ctx *common.CompileContext) error {
	ctx.RegisterNameTop(ctx.Stack.Top().PathItem)
	obj, err := b.build(ctx, bindingsKindServer, ctx.Stack.Top().PathItem)
	if err != nil {
		return err
	}
	ctx.PutObject(obj)
	return nil
}
