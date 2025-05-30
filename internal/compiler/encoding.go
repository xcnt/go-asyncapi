package compiler

import (
	"github.com/xcnt/go-asyncapi/internal/common"
	"github.com/xcnt/go-asyncapi/internal/render"
)

const encodingPackageName = "encoding"

func EncodingCompile(ctx *common.CompileContext) error {
	ctx.Logger.Trace("Encoding package")
	e, d := buildEncoderDecoder(ctx)
	ctx.Storage.AddObject(encodingPackageName, ctx.PathStack(), e)
	ctx.Storage.AddObject(encodingPackageName, ctx.PathStack(), d)
	return nil
}

func buildEncoderDecoder(ctx *common.CompileContext) (*render.EncodingEncode, *render.EncodingDecode) {
	allMessagesPrm := render.NewListCbPromise[*render.Message](func(item common.Renderer, _ []string) bool {
		_, ok := item.(*render.Message)
		return ok
	})
	ctx.PutListPromise(allMessagesPrm)

	return &render.EncodingEncode{
			AllMessages:        allMessagesPrm,
			DefaultContentType: ctx.Storage.DefaultContentType(),
		}, &render.EncodingDecode{
			AllMessages:        allMessagesPrm,
			DefaultContentType: ctx.Storage.DefaultContentType(),
		}
}
