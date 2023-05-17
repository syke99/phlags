package phlags

import (
	"errors"
	"github.com/syke99/phlags/internal"
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

func New(abrv string, name string, usage string, defaultValue any) (*Phlag, error) {
	if abrv == "" &&
		name == "" {
		return nil, internal.InvalidName
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
			return nil, internal.InvalidName
		}

		pShort[abrv] = flg
	}

	if name != "" {
		name, err = sanitizeName(0, name)
		if err != nil {
			return nil, internal.InvalidName
		}

		pShort[name] = flg
	}

	return flg, err
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
