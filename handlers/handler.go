package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const (
	//MAX_WORKERS defines maximum number of goroutines spwaned to process incoming request
	MAX_WORKERS = 10
	// MAX_BODY_LENGTH constrains body length to 5mb
	MAX_BODY_LENGTH = 5 << 20
)

func PrimeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("no methods except POST allowed"))
		return
	}

	contentLength, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{error: content-length is required}"))
		return
	}

	if contentLength > MAX_BODY_LENGTH {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{error: body is too large}"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{error: invalid input: %s}", err.Error())))
		return
	}

	input, err := CheckAndTrim(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	nums, idx, err := ConvertToNums(input)

	if err != nil {
		var resp string

		if errors.Is(err, strconv.ErrRange) {
			resp = fmt.Sprintf("{error: \"the given input is invalid. Number on index %d caused overflow\"}", idx)
		} else {
			resp = fmt.Sprintf("{error: \"the given input is invalid. Element on index %d is not a number\"}", idx)
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(resp))
		return
	}

	response := ExtractPrimes(nums)

	res := CompileResponse(len(nums), response)

	w.Write([]byte(res))
}

func IsPrime(num int) bool {
	if num <= 1 {
		return false
	}

	for n := 2; n <= num/2; n++ {
		if num%n == 0 {
			return false
		}
	}
	return true
}
func CheckAndTrim(body []byte) (string, error) {
	trimmed := strings.Trim(string(body), " ")

	if trimmed[0] != '[' || trimmed[len(trimmed)-1] != ']' {
		return "", errors.New("invalid syntax")
	}

	trimmed = strings.Trim(trimmed, "[]")
	return trimmed, nil
}

func ConvertToNums(input string) ([]int, int, error) {
	var nums []int

	// not to spawn any worker until we make sure input is correct
	for idx, strNum := range strings.Split(input, ",") {
		strNum = strings.Trim(strNum, " ")
		num, err := strconv.Atoi(strNum)

		if err != nil {
			return nil, idx, err
		}
		nums = append(nums, num)
	}
	return nums, 0, nil
}

func ExtractPrimes(nums []int) map[int]bool {
	buffer := make(chan bool, MAX_WORKERS)
	wg := &sync.WaitGroup{}

	response := make(map[int]bool, len(nums))
	mu := sync.Mutex{}

	for numIdx, num := range nums {
		wg.Add(1)
		buffer <- true
		go func(index, number int) {
			defer wg.Done()

			mu.Lock()
			response[index] = IsPrime(number)
			mu.Unlock()
			<-buffer
		}(numIdx, num)
	}

	wg.Wait()

	return response
}

func CompileResponse(n int, response map[int]bool) string {
	var respBody []string
	for i := 0; i < n; i++ {
		if response[i] {
			respBody = append(respBody, "true")
		} else {
			respBody = append(respBody, "false")
		}
	}

	res := "[" + strings.Join(respBody, ",") + "]"
	return res
}
