package tcp

import (
	"context"
	"io"

	"github.com/xcnt/go-asyncapi/run"
)

// Pub
type (
	Producer interface {
		Publisher(ctx context.Context, channelName string, bindings *ChannelBindings) (Publisher, error)
	}
	Publisher interface {
		Send(ctx context.Context, envelopes ...EnvelopeWriter) error
		Close() error
	}
	EnvelopeWriter interface {
		io.Writer
		ResetPayload()
		SetHeaders(headers run.Headers)
		SetContentType(contentType string)
		SetBindings(bindings MessageBindings)
	}
)

type EnvelopeMarshaler interface {
	MarshalTCPEnvelope(envelope EnvelopeWriter) error
}

// Sub
type (
	Consumer interface {
		Subscriber(ctx context.Context, channelName string, bindings *ChannelBindings) (Subscriber, error)
	}
	Subscriber interface {
		Receive(ctx context.Context, cb func(envelope EnvelopeReader)) error
		Close() error
	}
	EnvelopeReader interface {
		io.Reader
		Headers() run.Headers
	}
)

type EnvelopeUnmarshaler interface {
	UnmarshalTCPEnvelope(envelope EnvelopeReader) error
}
