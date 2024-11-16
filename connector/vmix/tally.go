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

type Target int

const (
	_ Target = iota
	Preview
	Program
)

// TODO: 実際の接続とリトライ処理
// TODO: 複数のOutputの保持
func NewVMixTallyConnector(
	v connection.VMixConnection,
	out output.Digital,
	determiner determiner.VMixTallyDeterminer,
	settings VMixTallyConnectorSettings,
) connector.Connector {
	return &vMixTallyConnector{
		in:         v.ToTallyInput(),
		out:        out,
		determiner: determiner,
		settings:   settings,
	}
}

type vMixTallyConnector struct {
	// in and out
	in  input.Input[*vmixtcp.TallyResponse]
	out output.Digital

	// determiner
	determiner determiner.VMixTallyDeterminer

	settings VMixTallyConnectorSettings
}

type VMixTallyConnectorSettings struct {
	Target Target
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
			if sd == nil {
				continue
			}
			log.Println("Determined tally for:", sd)

			switch t.settings.Target {

			case Preview:
				if sd.Preview {
					if err := t.out.On(); err != nil {
						return xerrors.Errorf("failed to turn on Preview tally light: %w", err)
					}
				} else {
					if err := t.out.Off(); err != nil {
						return xerrors.Errorf("failed to turn off Preview tally light: %w", err)
					}
				}

			case Program:
				if sd.Program {
					if err := t.out.On(); err != nil {
						return xerrors.Errorf("failed to turn on Program tally light: %w", err)
					}
				} else {
					if err := t.out.Off(); err != nil {
						return xerrors.Errorf("failed to turn off Program tally light: %w", err)
					}
				}
			}
		}
	}
}
