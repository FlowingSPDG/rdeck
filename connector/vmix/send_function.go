package vmix

import (
	"context"

	"github.com/FlowingSPDG/rdeck/connection"
	"github.com/FlowingSPDG/rdeck/connector"
	"github.com/FlowingSPDG/rdeck/input"
	"github.com/FlowingSPDG/rdeck/output"

	"golang.org/x/xerrors"
)

type sendFunction struct {
	in       input.Input[bool]
	out      output.VMixOutput
	settings sendFunctionSettings
}

type sendFunctionSettings struct {
	name  string
	query string
}

func NewSendFunction(
	in input.Input[bool],
	v connection.VMixConnection,
	funcName, query string,
) connector.Connector {
	return &sendFunction{
		in:  in,
		out: v.ToOutput(),
		settings: sendFunctionSettings{
			name:  funcName,
			query: query,
		},
	}
}

func (s *sendFunction) Start(ctx context.Context) error {
	data, err := s.in.Listen(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case e := <-err:
			return xerrors.Errorf("Unknown error for e: %w", e)
		case d := <-data:
			if !d {
				continue
			}
			if err := s.out.SendFunction(s.settings.name, s.settings.query); err != nil {
				return xerrors.Errorf("failed to send function: %w", err)
			}
			continue
		}
	}
}
