package vmix

import (
	"context"
	"log"

	"github.com/FlowingSPDG/rdeck/connection"
	"github.com/FlowingSPDG/rdeck/connector"
	"github.com/FlowingSPDG/rdeck/determiner"
	"github.com/FlowingSPDG/rdeck/input"
	"github.com/FlowingSPDG/rdeck/output"
	vmixtcp "github.com/FlowingSPDG/vmix-go/tcp"
	"golang.org/x/xerrors"
)

type vmixActivatorConnector struct {
	// in and out
	in  input.Input[*vmixtcp.ActsResponse]
	out output.Analog

	// determiner
	determiner determiner.VMixActivatorDeterminer
}

func NewVMixActivatorConnector(
	v connection.VMixConnection,
	out output.Analog,
	determiner determiner.VMixActivatorDeterminer,
) connector.Connector {
	return &vmixActivatorConnector{
		in:         v.ToActivatorInput(),
		out:        out,
		determiner: determiner,
	}
}

func (v *vmixActivatorConnector) Start(ctx context.Context) error {
	log.Println("STARTING vmixActivatorConnector.")
	log.Println("LISTENING...")
	data, err := v.in.Listen(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case e := <-err:
			return xerrors.Errorf("Unknown error for e: %w", e)
		case d := <-data:
			log.Println("RECEIVED DATA on vmixActivatorConnector.Start(). data:", d)
			sd := v.determiner.DetermineByActs(d)
			if sd == nil {
				continue
			}
			log.Println("Determined activator for:", sd)

			if sd.Program {
				if err := v.out.On(); err != nil {
					return xerrors.Errorf("failed to turn on activator light: %w", err)
				}
				continue
			}
			if err := v.out.Off(); err != nil {
				return xerrors.Errorf("failed to turn on activator light: %w", err)
			}
			continue
		}
	}
}
