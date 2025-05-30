package std

import (
	"bufio"
	"context"
	"net"

	"github.com/xcnt/go-asyncapi/run"
	runTCP "github.com/xcnt/go-asyncapi/run/tcp"
)

func NewChannel(conn *net.TCPConn, scanner *bufio.Scanner, maxEnvelopeSize int) *Channel {
	res := Channel{
		TCPConn:         conn,
		scanner:         scanner,
		maxEnvelopeSize: maxEnvelopeSize,
		items:           run.NewFanOut[runTCP.EnvelopeReader](),
	}
	res.ctx, res.cancel = context.WithCancelCause(context.Background())
	go res.run()
	return &res
}

type Channel struct {
	*net.TCPConn
	scanner         *bufio.Scanner
	maxEnvelopeSize int
	items           *run.FanOut[runTCP.EnvelopeReader]
	ctx             context.Context
	cancel          context.CancelCauseFunc
}

type ImplementationRecord interface {
	Bytes() []byte
}

func (c *Channel) Send(_ context.Context, envelopes ...runTCP.EnvelopeWriter) error {
	for _, envelope := range envelopes {
		ir := envelope.(ImplementationRecord)
		if _, err := c.TCPConn.Write(ir.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

func (c *Channel) Receive(ctx context.Context, cb func(envelope runTCP.EnvelopeReader)) error {
	el := c.items.Add(cb)
	defer c.items.Remove(el)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.ctx.Done():
		return context.Cause(c.ctx)
	}
}

func (c *Channel) Close() error {
	c.cancel(nil)
	return c.TCPConn.Close()
}

func (c *Channel) run() {
	for {
		// TODO: oob
		buf := make([]byte, c.maxEnvelopeSize) // TODO: sync.Pool
		switch {
		case c.scanner != nil:
			c.scanner.Buffer(buf, c.maxEnvelopeSize)
			if !c.scanner.Scan() {
				c.cancel(c.scanner.Err())
				return
			}
			c.items.Put(func() runTCP.EnvelopeReader { return NewEnvelopeIn(c.scanner.Bytes()) })
		default:
			n, err := c.TCPConn.Read(buf)
			if err != nil {
				c.cancel(err)
				return
			}
			c.items.Put(func() runTCP.EnvelopeReader { return NewEnvelopeIn(buf[:n]) })
		}
	}
}
