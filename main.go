package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"io/ioutil"
)

type Response struct {
	Number     int      `json:"number"`
	IsPrime    bool     `json:"is_prime"`
	IsPerfect  bool     `json:"is_perfect"`
	Properties []string `json:"properties"`
	DigitSum   int      `json:"digit_sum"`
	FunFact    string   `json:"fun_fact"`
}

type ErrorResponse struct {
	Number string `json:"number"`
	Error  bool   `json:"error"`
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func isPerfect(n int) bool {
	sum := 1
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			if i == n/i {
				sum += i
			} else {
				sum += i + n/i
			}
		}
	}
	return sum == n && n != 1
}

func isArmstrong(n int) bool {
	sum, temp := 0, n
	digits := len(strconv.Itoa(n))
	for temp > 0 {
		digit := temp % 10
		sum += int(math.Pow(float64(digit), float64(digits)))
		temp /= 10
	}
	return sum == n
}

func digitSum(n int) int {
	sum := 0
	for n > 0 {
		sum += n % 10
		n /= 10
	}
	return sum
}

func fetchFunFact(n int) string {
	url := fmt.Sprintf("http://numbersapi.com/%d/math", n)
	resp, err := http.Get(url)
	if err != nil {
		return "Could not retrieve fun fact"
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return strings.TrimSpace(string(body))
}

func classifyNumber(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Handle CORS

	numStr := r.URL.Query().Get("number")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Number: numStr, Error: true})
		return
	}

	properties := []string{}
	if isArmstrong(num) {
		properties = append(properties, "armstrong")
	}
	if num%2 == 0 {
		properties = append(properties, "even")
	} else {
		properties = append(properties, "odd")
	}

	response := Response{
		Number:     num,
		IsPrime:    isPrime(num),
		IsPerfect:  isPerfect(num),
		Properties: properties,
		DigitSum:   digitSum(num),
		FunFact:    fetchFunFact(num),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/classify-number", classifyNumber)
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
