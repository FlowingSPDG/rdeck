package determiner

import (
	vmixtcp "github.com/FlowingSPDG/vmix-go/tcp"
)

type ShouldTally struct {
	Preview bool
	Program bool
}

type VMixTallyDeterminer interface {
	DetermineByTally(*vmixtcp.TallyResponse) ShouldTally
}

type vmixTallyDeterminerSettings struct {
	target int
}

type vmixTallyDeterminer struct {
	s vmixTallyDeterminerSettings
}

func NewvMixTallyDeterminer(target int) VMixTallyDeterminer {
	return &vmixTallyDeterminer{
		s: vmixTallyDeterminerSettings{
			target: target,
		},
	}
}

func (v *vmixTallyDeterminer) DetermineByTally(resp *vmixtcp.TallyResponse) ShouldTally {
	if len(resp.Tally) < v.s.target-1 {
		return ShouldTally{}
	}
	switch resp.Tally[v.s.target-1] {
	case vmixtcp.Off:
		return ShouldTally{}
	case vmixtcp.Preview:
		return ShouldTally{
			Preview: true,
			Program: false,
		}
	case vmixtcp.Program:
		return ShouldTally{
			Preview: false,
			Program: true,
		}
	default:
		return ShouldTally{}
	}
}

type vMixActivatorDeterminer struct {
	s vMixActivatorDeterminerSettings
}

type vMixActivatorDeterminerSettings struct {
	// TODO...
	rawText string
}

type VMixActivatorDeterminer interface {
	DetermineByActs(*vmixtcp.ActsResponse) ShouldTally
}

// DetermineByActs implements VMixActivatorDeterminer.
func (v *vMixActivatorDeterminer) DetermineByActs(resp *vmixtcp.ActsResponse) ShouldTally {
	if v.s.rawText == resp.Response {
		return ShouldTally{
			Program: true,
		}
	}
	return ShouldTally{
		Program: false,
	}
}

func NewVMixActivatorDeterminer(rawTarget string) VMixActivatorDeterminer {
	return &vMixActivatorDeterminer{
		s: vMixActivatorDeterminerSettings{
			rawText: rawTarget,
		},
	}
}
