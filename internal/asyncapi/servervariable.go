package asyncapi

import (
	"github.com/xcnt/go-asyncapi/internal/common"
	"github.com/xcnt/go-asyncapi/internal/render"
	"github.com/xcnt/go-asyncapi/internal/utils"
)

type ServerVariable struct {
	Enum        []string `json:"enum" yaml:"enum"`
	Default     string   `json:"default" yaml:"default"`
	Description string   `json:"description" yaml:"description"`
	Examples    []string `json:"examples" yaml:"examples"`

	Ref string `json:"$ref" yaml:"$ref"`
}

func (sv ServerVariable) Compile(ctx *common.CompileContext) error {
	ctx.RegisterNameTop(ctx.Stack.Top().PathItem)
	obj, err := sv.build(ctx, ctx.Stack.Top().PathItem)
	if err != nil {
		return err
	}
	ctx.PutObject(obj)
	return nil
}

func (sv ServerVariable) build(ctx *common.CompileContext, serverVariableKey string) (common.Renderer, error) {
	if sv.Ref != "" {
		ctx.Logger.Trace("Ref", "$ref", sv.Ref)
		res := render.NewRendererPromise(sv.Ref, common.PromiseOriginUser)
		ctx.PutPromise(res)
		return res, nil
	}

	res := &render.ServerVariable{
		Name:        utils.ToGolangName(serverVariableKey, false),
		Enum:        sv.Enum,
		Default:     sv.Default,
		Description: sv.Description,
	}

	return res, nil
}
