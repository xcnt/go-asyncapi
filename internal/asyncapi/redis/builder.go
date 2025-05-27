package redis

import (
	"github.com/xcnt/go-asyncapi/internal/asyncapi"
)

type ProtoBuilder struct {
	asyncapi.BaseProtoBuilder
}

var Builder = ProtoBuilder{
	BaseProtoBuilder: asyncapi.BaseProtoBuilder{
		ProtoName:  "redis",
		ProtoTitle: "Redis",
	},
}
