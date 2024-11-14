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
	DetermineByActs(*vmixtcp.ActsResponse) ShouldTally
}

type settings struct {
	target int
}

type vmixTallyDeterminer struct {
	s settings
}

func NewvMixTallyDeterminer(target int) VMixTallyDeterminer {
	return &vmixTallyDeterminer{
		s: settings{
			target: target,
		},
	}
}

func (v *vmixTallyDeterminer) DetermineByActs(resp *vmixtcp.ActsResponse) ShouldTally {
	// TODO.
	return ShouldTally{}
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
