package render

import (
	"fmt"
	"strings"

	"github.com/xcnt/go-asyncapi/internal/types"

	"github.com/samber/lo"

	"github.com/xcnt/go-asyncapi/internal/common"
	"github.com/xcnt/go-asyncapi/internal/utils"
	"github.com/dave/jennifer/jen"
)

// GoStruct defines the data required to generate a struct in Go.
type GoStruct struct {
	BaseType
	Fields []GoStructField
}

func (s GoStruct) RenderDefinition(ctx *common.RenderContext) []*jen.Statement {
	var res []*jen.Statement
	ctx.LogStartRender("GoStruct", s.Import, s.Name, "definition", s.DirectRendering())
	defer ctx.LogFinishRender()

	if s.Description != "" {
		res = append(res, jen.Comment(s.Name+" -- "+utils.ToLowerFirstLetter(s.Description)))
	}
	code := lo.FlatMap(s.Fields, func(item GoStructField, _ int) []*jen.Statement {
		return item.renderDefinition(ctx)
	})
	res = append(res, jen.Type().Id(s.Name).Struct(utils.ToCode(code)...))
	return res
}

func (s GoStruct) RenderUsage(ctx *common.RenderContext) []*jen.Statement {
	ctx.LogStartRender("GoStruct", s.Import, s.Name, "usage", s.DirectRendering())
	defer ctx.LogFinishRender()

	if s.DirectRendering() {
		if s.Import != "" && s.Import != ctx.CurrentPackage {
			return []*jen.Statement{jen.Qual(ctx.GeneratedModule(s.Import), s.Name)}
		}
		return []*jen.Statement{jen.Id(s.Name)}
	}

	code := lo.FlatMap(s.Fields, func(item GoStructField, _ int) []*jen.Statement {
		return item.renderDefinition(ctx)
	})

	return []*jen.Statement{jen.Struct(utils.ToCode(code)...)}
}

func (s GoStruct) NewFuncName() string {
	return "New" + s.Name
}

func (s GoStruct) ReceiverName() string {
	return strings.ToLower(string(s.Name[0]))
}

func (s GoStruct) MustGetField(name string) GoStructField {
	f, ok := lo.Find(s.Fields, func(item GoStructField) bool {
		return item.Name == name
	})
	if !ok {
		panic(fmt.Errorf("field %s.%s not found", s.Name, name))
	}
	return f
}

func (s GoStruct) IsStruct() bool {
	return true
}

// GoStructField defines the data required to generate a field in Go.
type GoStructField struct {
	Name           string
	MarshalName    string
	Description    string
	Type           common.GolangType
	TagsSource     *ListPromise[*Message]
	ExtraTags      types.OrderedMap[string, string] // Just append these tags as constant, overwrite other tags on overlap
	ExtraTagNames  []string                         // Append these tags and fill them the same value as others
	ExtraTagValues []string                         // Add these comma-separated values to all tags (excluding ExtraTags)
}

func (f GoStructField) renderDefinition(ctx *common.RenderContext) []*jen.Statement {
	var res []*jen.Statement
	ctx.LogStartRender("GoStructField", "", f.Name, "definition", false)
	defer ctx.LogFinishRender()

	if f.Description != "" {
		res = append(res, jen.Comment(f.Name+" -- "+utils.ToLowerFirstLetter(f.Description)))
	}

	stmt := jen.Id(f.Name)

	items := utils.ToCode(f.Type.RenderUsage(ctx))
	stmt = stmt.Add(items...)

	if f.TagsSource != nil {
		tagValues := append([]string{f.MarshalName}, f.ExtraTagValues...)
		tagNames := lo.Uniq(lo.FilterMap(f.TagsSource.Targets(), func(item *Message, _ int) (string, bool) {
			format := getFormatByContentType(item.ContentType)
			return format, format != ""
		}))
		tagNames = append(tagNames, f.ExtraTagNames...)

		tags := lo.SliceToMap(tagNames, func(item string) (string, string) {
			return item, strings.Join(tagValues, ",")
		})
		tags = lo.Assign(tags, lo.FromEntries(f.ExtraTags.Entries()))
		stmt = stmt.Tag(tags)
	}

	res = append(res, stmt)

	return res
}

