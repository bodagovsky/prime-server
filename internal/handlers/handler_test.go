package handlers

import (
	"testing"

	"primeServer/internal/errors"
	"primeServer/internal/helpers"
	ops "primeServer/pkg/math"

	"github.com/stretchr/testify/assert"
)

func TestWholesIsPrime(t *testing.T) {
	expectedOutput := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}

	var output []int

	for n := -1; n < 100; n++ {
		if ops.IsPrime(float64(n)) {
			output = append(output, n)
		}
	}
	assert.ElementsMatch(t, expectedOutput, output)
}

func TestFractionalIsPrime(t *testing.T) {
	expectedOutPut := []int{2, 3, 5}

	input := []float64{1, 2.0, 3, 4, 5, 5.5}

	var output = make([]int, 0, len(expectedOutPut))

	for _, num := range input {
		if ops.IsPrime(num) {
			output = append(output, int(num))
		}
	}
	assert.ElementsMatch(t, expectedOutPut, output)
}

func TestConvertToNums(t *testing.T) {

	cases := []struct {
		title          string
		input          []interface{}
		idx            int
		expectedOutput []float64
		err            error
	}{
		{
			title:          "success",
			input:          []interface{}{float64(1), float64(2), float64(3), float64(4)},
			idx:            0,
			err:            nil,
			expectedOutput: []float64{1, 2, 3, 4},
		},
		{
			title:          "nan_value",
			input:          []interface{}{float64(1), "nan", float64(3), float64(4)},
			idx:            1,
			err:            errors.ErrNotNumber,
			expectedOutput: nil,
		},
		{
			title:          "null_value",
			input:          []interface{}{float64(1), float64(2), float64(3), nil},
			idx:            3,
			err:            errors.ErrNotNumber,
			expectedOutput: nil,
		},
		{
			title:          "success_fractional",
			input:          []interface{}{1.0, 2.0, 3.0, 4.4, 5.5},
			idx:            0,
			err:            nil,
			expectedOutput: []float64{1, 2, 3, 4.4, 5.5},
		},
	}

	for _, test := range cases {
		t.Run(test.title, func(t *testing.T) {
			nums, idx, err := helpers.ConvertToNums(test.input)

			assert.ElementsMatch(t, test.expectedOutput, nums)
			assert.Equal(t, test.idx, idx)
			assert.ErrorIs(t, err, test.err)
		})

	}
}

func TestExtractPrimes(t *testing.T) {
	cases := []struct {
		title          string
		nums           []float64
		expectedOutput []bool
	}{
		{
			title:          "success",
			nums:           []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			expectedOutput: []bool{false, true, true, false, true, false, true, false, false, false, true, false, true, false, false},
		},
	}

	for _, test := range cases {
		t.Run(test.title, func(t *testing.T) {
			primes := extractPrimes(test.nums)
			seqPrimes := extractPrimesSequantially(test.nums)

			for i := 0; i < len(test.nums); i++ {
				assert.Equal(t, test.expectedOutput[i], primes[i])
				assert.Equal(t, test.expectedOutput[i], seqPrimes[i])
			}
		})
	}
}
