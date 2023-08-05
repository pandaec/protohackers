package main

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestIsPrime(t *testing.T) {
	// hex1 := "51"
	s := "49"
	data, _ := hex.DecodeString(s)
	fmt.Printf("% x", data)
}
