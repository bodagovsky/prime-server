package test

import (
	"errors"
	"strconv"
	"testTask/handlers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPrime(t *testing.T) {
	expectedOutput := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}

	var output []int

	for n := -1; n < 100; n++ {
		if handlers.IsPrime(n) {
			output = append(output, n)
		}
	}
	assert.ElementsMatch(t, expectedOutput, output)
}

func TestCheckAndTrim(t *testing.T) {

	cases := map[string]struct {
		title          string
		expectedOutput string
		err            error
	}{
		"[1,2,3,4,5]": {
			title:          "case_1",
			expectedOutput: "1,2,3,4,5",
			err:            nil,
		},
		"[1,2,3,4,5": {
			title:          "case_2",
			expectedOutput: "",
			err:            errors.New("invalid syntax"),
		},
		" [ 1, 2 ,3 ,4 ]": {
			title:          "case_3",
			expectedOutput: " 1, 2 ,3 ,4 ",
			err:            nil,
		},
	}

	for input, test := range cases {
		t.Run(test.title, func(t *testing.T) {
			out, err := handlers.CheckAndTrim([]byte(input))
			assert.Equal(t, test.expectedOutput, out)
			assert.Equal(t, test.err, err)
		})
	}

}

func TestConvertToNums(t *testing.T) {

	cases := []struct {
		title          string
		input          string
		idx            int
		expectedOutput []int
		err            error
	}{
		{
			title:          "case_1",
			input:          " 1, 2 ,3 ,4 ",
			idx:            0,
			expectedOutput: []int{1, 2, 3, 4},
		},
		{
			title:          "case_2",
			input:          "1,2,3,4",
			idx:            0,
			expectedOutput: []int{1, 2, 3, 4},
		},
		{
			title:          "case_3",
			input:          " 1, \"nan\" ,3 ,4 ",
			idx:            1,
			err:            strconv.ErrSyntax,
			expectedOutput: nil,
		},
		{
			title:          "case_4",
			input:          " 1, 2 ,3 ,9999999999999999999999999 ",
			idx:            3,
			err:            strconv.ErrRange,
			expectedOutput: nil,
		},
	}

	for _, test := range cases {
		t.Run(test.title, func(t *testing.T) {
			nums, idx, err := handlers.ConvertToNums(test.input)

			assert.ElementsMatch(t, test.expectedOutput, nums)
			assert.Equal(t, test.idx, idx)
			assert.True(t, errors.Is(err, test.err))
		})

	}
}

func TestExtractPrimes(t *testing.T) {
	cases := []struct {
		title          string
		nums           []int
		expectedOutput map[int]bool
	}{
		{
			title: "success",
			nums:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			expectedOutput: map[int]bool{
				0:  false,
				1:  true,
				2:  true,
				3:  false,
				4:  true,
				5:  false,
				6:  true,
				7:  false,
				8:  false,
				9:  false,
				10: true,
				11: false,
				12: true,
				13: false,
				14: false,
			},
		},
	}

	for _, test := range cases {
		t.Run(test.title, func(t *testing.T) {
			primes := handlers.ExtractPrimes(test.nums)

			for i := 0; i < len(test.nums); i++ {
				assert.Equal(t, test.expectedOutput[i], primes[i])
			}
		})
	}
}

func TestCompileResponse(t *testing.T) {
	cases := []struct{
		title string
		n int
		expectedOutput string
		response map[int]bool
	}{
		{
			title: "case_1",
			n: 15,
			response: map[int]bool{
				0:  false,
				1:  true,
				2:  true,
				3:  false,
				4:  true,
				5:  false,
				6:  true,
				7:  false,
				8:  false,
				9:  false,
				10: true,
				11: false,
				12: true,
				13: false,
				14: false,
			},
			expectedOutput: "[false,true,true,false,true,false,true,false,false,false,true,false,true,false,false]",
		},
	}
	
	
	for _, test := range cases {
		t.Run(test.title, func(t *testing.T) {
			assert.Equal(t, test.expectedOutput, handlers.CompileResponse(test.n, test.response))
		})
	}
}
