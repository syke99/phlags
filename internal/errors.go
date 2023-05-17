package internal

import (
	"errors"
	"fmt"
)

var InvalidName = errors.New("invalid name provided")
var InvalidNameTooManyDashes = errors.New("invalid flag, too many \"-\"")

func NoPhlagsForSet(set string) error {
	return fmt.Errorf("no phlags set for PhlagSet: %s", set)
}
