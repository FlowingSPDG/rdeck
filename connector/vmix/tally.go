package vmix

import (
	"context"
	"log"

	"github.com/FlowingSPDG/rdeck/connector"
	"github.com/FlowingSPDG/rdeck/determiner"
	"github.com/FlowingSPDG/rdeck/input"
	"github.com/FlowingSPDG/rdeck/output"

	vmixtcp "github.com/FlowingSPDG/vmix-go/tcp"

	"golang.org/x/xerrors"
)

// TODO: 実際の接続とリトライ処理
// TODO: 複数のOutputの保持
func NewVMixTallyConnector(
	in input.Input[*vmixtcp.TallyResponse],
	out output.Analog,
	determiner determiner.VMixTallyDeterminer,
) connector.Connector {
	return &vMixTallyConnector{
		in:         in,
		out:        out,
		determiner: determiner,
	}
}

type vMixTallyConnector struct {
	// in and out
	in  input.Input[*vmixtcp.TallyResponse]
	out output.Analog

	// determiner
	determiner determiner.VMixTallyDeterminer
}

func (t *vMixTallyConnector) Start(ctx context.Context) error {
	log.Println("STARTING vMixTallyConnecor.")
	log.Println("LISTENING...")
	data, err := t.in.Listen(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case e := <-err:
			return xerrors.Errorf("Unknown error for e: %w", e)
		case d := <-data:
			log.Println("RECEIVED DATA on vMixTallyConnecor.Start(). data:", d)
			sd := t.determiner.DetermineByTally(d)
			log.Println("Determined tally for:", sd)

			if sd.Program {
				if err := t.out.On(); err != nil {
					return xerrors.Errorf("failed to turn on tally light: %w", err)
				}
				continue
			}
			if err := t.out.Off(); err != nil {
				return xerrors.Errorf("failed to turn on tally light: %w", err)
			}
			continue
		}
	}
}
