package verbosity

import (
	"fmt"
	"strconv"
)

type Verbosity int

const (
	Quiet Verbosity = iota
	Minimal
	Normal
	Detailed
	Diagnostic
)

const (
	minVerbosity = Quiet
	maxVerbosity = Diagnostic
)

var oorError = fmt.Errorf("verbosity must be in range [%d, %d]", minVerbosity, maxVerbosity)

func FromInt(v int) (Verbosity, error) {
	if v < int(minVerbosity) || v > int(maxVerbosity) {
		return 0, oorError
	}
	return Verbosity(v), nil
}

func FromString(v string) (Verbosity, error) {
	vInt, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	return FromInt(vInt)
}
