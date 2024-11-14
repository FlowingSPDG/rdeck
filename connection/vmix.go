package connection

import (
	"context"
	"time"

	"github.com/FlowingSPDG/rdeck/device/vmix"
	"github.com/FlowingSPDG/rdeck/input"
	"github.com/FlowingSPDG/rdeck/output"
	"golang.org/x/xerrors"

	vmixtcp "github.com/FlowingSPDG/vmix-go/tcp"
)

type VMixConnection interface {
	Start(ctx context.Context) error
	ToTallyInput() input.Input[*vmixtcp.TallyResponse]
	ToOutput() output.VMixOutput
}

func NewVMixConnection(addr string) VMixConnection {
	v := vmixtcp.New(addr)
	return &vMixConnection{
		v: v,
	}
}

type vMixConnection struct {
	v vmixtcp.Vmix
}

// Start implements VMixConnection.
func (v *vMixConnection) Start(ctx context.Context) error {
	if err := v.v.Connect(ctx, 3*time.Second); err != nil {
		return xerrors.Errorf("failed to connect: %w", err)
	}
	if err := v.v.Run(ctx); err != nil {
		return xerrors.Errorf("unknown error on running vmix: %w", err)
	}
	return nil
}

// ToInput implements VMixConnection.
func (v *vMixConnection) ToTallyInput() input.Input[*vmixtcp.TallyResponse] {
	return vmix.NewvMixTallyInput(v.v)
}

// ToOutput implements VMixConnection.
func (v *vMixConnection) ToOutput() output.VMixOutput {
	return vmix.NewVMixOutput(v.v)
}
