package common

import (
	"regexp"
	"strconv"
)

func CheckLuhn(s string) (bool, error) {
	ctr := 0

	lengthIsOdd := (len(s)-1)%2 != 0

	for i := range len(s) {
		n, err := strconv.Atoi(string(s[i]))
		if err != nil {
			return false, err
		}

		posIsOdd := i%2 != 0
		mustDouble := !lengthIsOdd && posIsOdd || lengthIsOdd && !posIsOdd

		if mustDouble {
			if n < 5 {
				ctr += n * 2
			} else {
				ctr += n*2 - 9
			}
		} else {
			ctr += n
		}
	}
	return ctr%10 == 0, nil
}

func CheckForAllDigits(s string) bool {
	pattern := `^[0-9]*$`

	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return false
	}

	return matched
}

func CheckOrderNumberFormat(number string) (bool, error) {
	if number == "" {
		return false, nil
	}

	if !CheckForAllDigits(number) {
		return false, nil
	}

	valid, err := CheckLuhn(number)

	if err != nil {
		return false, err
	}

	if !valid {
		return false, err
	}

	return true, nil
}
