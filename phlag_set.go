package phlags

type phlagSet map[string]any

var flgs phlagSet

type PhlagSet struct {
	set []*Phlag
}

func NewSet(cmd string) *PhlagSet {
	if flgs == nil {
		flgs = make(phlagSet)
	}

	ps := &PhlagSet{
		set: make([]*Phlag, 0),
	}

	flgs[cmd] = ps

	return ps
}

func (ps *PhlagSet) AddPhlag(phlag *Phlag) *PhlagSet {
	ps.set = append(ps.set, phlag)
	return ps
}
