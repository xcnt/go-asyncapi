package render

import (
	"github.com/xcnt/go-asyncapi/internal/common"
	"github.com/xcnt/go-asyncapi/internal/utils"
	"github.com/dave/jennifer/jen"
)

type GoPointer struct {
	Type         common.GolangType
	DirectRender bool
}

func (p GoPointer) DirectRendering() bool {
	return p.DirectRender
}

func (p GoPointer) RenderDefinition(ctx *common.RenderContext) []*jen.Statement {
	ctx.LogStartRender("GoPointer", "", "", "definition", p.DirectRendering())
	defer ctx.LogFinishRender()

	return p.Type.RenderDefinition(ctx)
}

func (p GoPointer) RenderUsage(ctx *common.RenderContext) []*jen.Statement {
	ctx.LogStartRender("GoPointer", "", "", "usage", p.DirectRendering())
	defer ctx.LogFinishRender()

	isPtr := true
	switch v := p.Type.(type) {
	case *GoInterface: // Prevent pointer to interface
		isPtr = false
	case golangPointerType:
		isPtr = !v.IsPointer() // Prevent appearing pointer to pointer
	case *GoSimple:
		isPtr = !v.IsIface
	}
	if isPtr {
		return []*jen.Statement{jen.Op("*").Add(utils.ToCode(p.Type.RenderUsage(ctx))...)}
	}
	return p.Type.RenderUsage(ctx)
}

func (p GoPointer) TypeName() string {
	return p.Type.TypeName()
}

func (p GoPointer) ID() string {
	return p.Type.ID()
}

func (p GoPointer) String() string {
	return "GoPointer -> " + p.Type.String()
}

func (p GoPointer) WrappedGolangType() (common.GolangType, bool) {
	return p.Type, p.Type != nil
}

func (p GoPointer) IsPointer() bool {
	return true
}
