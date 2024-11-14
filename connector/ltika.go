package connector

import (
	"context"

	"github.com/FlowingSPDG/rdeck/input"
	"github.com/FlowingSPDG/rdeck/output"
	"golang.org/x/xerrors"
)

type ltika struct {
	in  input.Input[bool]
	out output.Analog
}

func NewLTikaConnector(in input.Input[bool], out output.Analog) Connector {
	return &ltika{
		in:  in,
		out: out,
	}
}

func (l *ltika) Start(ctx context.Context) error {
	data, err := l.in.Listen(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case e := <-err:
			return xerrors.Errorf("Unknown error for e: %w", e)
		case d := <-data:
			if d {
				if err := l.out.On(); err != nil {
					return err
				}
				continue
			}
			if err := l.out.Off(); err != nil {
				return err
			}
			continue
		}
	}
}
