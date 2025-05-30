package std

import (
	"bufio"
	"context"
	"net"
	"net/url"

	runTCP "github.com/xcnt/go-asyncapi/run/tcp"
)

func NewProducer(serverURL string) (*ProduceClient, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}
	la, err := net.ResolveTCPAddr(u.Scheme, u.Host)
	if err != nil {
		return nil, err
	}
	d := net.Dialer{LocalAddr: la}

	return &ProduceClient{
		Dialer:          d,
		Scanner:         bufio.NewScanner(nil),
		MaxEnvelopeSize: DefaultMaxEnvelopeSize,
		address:         u.Host,
		protocolFamily:  u.Scheme,
	}, nil
}

type ProduceClient struct {
	net.Dialer
	// Scanner splits the incoming data into Envelopes. If equal to nil, the data is
	// split on chunks of MaxEnvelopeSize bytes, which is equal to bufio.MaxScanTokenSize by default.
	Scanner         *bufio.Scanner
	MaxEnvelopeSize int

	address        string
	protocolFamily string
}

func (p ProduceClient) Publisher(ctx context.Context, _ string, _ *runTCP.ChannelBindings) (runTCP.Publisher, error) {
	conn, err := p.DialContext(ctx, p.protocolFamily, p.address)
	if err != nil {
		return nil, err
	}

	return NewChannel(conn.(*net.TCPConn), p.Scanner, p.MaxEnvelopeSize), nil
}
