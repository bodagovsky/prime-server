package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"primeServer/internal/errors"
	"primeServer/internal/helpers"
	ops "primeServer/internal/math"
)

const (
	//MAX_WORKERS defines maximum number of goroutines spwaned to process incoming request
	MAX_WORKERS = 10

	// MAX_BODY_LENGTH constrains body length to 5mb
	MAX_BODY_LENGTH = 5 << 20

	//BIG_INPUT const defines the minimum size if an input to calculate response in parallel for. For input less than BIG_INPUT calculate response sequantially
	BIG_INPUT = 100
)

type PrimeIncomingRequest []interface{}
type PrimeIncomingResponse []bool

type PrimeIncomingErrorResponse struct {
	Error string `json:"error"`
}

//PrimeHandler takes a slice of values and responds with slice with boolean values where each boolean denotes if the input number is prime
//if non number input is given, such as [1,2,"nan"] or [1,2,3,4, null] responds with error message
func PrimeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(PrimeIncomingErrorResponse{
			Error: errors.ErrWrongMethod.Error(),
		})
		return
	}

	var header int
	var response interface{}

	defer func() {
		w.Header().Add("Content-Type", "application/json")
		if header != 0 {
			w.WriteHeader(header)
		}
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(PrimeIncomingErrorResponse{
				Error: errors.ErrInternalServer.Error(),
			})
			return
		}
	}()

	contentLength, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		header = http.StatusBadRequest
		response = PrimeIncomingErrorResponse{
			Error: errors.ErrContentLengthNotProvided.Error(),
		}
		return
	}

	if contentLength > MAX_BODY_LENGTH {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PrimeIncomingErrorResponse{
			Error: errors.ErrBodyTooLarge.Error(),
		})
		return
	}

	var reqIn PrimeIncomingRequest

	err = json.NewDecoder(r.Body).Decode(&reqIn)
	if err != nil {
		header = http.StatusBadRequest
		response = PrimeIncomingErrorResponse{
			Error: errors.ErrWrongBodySyntax.Error(),
		}
		return
	}

	nums, idx, err := helpers.ConvertToNums(reqIn)

	if err != nil {
		header = http.StatusBadRequest
		response = PrimeIncomingErrorResponse{
			Error: fmt.Sprintf("the given input is invalid. Element on index %d is not a number", idx),
		}
		return
	}

	var resp []bool
	if len(nums) >= BIG_INPUT {
		//run in parallel
		resp = extractPrimes(nums)
	} else {
		resp = extractPrimesSequantially(nums)
	}

	// err = json.NewEncoder(w).Encode(resp)
	response = resp
}

//extractPrimesSequantially returns slice of boolean results of IsPrime func for each number in nums
func extractPrimesSequantially(nums []float64) []bool {
	response := make([]bool, len(nums))

	for idx, num := range nums {
		response[idx] = ops.IsPrime(num)
	}

	return response
}

//extractPrimes does exactly the same as extractPrimesSequantially, but in parallel spawning MAX_WORKERS goroutines
func extractPrimes(nums []float64) []bool {
	buffer := make(chan bool, MAX_WORKERS)
	wg := &sync.WaitGroup{}

	response := make([]bool, len(nums))
	mu := sync.Mutex{}

	for numIdx, num := range nums {
		wg.Add(1)
		buffer <- true
		go func(index int, number float64) {
			defer wg.Done()

			mu.Lock()
			response[index] = ops.IsPrime(number)
			mu.Unlock()
			<-buffer
		}(numIdx, num)
	}

	wg.Wait()

	return response
}
