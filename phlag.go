package phlags

import (
	"errors"
	"os"
	"strings"
)

type phlgs map[string]any

var pShort phlgs
var pFull phlgs

type Phlag struct {
	abrv   string
	name   string
	usage  string
	args   []PositionalArgument
	dValue any
	value  any
}

func sanitizeName(name string) string {
	if strings.HasPrefix(name, "-") {
		name = strings.TrimPrefix(name, "-")
	}

	if strings.HasPrefix(name, "-") {
		sanitizeName(name)
	}

	return name
}

func New(abrv string, name string, usage string, defaultValue any) *Phlag {
	if abrv == "" &&
		name == "" {
		return nil
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

	if abrv != "" {
		abrv = sanitizeName(abrv)

		pShort[abrv] = flg
	}

	if name != "" {
		name = sanitizeName(name)

		pShort[name] = flg
	}

	return flg
}

func parseIntoFlags(args []string) error {
	flgSplits := make([]int, 0)

	// find the os.Args index of each flag
	for i := range args {
		v := args[i]
		if strings.HasPrefix(v, "--") ||
			strings.HasPrefix(v, "-") {
			flgSplits = append(flgSplits, i)
		}
	}

	// separate os.Args into sections based on
	// flag indexes
	sections := make([][]string, len(flgSplits))

	for i, split := range flgSplits {
		var next int

		if i != len(flgSplits)-1 {
			next = flgSplits[i+1]
		}

		// if it's the very last section,
		// just grab the remaining os.Args
		if i == len(flgSplits)-1 {
			sections[i] = args[split:]
			continue
		}

		// otherwise, grab the args between
		// this section's beginning and the
		// next section
		sections[i] = args[split:next]
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
		name := section[0]

		// trim dash prefixes from the names
		if strings.HasPrefix(name, "--") {
			name, _ = strings.CutPrefix(name, "--")
		} else if strings.HasPrefix(name, "-") {
			name, _ = strings.CutPrefix(name, "-")
			if strings.HasPrefix(name, "-") {
				return errors.New("invalid flag, too many \"-\"")
			}
		} else {
			return errors.New("invalid flag")
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
		if phlg, ok := pFull[name]; ok {
			f = phlg.(*Phlag)
		}

		if phlg, ok := pShort[name]; ok {
			f = phlg.(*Phlag)
		}

		// set the value to the flag
		if any(val) != nil {
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

				f.args[pI] = PositionalArgument{
					value: pV,
				}
			}
		}
	}

	return nil
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
		cmd := os.Args[1]

		ps := flgs[cmd].(*PhlagSet).set

		// override pFull and pShort
		// with only the appropriate
		// flags
		pShort = make(phlgs)
		pShort = make(phlgs)

		for i := range ps {
			if ps[i].abrv != "" {
				v := ps[i]
				abrv := sanitizeName(ps[i].abrv)
				pShort[abrv] = v
			}

			if ps[i].name != "" {
				v := ps[i]
				name := sanitizeName(ps[i].name)
				pFull[name] = v
			}
		}

		args = os.Args[1:]
	}

	short := pShort
	long := pFull

	flgSplits := make([]int, 0)

	// find the os.Args index of each flag
	for i := range args {
		v := args[i]
		if strings.HasPrefix(v, "--") ||
			strings.HasPrefix(v, "-") {
			flgSplits = append(flgSplits, i)
		}
	}

	// separate os.Args into sections based on
	// flag indexes
	sections := make([][]string, len(flgSplits))

	for i, split := range flgSplits {
		var next int

		if i != len(flgSplits)-1 {
			next = flgSplits[i+1]
		}

		// if it's the very last section,
		// just grab the remaining os.Args
		if i == len(flgSplits)-1 {
			sections[i] = args[split:]
			continue
		}

		// otherwise, grab the args between
		// this section's beginning and the
		// next section
		sections[i] = args[split:next]
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
		name := section[0]

		// trim dash prefixes from the names
		if strings.HasPrefix(name, "--") {
			name, _ = strings.CutPrefix(name, "--")
		} else if strings.HasPrefix(name, "-") {
			name, _ = strings.CutPrefix(name, "-")
			if strings.HasPrefix(name, "-") {
				return errors.New("invalid flag, too many \"-\"")
			}
		} else {
			return errors.New("invalid flag")
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
		if phlg, ok := long[name]; ok {
			f = phlg.(*Phlag)
		}

		if phlg, ok := short[name]; ok {
			f = phlg.(*Phlag)
		}

		if f == nil {
			continue
		}

		// set the value to the flag
		if val != nil {
			f.value = val.(string)
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

				f.args[pI] = PositionalArgument{
					value: pV,
				}
			}
		}
	}

	pFull = nil
	pFull = nil
	flgs = nil

	return nil
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
