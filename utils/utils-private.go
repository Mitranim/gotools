package utils

// Private helpers for other utilities.

import ()

/********************************* Utilities *********************************/

// Converts a string to a base 10 integer. Scans the string until the first non-
// digit and ignores the rest.
func atoi(str string) (result int) {
	// Loop over the bytes, building the result, until we hit a non-digit.
	for i := 0; i < len(str) && str[i] >= '0' && str[i] <= '9'; i++ {
		if result == 0 {
			result = int(str[i]) - '0'
		} else {
			result = result*10 + int(str[i]) - '0'
		}
	}
	return
}

// Converts a positive integer into a same-looking string.
func itoa(num int) (result string) {
	for num > 0 {
		result = string('0'+num%10) + result
		num /= 10
	}
	return
}
