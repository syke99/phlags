package phlags

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type phlgs map[string]any

var pShort phlgs
var pFull phlgs

type Phlag struct {
	present bool
	abrv    string
	name    string
	usage   string
	args    []PositionalArgument
	dValue  any
	value   any
}

func sanitizeName(counter int, name string) (string, error) {

	if counter == 0 &&
		!strings.HasPrefix(name, "-") {
		return "", errors.New("invalid flag")
	}

	if counter == 2 {
		return name, nil
	}

	counter++

	if strings.HasPrefix(name, "-") {
		name = strings.TrimPrefix(name, "-")
	}

	return sanitizeName(counter, name)
}

func New(abrv string, name string, usage string, defaultValue any) (*Phlag, error) {
	if abrv == "" &&
		name == "" {
		return nil, errors.New("invalid name provided")
	}

	if pFull == nil {
		pFull = make(phlgs)
	}

	if pShort == nil {
		pShort = make(phlgs)
	}

	flg := &Phlag{
		abrv:   abrv,
		name:   name,
		usage:  usage,
		dValue: defaultValue,
	}

	var err error

	if abrv != "" {
		abrv, err = sanitizeName(0, abrv)
		if err != nil {
			return nil, errors.New("invalid name provided")
		}

		pShort[abrv] = flg
	}

	if name != "" {
		name, err = sanitizeName(0, name)
		if err != nil {
			return nil, errors.New("invalid name provided")
		}

		pShort[name] = flg
	}

	return flg, err
}

func parseIntoFlags(args []string, flgSplits []int, short phlgs, full phlgs) error {
	// separate os.Args into sections based on
	// flag indexes
	sections := make([][]string, len(flgSplits))

	for i, split := range flgSplits {
		current := split - 1

		if i != len(flgSplits)-1 {
			next := flgSplits[i+1] - 1

			sections[i] = args[current:next]
			continue
		}

		sections[i] = args[current:]
	}

	// loop through the sections and assign
	// all args and values to the appropriate
	// Phlag
	for _, section := range sections {
		// if a section is empty for some reason, break
		if len(section) == 0 {
			break
		}

		// the name (and value) of the flag
		// will be the first value in a section
		name, err := sanitizeName(0, section[0])
		if err != nil {
			return errors.New("invalid name provided")
		}

		if strings.HasPrefix(name, "-") {
			if strings.HasPrefix(name, "-") {
				return errors.New("invalid flag, too many \"-\"")
			}
		}

		var val any

		// separate the flag name from the value provided
		if n, v, found := strings.Cut(name, "="); found {
			name = n
			val = any(v)
		}

		var f *Phlag

		// find which map the flag was stored in
		// based on the name, either full or
		// abbreviated
		if phlg, ok := full[name]; ok {
			f = phlg.(*Phlag)
			f.present = true
		}

		if phlg, ok := short[name]; ok {
			f = phlg.(*Phlag)
			f.present = true
		}

		// set the value to the flag
		if val != nil &&
			f != nil {
			f.value = val
		}

		// if the length of the section
		// is greater than 1, that means
		// positional arguments
		// were provided to the flag
		if len(section) > 1 {
			f.args = make([]PositionalArgument, len(section)-1)

			positional := section[1:]

			for pI := range positional {
				pV := positional[pI]

				// make sure the trailing values are a subcommand
				if _, ok := flgs[pV]; !ok {
					f.args[pI] = PositionalArgument{
						value: positional[pI],
					}
					continue
				}
			}
		}
	}

	return nil
}

func overridePShortPFullSanitized(ps []*Phlag) (phlgs, phlgs) {
	short := make(phlgs)
	full := make(phlgs)

	for pI := range ps {
		if ps[pI].abrv != "" {
			abrv, err := sanitizeName(0, ps[pI].abrv)
			if err != nil {
				return nil, nil
			}
			short[abrv] = ps[pI]
		}

		if ps[pI].name != "" {
			name, err := sanitizeName(0, ps[pI].name)
			if err != nil {
				return nil, nil
			}
			full[name] = ps[pI]
		}
	}
	return short, full
}

func separateByPhlagSet(sections []*setGroup, args []string, group *setGroup) []*setGroup {
	if len(args) == 0 {
		return sections
	}

	// the first pass through, signaling possible
	// base flags, group will be nil, so set it
	// to a pointer to a newly initialized
	// setGroup
	if group == nil {
		group = &setGroup{
			args: nil,
			ps:   flgs["plagBase"].(*PhlagSet),
		}
	}

	// create a slice to hold all the arguments
	// for a setGroup
	vals := make([]string, 0)

	recurse := false
	var ptr any
	var ok bool

	// loop through all the arguments
	for i, arg := range args {
		// if the argument isn't a
		// subcommand (a key in flgs
		// set when calling NewSet)
		// build out vals
		if ptr, ok = flgs[arg]; !ok {
			// if the arg isn't a subcommand
			// and the end of args has been
			// reached, append the group to
			// the setGroup
			if i == len(args)-1 {
				vals = append(vals, arg)
				group.args = vals
				sections = append(sections, group)
				break
			}

			vals = append(vals, arg)
		} else {
			// else slice the remaining args,
			// not including the subcommand
			// arg, and recurse
			args = args[i+1:]
			recurse = true
			break
		}
	}

	if recurse {
		sections = separateByPhlagSet(sections, args, &setGroup{
			args: nil,
			ps:   ptr.(*PhlagSet),
		})
	}

	return sections
}

type setGroup struct {
	args []string
	ps   *PhlagSet
}

func Parse() error {
	args := os.Args

	// if flag sets were made, then grab the
	// appropriate set of flags by extracting
	// the subcommand and then passing it to
	// flgs, then shift the args to not include
	// the subcommand and then parse the flag
	// values (and positional args if they exist)
	// to the appropriate flags
	if len(flgs) != 0 {

		var err error

		args = args[1:]

		phlagSetSections := separateByPhlagSet(make([]*setGroup, 0), args, nil)

		for _, ps := range phlagSetSections {
			short, long := overridePShortPFullSanitized(ps.ps.set)
			if short == nil &&
				long == nil {
				return fmt.Errorf("no phlags set for PhlagSet: %s", ps.ps.cmd)
			}

			flgSplits := findSplits(args)

			err = parseIntoFlags(ps.args, flgSplits, short, long)
			if err != nil {
				break
			}
		}

		return err
	}

	flgSplits := findSplits(args)

	err := parseIntoFlags(args, flgSplits, pFull, pShort)
	if err != nil {
		cleanUp()
		return err
	}

	cleanUp()

	return nil
}

func findSplits(args []string) []int {
	flgSplits := make([]int, 0)

	// find the index of each flag
	for i := range args {
		v := args[i]
		if strings.HasPrefix(v, "--") ||
			strings.HasPrefix(v, "-") {
			flgSplits = append(flgSplits, i)
		}
	}

	return flgSplits
}

func cleanUp() {
	pFull = nil
	pFull = nil
	flgs = nil
}

func (p *Phlag) String() *string {
	var v string
	if p.value == nil &&
		p.dValue != nil {

		switch p.dValue.(type) {
		case string:
			v = p.dValue.(string)
			// TODO: format other types to string
		}
		return &v
	}

	switch p.value.(type) {
	case string:
		v = p.value.(string)

		return &v
	}
	return nil
}

func (p *Phlag) Int() *int {
	var v int
	if p.value == nil &&
		p.dValue != nil {
		switch p.dValue.(type) {
		case string:
			v = p.dValue.(int)
			// TODO: format other types to int
		}

		return &v
	}

	switch any(p.value).(type) {
	case string:
		v = any(p.value).(int)
		// TODO: format other types to int
		return &v
	}
	return nil
}

func (p *Phlag) PositionalArguments() []PositionalArgument {
	if p.args != nil ||
		len(p.args) != 0 {
		args := make([]PositionalArgument, len(p.args))

		for i, v := range p.args {
			args[i] = v
		}

		return args
	}
	return nil
}

func (p *Phlag) PositionalArgumentByIndex(idx int) (PositionalArgument, error) {
	if idx >= len(p.args) {
		return PositionalArgument{}, errors.New("provided index out of range of Phlag Positional Arguments")
	}

	return p.args[idx], nil
}
