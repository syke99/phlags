package phlags

import "os"

type phlagSet map[string]any

var flgs phlagSet

type PhlagSet struct {
	cmd string
	set []*Phlag
}

func NewSet(cmd string) *PhlagSet {
	if flgs == nil {
		flgs = make(phlagSet)
		flgs["plagBase"] = &PhlagSet{
			cmd: os.Args[0],
			set: make([]*Phlag, 0),
		}
	}

	ps := &PhlagSet{
		cmd: cmd,
		set: make([]*Phlag, 0),
	}

	flgs[cmd] = ps

	return ps
}

func AddBaseSetPhlag(phlag *Phlag) {
	set := flgs["plagBase"].(*PhlagSet)

	set.set = append(set.set, phlag)
}

func (ps *PhlagSet) AddPhlag(phlag *Phlag) *PhlagSet {
	ps.set = append(ps.set, phlag)
	return ps
}
