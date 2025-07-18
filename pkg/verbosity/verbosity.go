package verbosity

import (
	"errors"
	"strings"
)

type Verbosity int

const (
	Quiet Verbosity = iota
	Minimal
	Normal
	Detailed
	Diagnostic
)

func FromString(s string) (v Verbosity, err error) {
	switch strings.ToLower(s) {
	case "quiet", "0":
		v = Quiet
	case "minimal", "1":
		v = Minimal
	case "normal", "2":
		v = Normal
	case "detailed", "3":
		v = Detailed
	case "diagnostic", "4":
		v = Diagnostic
	default:
		return v, errors.New(
			"invalid format: allowed formats are 'quiet' (0), 'minimal' (1), " +
				"'normal' (2), 'detailed' (3), 'diagnostic' (4)")
	}
	return v, nil
}

func (t *Verbosity) UnmarshalText(data []byte) (err error) {
	*t, err = FromString(string(data))
	return err
}
