package handlers

import "strconv"

func parsePaginationParams(str string, defaultValue uint) uint {
	number, err := strconv.Atoi(str)
	if err != nil || number < int(defaultValue) {
		return defaultValue
	}

	return uint(number)
}
