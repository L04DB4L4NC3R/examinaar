package controller

import "strconv"

func CheckPort(p1 string, p2 string) bool {
	a, err := strconv.Atoi(p1)

	if err != nil {
		return false
	}
	if a <= 1024 || a > 65535 {
		return false
	}

	a, err = strconv.Atoi(p2)

	if err != nil {
		return false
	}
	if a <= 1024 || a > 65535 {
		return false
	}

	return true
}
