package render

import (
	"github.com/samber/lo"

	j "github.com/dave/jennifer/jen"
	"github.com/xcnt/go-asyncapi/internal/common"
	"github.com/xcnt/go-asyncapi/internal/utils"
)

const encodingPackageName = "encoding"

var encodingEncoders = map[string]j.Code{
	"application/json":      j.Op(`func(w io.Writer) Encoder`).Block(j.Return(j.Qual("encoding/json", "NewEncoder").Call(j.Id("w")))),
	"application/yaml":      j.Op(`func(w io.Writer) Encoder`).Block(j.Return(j.Qual("gopkg.in/yaml.v3", "NewEncoder").Call(j.Id("w")))),
	"application/x-msgpack": j.Op(`func(w io.Writer) Encoder`).Block(j.Return(j.Qual("github.com/vmihailenco/msgpack/v5", "NewEncoder").Call(j.Id("w")))),
	"text/plain":            j.Op(`func(w io.Writer) Encoder`).Block(j.Return(j.Qual("github.com/xcnt/client-go/pkg/text", "NewEncoder").Call(j.Id("w")))),
	// TODO: add other encoders: protobuf, avro, etc.
}

var encodingDecoders = map[string]j.Code{
	"application/json":      j.Op(`func(r io.Reader) Decoder`).Block(j.Return(j.Qual("encoding/json", "NewDecoder").Call(j.Id("r")))),
	"application/yaml":      j.Op(`func(r io.Reader) Decoder`).Block(j.Return(j.Qual("gopkg.in/yaml.v3", "NewDecoder").Call(j.Id("r")))),
	"application/x-msgpack": j.Op(`func(r io.Reader) Decoder`).Block(j.Return(j.Qual("github.com/vmihailenco/msgpack/v5", "NewDecoder").Call(j.Id("r")))),
	"text/plain":            j.Op(`func(r io.Reader) Decoder`).Block(j.Return(j.Qual("github.com/xcnt/client-go/pkg/text", "NewDecoder").Call(j.Id("r")))),
	// TODO: add other decoders: protobuf, avro, etc.
}

type EncodingEncode struct {
	AllMessages        *ListPromise[*Message]
	DefaultContentType string
}

func (e EncodingEncode) DirectRendering() bool {
	return true
}

func (e EncodingEncode) RenderDefinition(ctx *common.RenderContext) []*j.Statement {
	ctx.LogStartRender("EncodingEncode", "", "", "definition", e.DirectRendering())
	defer ctx.LogFinishRender()

	contentTypes := lo.Uniq(lo.FilterMap(e.AllMessages.Targets(), func(item *Message, _ int) (string, bool) {
		return item.ContentType, item.ContentType != ""
	}))
	return []*j.Statement{
		j.Op(`
			type Encoder interface {
				Encode(v any) error
			}`),

		j.Add(utils.QualSprintf(`var Encoders = map[string]func(w %Q(io,Writer)) Encoder`)).Values(j.DictFunc(func(d j.Dict) {
			for _, ct := range contentTypes {
				format := getFormatByContentType(ct)
				if format != "" {
					if v, ok := encodingEncoders[format]; ok {
						d[j.Lit(ct)] = v
					}
				} else {
					d[j.Lit(ct)] = j.Add(utils.QualSprintf(`func(_ %Q(io,Writer)) Encoder { panic("No encoder is set for content type %s") }`, ct))
				}
			}
		})),

		j.Add(utils.QualSprintf(`
			func NewEncoder(contentType string, w %Q(io,Writer)) Encoder {
				if v, ok := Encoders[contentType]; ok {
					return v(w)
				}
				panic("Unknown content type " + contentType)
			}`)),
	}
}

func (e EncodingEncode) RenderUsage(_ *common.RenderContext) []*j.Statement {
	panic("not implemented")
}

func (e EncodingEncode) ID() string {
	return "Encode"
}

func (e EncodingEncode) String() string {
	return "EncodingEncode"
}

type EncodingDecode struct {
	AllMessages        *ListPromise[*Message]
	DefaultContentType string
}

func (e EncodingDecode) DirectRendering() bool {
	return true
}

func (e EncodingDecode) RenderDefinition(ctx *common.RenderContext) []*j.Statement {
	ctx.LogStartRender("EncodingDecode", "", "", "definition", e.DirectRendering())
	defer ctx.LogFinishRender()

	contentTypes := lo.Uniq(lo.FilterMap(e.AllMessages.Targets(), func(item *Message, _ int) (string, bool) {
		return item.ContentType, item.ContentType != ""
	}))
	return []*j.Statement{
		j.Op(`
			type Decoder interface {
				Decode(v any) error
			}`),

		j.Add(utils.QualSprintf(`var Decoders = map[string]func(r %Q(io,Reader)) Decoder`)).Values(j.DictFunc(func(d j.Dict) {
			for _, ct := range contentTypes {
				format := getFormatByContentType(ct)
				if format != "" {
					if v, ok := encodingDecoders[format]; ok {
						d[j.Lit(ct)] = v
					}
				} else {
					d[j.Lit(ct)] = j.Add(utils.QualSprintf(`func(_ %Q(io,Reader)) Decoder { panic("No decoder is set for content type %s") }`, ct))
				}
			}
		})),

		j.Add(utils.QualSprintf(`
			func NewDecoder(contentType string, r %Q(io,Reader)) Decoder {
				if v, ok := Decoders[contentType]; ok {
					return v(r)
				}
				panic("Unknown content type " + contentType)
			}`)),
	}
}

func (e EncodingDecode) RenderUsage(_ *common.RenderContext) []*j.Statement {
	panic("not implemented")
}

func (e EncodingDecode) ID() string {
	return "Decode"
}

func (e EncodingDecode) String() string {
	return "EncodingDecode"
}

func getFormatByContentType(contentType string) string {
	// TODO: add other formats: protobuf, avro, etc.
	switch {
	case contentType == "application/json":
		return "application/json"
	case contentType == "application/yaml":
		return "application/yaml"
	case contentType == "application/x-msgpack":
		return "application/x-msgpack"
	case contentType == "text/plain":
		return "text/plain"
	}
	return ""
}
