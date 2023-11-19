package mask

import (
	"errors"
	"regexp"
)

var (
	validatePattern      *regexp.Regexp
	ErrInvalidMaskFormat = errors.New("invalid mask format")
)

func init() {
	validatePattern = regexp.MustCompile(`(?m)^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`)
}

type Mask struct {
	value string
}

func New(mask string) (Mask, error) {
	if !isCorrectMask(mask) {
		return Mask{}, ErrInvalidMaskFormat
	}

	return Mask{value: mask}, nil
}

func (m Mask) String() string {
	return m.value
}

func isCorrectMask(mask string) bool {
	return validatePattern.MatchString(mask)
}
