package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"testing"
)

func TestIsPrime(t *testing.T) {
	primes := []int{2, 3, 5, 47, 211, 1117, 7741, 7879, 494933, 63480017}
	nonprimes := []int{-1, 0, 1, 8, 16, 69, 77, 5823, 6629, 7721, 459655}
	for _, x := range primes {
		n := strconv.Itoa(x)
		if !IsPrime(n) {
			t.Fatal("Should be prime but return false", n)
		}
	}
	for _, x := range nonprimes {
		n := strconv.Itoa(x)
		if IsPrime(n) {
			t.Fatal("Should not be prime but return true", n)
		}
	}
}

func TestValidNumberString(t *testing.T) {
	validNumberString := func(s string) bool {
		r, _ := regexp.Compile(`^[+-]?\d+(?:\.\d+)?$`)
		return r.MatchString(s)
	}

	valids := []string{"123", "0", "-1", "-2312", "16.33344"}
	invalids := []string{"", "1_123", "1,000", `"9999"`}
	for _, s := range valids {
		if !validNumberString(s) {
			t.Fatal("Should be valid but return false", s)
		}
	}
	for _, s := range invalids {
		if validNumberString(s) {
			t.Fatal("Should be invalid but return true", s)
		}
	}
}

func TestJSON(t *testing.T) {
	// s := "{\"method\":\"isPrime\",\"prime\":false}"
	s := `{"method":"isPrime","number":99028461091604130867360308802317059320180412877139314398,"bignumber":true}`

	type reqs struct {
		Method string      `json:"method"`
		Number json.Number `json:"number"`
	}

	var req *reqs

	if err := json.Unmarshal([]byte(s), &req); err != nil {
		panic(err)
	}

	if req.Method != "isPrime" {
		t.Fatal("Unexpected unmarshal result")
	}
}

func TestJSONStringNum(t *testing.T) {
	s := `{"method":"isPrime","number":"7474024", "number1":7474024}`

	type reqs struct {
		Method  string          `json:"method"`
		Number  json.RawMessage `json:"number"`
		Number1 json.RawMessage `json:"number1"`
	}

	var req reqs

	if err := json.Unmarshal([]byte(s), &req); err != nil {
		panic(err)
	}
	v := string(req.Number)
	v1 := string(req.Number1)
	fmt.Println(v, v1)
	if req.Method != "isPrime" {
		t.Fatal("Unexpected unmarshal result")
	}
}
