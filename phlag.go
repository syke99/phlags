package phlags

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type phlgs[T any] map[string]any

var pShort phlgs[any]
var pFull phlgs[any]

func Init() {
	pShort = make(map[string]any)
	pFull = make(map[string]any)
}

type Phlag[T any, A any] struct {
	usage  string
	args   []A
	dValue T
	value  T
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

func New[T any, A any](abrv string, name string, usage string, defaultValue T) *Phlag[T, A] {
	flg := &Phlag[T, A]{
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

func Parse[T any, A any]() error {
	flgSplits := make([]int, 0)

	// find the os.Args index of each flag
	for i := range os.Args {
		v := os.Args[i]
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
			sections[i] = os.Args[split:]
			continue
		}

		// otherwise, grab the args between
		// this section's beginning and the
		// next section
		sections[i] = os.Args[split:next]
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

		var val T

		// separate the flag name from the value provided
		if n, v, found := strings.Cut(name, "="); found {
			name = n
			val = any(v).(T)
		}

		var f *Phlag[T, A]

		// find which map the flag was stored in
		// based on the name, either full or
		// abbreviated
		if phlg, ok := pFull[name]; ok {
			f = phlg.(*Phlag[T, A])
		}

		if phlg, ok := pShort[name]; ok {
			f = phlg.(*Phlag[T, A])
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
			f.args = make([]A, len(section)-1)

			positional := section[1:]

			for pI := range positional {
				pV := positional[pI]

				f.args[pI] = any(pV).(A)
			}
		}
	}

	return nil
}

func (p *Phlag[T, A]) String() *string {
	if any(p.value) == nil &&
		any(p.dValue) != nil {
		v := any(p.dValue).(string)

		return &v
	}

	switch any(p.value).(type) {
	case string:
		v := any(p.value).(string)

		return &v
	}
	return nil
}

func (p *Phlag[T, A]) Int() *int {
	if any(p.value) == nil &&
		any(p.dValue) != nil {
		v := any(p.dValue).(int)

		return &v
	}

	switch any(p.value).(type) {
	case string:
		v := any(p.value).(int)

		return &v
	}
	return nil
}

func (p *Phlag[T, A]) Bool() bool {
	switch any(p.value).(type) {
	case bool:
		v := any(p.value).(bool)

		return v
	}
	return false
}

func (p *Phlag[T, A]) Args() []PositionalArgument[A] {
	if p.args != nil ||
		len(p.args) != 0 {
		args := make([]PositionalArgument[A], len(p.args))

		for i, v := range p.args {
			pA := PositionalArgument[A]{
				value: v,
			}

			args[i] = pA
		}

		return args
	}
	return nil
}

type PositionalArgument[T any] struct {
	value T
}

func (pa PositionalArgument[T]) String() *string {
	v := any(pa.value).(string)

	return &v
}

func (pa PositionalArgument[T]) Int() *int {
	v, err := strconv.Atoi(any(pa.value).(string))
	if err != nil {
		return nil
	}

	return &v
}

func (pa PositionalArgument[T]) Bool() bool {
	b, err := strconv.ParseBool(any(pa.value).(string))
	if err != nil {
		return false
	}
	return b
}
