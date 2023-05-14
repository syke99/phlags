package phlags

import "strconv"

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
