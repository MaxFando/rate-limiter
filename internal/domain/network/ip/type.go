package ip

import (
	"errors"
	"regexp"
)

var (
	validatePattern    *regexp.Regexp
	ErrInvalidIPFormat = errors.New("invalid ip format")
)

func init() {
	validatePattern = regexp.MustCompile(`(?m)^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`)
}

type Ip struct {
	value string
}

func New(ip string) (Ip, error) {
	if !isCorrectIP(ip) {
		return Ip{}, ErrInvalidIPFormat
	}

	return Ip{value: ip}, nil
}

func (i Ip) String() string {
	return i.value
}

func isCorrectIP(ip string) bool {
	return validatePattern.MatchString(ip)
}
