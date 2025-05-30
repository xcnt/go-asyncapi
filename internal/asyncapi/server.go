package asyncapi

import (
	"github.com/xcnt/go-asyncapi/internal/specurl"
	"github.com/samber/lo"

	"github.com/xcnt/go-asyncapi/internal/types"

	"github.com/xcnt/go-asyncapi/internal/common"
	"github.com/xcnt/go-asyncapi/internal/render"
)

type Server struct {
	URL             string                                   `json:"url" yaml:"url"`
	Protocol        string                                   `json:"protocol" yaml:"protocol"`
	ProtocolVersion string                                   `json:"protocolVersion" yaml:"protocolVersion"`
	Description     string                                   `json:"description" yaml:"description"`
	Variables       types.OrderedMap[string, ServerVariable] `json:"variables" yaml:"variables"`
	Security        []SecurityRequirement                    `json:"security" yaml:"security"`
	Tags            []Tag                                    `json:"tags" yaml:"tags"`
	Bindings        *ServerBindings                          `json:"bindings" yaml:"bindings"`

	XGoName string `json:"x-go-name" yaml:"x-go-name"`
	XIgnore bool   `json:"x-ignore" yaml:"x-ignore"`

	Ref string `json:"$ref" yaml:"$ref"`
}

func (s Server) Compile(ctx *common.CompileContext) error {
	ctx.RegisterNameTop(ctx.Stack.Top().PathItem)
	obj, err := s.build(ctx, ctx.Stack.Top().PathItem)
	if err != nil {
		return err
	}
	ctx.PutObject(obj)

	return nil
}

func (s Server) build(ctx *common.CompileContext, serverKey string) (common.Renderer, error) {
	_, isComponent := ctx.Stack.Top().Flags[common.SchemaTagComponent]
	ignore := s.XIgnore || !ctx.CompileOpts.ServerOpts.IsAllowedName(serverKey)
	if ignore {
		ctx.Logger.Debug("Server denoted to be ignored")
		return &render.Server{Dummy: true}, nil
	}
	if s.Ref != "" {
		ctx.Logger.Trace("Ref", "$ref", s.Ref)
		prm := render.NewRendererPromise(s.Ref, common.PromiseOriginUser)
		// Set a server to be rendered if we reference it from `servers` document section
		prm.DirectRender = !isComponent
		ctx.PutPromise(prm)
		return prm, nil
	}

	srvName, _ := lo.Coalesce(s.XGoName, serverKey)
	// Render only the servers defined directly in `servers` document section, not in `components`
	res := render.Server{
		Name:            srvName,
		RawName:         serverKey,
		GolangName:      ctx.GenerateObjName(srvName, ""),
		DirectRender:    !isComponent,
		URL:             s.URL,
		Protocol:        s.Protocol,
		ProtocolVersion: s.ProtocolVersion,
	}

	// Channels which are connected to this server
	prms := lo.Map(ctx.Storage.ActiveChannels(), func(item string, _ int) *render.Promise[*render.Channel] {
		ref := specurl.BuildRef("channels", item)
		prm := render.NewPromise[*render.Channel](ref, common.PromiseOriginInternal)
		ctx.PutPromise(prm)
		return prm
	})
	res.AllChannelsPromises = prms

	// Bindings
	if s.Bindings != nil {
		ctx.Logger.Trace("Server bindings")
		res.BindingsStruct = &render.GoStruct{
			BaseType: render.BaseType{
				Name:         ctx.GenerateObjName(srvName, "Bindings"),
				DirectRender: true,
				Import:       ctx.CurrentPackage(),
			},
			Fields: nil,
		}

		ref := ctx.PathStackRef("bindings")
		res.BindingsPromise = render.NewPromise[*render.Bindings](ref, common.PromiseOriginInternal)
		ctx.PutPromise(res.BindingsPromise)
	}

	// Server variables
	for _, v := range s.Variables.Entries() {
		ctx.Logger.Trace("Server variable", "name", v.Key)
		ref := ctx.PathStackRef("variables", v.Key)
		prm := render.NewPromise[*render.ServerVariable](ref, common.PromiseOriginInternal)
		ctx.PutPromise(prm)
		res.Variables.Set(v.Key, prm)
	}

	protoBuilder, ok := ProtocolBuilders[s.Protocol]
	if !ok {
		ctx.Logger.Warn("Skip unsupported server protocol", "proto", s.Protocol)
	} else {
		var err error
		ctx.Logger.Trace("Server", "proto", protoBuilder.ProtocolName())
		res.ProtoServer, err = protoBuilder.BuildServer(ctx, &s, &res)
		if err != nil {
			return nil, err
		}

		// Register protocol only for servers in `servers` document section, not in `components`
		if !isComponent {
			ctx.Storage.RegisterProtocol(s.Protocol)
		}
	}

	return &res, nil
}
