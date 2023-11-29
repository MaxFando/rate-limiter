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

type IP struct {
	value string
}

func New(ip string) (IP, error) {
	if !isCorrectIP(ip) {
		return IP{}, ErrInvalidIPFormat
	}

	return IP{value: ip}, nil
}

func (i IP) String() string {
	return i.value
}

func isCorrectIP(ip string) bool {
	return validatePattern.MatchString(ip)
}
