package vmix

import (
	"github.com/FlowingSPDG/rdeck/output"
	vmixtcp "github.com/FlowingSPDG/vmix-go/tcp"
	"golang.org/x/xerrors"
)

type vMixOutput struct {
	// vMix as an output
	v vmixtcp.Vmix
}

func (v *vMixOutput) SendFunction(name, query string) error {
	if err := v.v.Function(name, query); err != nil {
		return xerrors.Errorf("failed to send function: %w", err)
	}
	return nil
}

func NewVMixOutput(v vmixtcp.Vmix) output.VMixOutput {
	return &vMixOutput{
		v: v,
	}
}
