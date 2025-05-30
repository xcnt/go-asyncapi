package render

import (
	"github.com/samber/lo"

	"github.com/xcnt/go-asyncapi/internal/common"
	"github.com/xcnt/go-asyncapi/internal/utils"
	"github.com/dave/jennifer/jen"
)

type GoSimple struct {
	Name            string            // type name
	IsIface         bool              // true if type is interface, which means it cannot be rendered as pointer
	Import          string            // optional generated package name or module to import a type from
	TypeParamValues []common.Renderer // optional type parameter types to be filled in definition and usage
}

func (p GoSimple) DirectRendering() bool {
	return false
}

func (p GoSimple) RenderDefinition(ctx *common.RenderContext) []*jen.Statement {
	ctx.LogStartRender("GoSimple", p.Import, p.Name, "definition", p.DirectRendering())
	defer ctx.LogFinishRender()

	stmt := jen.Id(p.Name)
	if len(p.TypeParamValues) > 0 {
		typeParams := lo.FlatMap(p.TypeParamValues, func(item common.Renderer, _ int) []jen.Code {
			return utils.ToCode(item.RenderUsage(ctx))
		})
		stmt = stmt.Types(typeParams...)
	}
	return []*jen.Statement{stmt}
}

func (p GoSimple) RenderUsage(ctx *common.RenderContext) []*jen.Statement {
	ctx.LogStartRender("GoSimple", p.Import, p.Name, "usage", p.DirectRendering())
	defer ctx.LogFinishRender()

	stmt := &jen.Statement{}
	switch {
	case p.Import != "" && p.Import != ctx.CurrentPackage:
		stmt = stmt.Qual(p.Import, p.Name)
	default:
		stmt = stmt.Id(p.Name)
	}

	if len(p.TypeParamValues) > 0 {
		typeParams := lo.FlatMap(p.TypeParamValues, func(item common.Renderer, _ int) []jen.Code {
			return utils.ToCode(item.RenderUsage(ctx))
		})
		stmt = stmt.Types(typeParams...)
	}

	return []*jen.Statement{stmt}
}

func (p GoSimple) TypeName() string {
	return p.Name
}

func (p GoSimple) ID() string {
	return p.Name
}

func (p GoSimple) String() string {
	if p.Import != "" {
		return "GoSimple /" + p.Import + "." + p.Name
	}
	return "GoSimple " + p.Name
}

