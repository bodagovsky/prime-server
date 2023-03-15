package ops

import "math"

//IsPrime returns true for prime numbers, false otherwise 
//Prime numbers considered only real numbers greater or equal than 2 that have only two divisors: 1 and number itself
func IsPrime(num float64) bool {
	if num <= 1 {
		return false
	}

	if num != math.Floor(num) {
		return false
	}

	intNum := int(num)

	for n := 2; n <= intNum/2; n++ {
		if intNum%n == 0 {
			return false
		}
	}
	return true
}
