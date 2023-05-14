package phlags

import (
	"strconv"
)

// TODO: certain methods should really be changed to bitwise operations
// TODO: so they're accurate and efficiently stored (this will require
// TODO: actually learning bitwise operations ;P )

type PositionalArgument struct {
	value string
}

func (pa PositionalArgument) String() *string {
	return &pa.value
}

func (pa PositionalArgument) Int() *int {
	v, err := strconv.Atoi(pa.value)
	if err != nil {
		return nil
	}

	return &v
}

func (pa PositionalArgument) Int8() *int8 {
	v, err := strconv.Atoi(pa.value)
	if err != nil {
		return nil
	}

	i := int8(v)

	return &i
}

func (pa PositionalArgument) Int16() *int16 {
	v, err := strconv.Atoi(pa.value)
	if err != nil {
		return nil
	}

	i := int16(v)

	return &i
}

func (pa PositionalArgument) Int32() *int32 {
	v, err := strconv.Atoi(pa.value)
	if err != nil {
		return nil
	}

	i := int32(v)

	return &i
}

func (pa PositionalArgument) Int64() *int64 {
	v, err := strconv.Atoi(pa.value)
	if err != nil {
		return nil
	}

	i := int64(v)

	return &i
}

func (pa PositionalArgument) Uint() *uint {
	v, err := strconv.Atoi(pa.value)
	if err != nil {
		return nil
	}

	i := uint(v)

	return &i
}

func (pa PositionalArgument) Uint8() *uint8 {
	v, err := strconv.Atoi(pa.value)
	if err != nil {
		return nil
	}

	i := uint8(v)

	return &i
}

func (pa PositionalArgument) Uint16() *uint16 {
	v, err := strconv.Atoi(pa.value)
	if err != nil {
		return nil
	}

	i := uint16(v)

	return &i
}

// Uint32 returns uint32 stored in uint64
func (pa PositionalArgument) Uint32() *uint64 {
	v, err := strconv.ParseUint(pa.value, 10, 64)
	if err != nil {
		return nil
	}

	return &v
}

func (pa PositionalArgument) Uint64() *uint64 {
	v, err := strconv.ParseUint(pa.value, 10, 64)
	if err != nil {
		return nil
	}

	return &v
}

// Float32 returns float32 stored in float64
func (pa PositionalArgument) Float32() *float64 {
	v, err := strconv.ParseFloat(pa.value, 32)
	if err != nil {
		return nil
	}

	return &v
}

func (pa PositionalArgument) Float64() *float64 {
	v, err := strconv.ParseFloat(pa.value, 64)
	if err != nil {
		return nil
	}

	return &v
}

func (pa PositionalArgument) Bool() *bool {
	b, err := strconv.ParseBool(pa.value)
	if err != nil {
		return nil
	}
	return &b
}
