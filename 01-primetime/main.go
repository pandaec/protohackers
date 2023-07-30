package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
)

const DEBUG_MODE = false

func main() {
	port := ":56998"
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("listen: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("listening on port ", port)

	var notPrimeResponse = &Response{
		Method: "isPrime",
		Prime:  false,
	}

	notPrimeRes, err := json.Marshal(notPrimeResponse)
	if err != nil {
		panic(err)
	}

	var primeResponse = &Response{
		Method: "isPrime",
		Prime:  true,
	}

	primeRes, err := json.Marshal(primeResponse)
	if err != nil {
		panic(err)
	}

	// fmt.Println("Pre-calculate prime cache")
	isPrime := IsPrime()
	// isPrime(100_000_001)
	// fmt.Println("End Pre-calculate prime cache")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("connection from ", conn.RemoteAddr())
		go handle(conn, isPrime, primeRes, notPrimeRes)
	}
}

type Request struct {
	Method *string `json:"method"`
	Number *int    `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func IsPrime() func(n int) bool {
	// var cache = make(map[int]bool)
	return func(n int) bool {
		// if ret, ok := cache[n]; ok {
		// 	return ret
		// }
		if n == 2 || n == 3 {
			return true
		}
		if n <= 1 {
			return false
		}
		if n%2 == 0 {
			return false
		}
		for i := 3; i <= int(math.Sqrt(float64(n))); i += 2 {
			if n%i == 0 {
				// cache[n] = false
				return false
			}
		}
		// cache[n] = true
		return true
	}
}

func handle(conn net.Conn, isPrime func(n int) bool, primeRes []byte, notPrimeRes []byte) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		in := scanner.Bytes()
		if DEBUG_MODE {
			fmt.Println(string(in))
		}

		req := Request{}
		if err := json.Unmarshal(in, &req); err != nil {
			if _, err := conn.Write(in); err != nil {
				if DEBUG_MODE {
					fmt.Printf("Write failed (malform)")
				}
			}
			return
		}
		if req.Method == nil || req.Number == nil || *req.Method != "isPrime" {
			if _, err := conn.Write(in); err != nil {
				if DEBUG_MODE {
					fmt.Printf("Write failed (malform)")
				}
				return
			}
		}

		var res []byte
		if isPrime(*req.Number) {
			res = primeRes
		} else {
			res = notPrimeRes
		}
		res = append(res, byte('\n'))
		if DEBUG_MODE {
			fmt.Println(string(res))
		}
		if _, err := conn.Write(res); err != nil {
			if DEBUG_MODE {
				fmt.Printf("Write failed (res)")
			}
			return
		}
	}
}
