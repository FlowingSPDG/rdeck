package determiner

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	vmixtcp "github.com/FlowingSPDG/vmix-go/tcp"
)

type ShouldTally struct {
	Preview bool
	Program bool
}

type VMixTallyDeterminer interface {
	// DetermineByTally Determines tally by TallyResponse. Potentially returns nil.
	DetermineByTally(*vmixtcp.TallyResponse) *ShouldTally
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

func (v *vmixTallyDeterminer) DetermineByTally(resp *vmixtcp.TallyResponse) *ShouldTally {
	if len(resp.Tally) < v.s.target-1 {
		return nil
	}
	switch resp.Tally[v.s.target-1] {
	case vmixtcp.Off:
		return &ShouldTally{}
	case vmixtcp.Preview:
		return &ShouldTally{
			Preview: true,
			Program: false,
		}
	case vmixtcp.Program:
		return &ShouldTally{
			Preview: false,
			Program: true,
		}
	}
	return nil
}

type vMixActivatorInputDeterminer[T comparable] struct {
	s vMixActivatorInputDeterminerSettings[T]
}

type vMixActivatorInputDeterminerSettings[T comparable] struct {
	target string
	input  int
	state  T
}

type VMixActivatorDeterminer interface {
	// DetermineByActs Determines tally by ActsResponse. Potentially returns nil.
	DetermineByActs(*vmixtcp.ActsResponse) *ShouldTally
}

// DetermineByActs implements VMixActivatorDeterminer.
func (v *vMixActivatorInputDeterminer[T]) DetermineByActs(resp *vmixtcp.ActsResponse) *ShouldTally {
	log.Println("DetermineByActs:", resp)

	strs := strings.Split(resp.Response, " ")
	if len(strs) < 3 {
		return nil
	}
	if strs[0] != v.s.target {
		return nil
	}

	input, err := strconv.Atoi(strs[1])
	if err != nil {
		return nil
	}
	if input != v.s.input {
		return nil
	}

	state := new(T)

	if _, err := fmt.Sscanf(strs[2], "%v", state); err != nil {
		return nil
	}

	st := &ShouldTally{
		Preview: false,
		Program: *state == v.s.state,
	}
	log.Println("DetermineByActs result:", st)
	return st
}

func NewVMixActivatorDeterminer[T comparable](target string, input int, state T) VMixActivatorDeterminer {
	return &vMixActivatorInputDeterminer[T]{
		s: vMixActivatorInputDeterminerSettings[T]{
			target: target,
			input:  input,
			state:  state,
		},
	}
}
