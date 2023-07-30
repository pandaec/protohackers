package main

import "testing"

func TestIsPrime(t *testing.T) {
	primes := []int{2, 3, 5, 47, 211, 1117, 7741, 7879, 494933}
	nonprimes := []int{8, 16, 69, 77, 5823, 6629, 7721}

	for _, n := range primes {
		if !IsPrime(n) {
			t.Fatal("Should be prime but return false", n)
		}
	}
	for _, n := range nonprimes {
		if IsPrime(n) {
			t.Fatal("Should not be prime but return true", n)
		}
	}
}
