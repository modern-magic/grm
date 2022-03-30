package internal

import "strings"

func isBoolFlag(arg string, flag string) bool {

	if strings.HasPrefix(arg, flag) {
		rem := arg[len(flag):]
		return len(rem) == 0 || rem[0] == '='
	}
	return false
}
