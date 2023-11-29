package utils

import (
	"strconv"
	"strings"
)

func GetPrefix(inputIP string, inputMask string) (string, error) {
	ip := strings.Split(inputIP, ".")
	mask := strings.Split(inputMask, ".")
	var prefix string
	for index, ipOct := range ip {
		intIPOct, err := strconv.Atoi(ipOct)
		if err != nil {
			return "", err
		}
		intMaskOct, err := strconv.Atoi(mask[index])
		if err != nil {
			return "", err
		}
		prefix += strconv.Itoa(intIPOct & intMaskOct)
		if index != len(ip)-1 {
			prefix += "."
		}
	}
	return prefix, nil
}
