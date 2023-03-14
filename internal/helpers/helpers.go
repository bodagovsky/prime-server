package helpers

import "primeServer/internal/errors"


//ConvertToNums takes slice of interface values as input and converts them into slice of float64 values
//if there is non-float64 value in input, returns ErrNotNumber error and index of such value as idx variable
func ConvertToNums(input []interface{}) (nums []float64, idx int, err error) {
	nums = make([]float64, 0, len(input))

	// not to spawn any worker until we make sure input is correct
	for idx, maybeNum := range input {
		if num, ok := maybeNum.(float64); ok {
			nums = append(nums, num)
		} else {
			return nil, idx, errors.ErrNotNumber
		}

	}
	return nums, 0, nil
}

