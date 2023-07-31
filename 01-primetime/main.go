package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
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

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("connection from ", conn.RemoteAddr())
		go handle(conn, primeRes, notPrimeRes)
	}
}

type Request struct {
	Method string      `json:"method"`
	Number json.Number `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func IsPrime(n json.Number) bool {
	for _, c := range n {
		if c == '-' || c == '.' {
			return false
		}
	}

	// Large int handling
	// Assume they won't request with large integer as IsPrime() can't handle it effectivly anyway
	lastDigit, _ := strconv.Atoi(n[len(n)-1:].String())
	if len(n) > 1 && lastDigit%2 == 0 {
		return false
	}

	x, err := n.Int64()
	if err != nil {
		fmt.Println(err)
		return false
	}
	if x == 2 || x == 3 {
		return true
	}
	if x <= 1 {
		return false
	}
	if x%2 == 0 {
		return false
	}
	for i := int64(3); i <= int64(math.Sqrt(float64(x))); i += 2 {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func handle(conn net.Conn, primeRes []byte, notPrimeRes []byte) {
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
		if req.Method != "isPrime" {
			if _, err := conn.Write(in); err != nil {
				if DEBUG_MODE {
					fmt.Printf("Write failed (malform)")
				}
			}
			return
		}

		var res []byte
		if IsPrime(req.Number) {
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
