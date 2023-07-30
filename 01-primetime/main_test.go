package main

import (
	"encoding/json"
	"testing"
)

// func TestIsPrime(t *testing.T) {
// 	primes := []int{2, 3, 5, 47, 211, 1117, 7741, 7879, 494933, 63480017}
// 	nonprimes := []int{-1, 0, 1, 8, 16, 69, 77, 5823, 6629, 7721, 459655}
// 	isPrime := IsPrime()
// 	for _, n := range primes {
// 		if !isPrime(n) {
// 			t.Fatal("Should be prime but return false", n)
// 		}
// 	}
// 	for _, n := range nonprimes {
// 		if isPrime(n) {
// 			t.Fatal("Should not be prime but return true", n)
// 		}
// 	}
// }

func TestJSON(t *testing.T) {
	// s := "{\"method\":\"isPrime\",\"prime\":false}"
	s := `{"method":"isPrime","number":99028461091604130867360308802317059320180412877139314398,"bignumber":true}`
	var req *Request

	if err := json.Unmarshal([]byte(s), &req); err != nil {
		panic(err)
	}

	if *req.Method != "isPrime" && req.Number != nil {
		t.Fatal("Unexpected unmarshal result")
	}
}
