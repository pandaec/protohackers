package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"net"
	"os"
	"regexp"
	"strconv"
)

var debugMode = false

func main() {
	portArg := flag.Int("port", 56998, "port")
	debugModeArg := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	port := *portArg
	debugMode = *debugModeArg

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("listen: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("listening on port ", port)

	notPrimeRes, err := json.Marshal(&Response{
		Method: "isPrime",
		Prime:  false,
	})
	if err != nil {
		panic(err)
	}

	primeRes, err := json.Marshal(&Response{
		Method: "isPrime",
		Prime:  true,
	})
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
	Method string          `json:"method"`
	Number json.RawMessage `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func IsPrime(n string) bool {
	for _, c := range n {
		if c == '-' || c == '.' {
			return false
		}
	}

	// Large int handling
	// Assume they won't request with large integer as IsPrime() can't handle it effectivly anyway
	lastDigit, _ := strconv.Atoi(n[len(n)-1:])
	if len(n) > 1 && lastDigit%2 == 0 {
		return false
	}

	x, _ := strconv.ParseInt(n, 10, 64)
	if x == 2 || x == 3 {
		return true
	}
	if x <= 1 || x%2 == 0 || x%3 == 0 {
		return false
	}
	for i := int64(5); i <= int64(math.Sqrt(float64(x))); i += 2 {
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
		if debugMode {
			fmt.Println(string(in))
		}

		req := Request{}
		if err := json.Unmarshal(in, &req); err != nil {
			if _, err := conn.Write([]byte("🦆")); err != nil {
				if debugMode {
					fmt.Printf("Write failed (malform)")
				}
			}
			return
		}
		numberStr := string(req.Number)
		if req.Method != "isPrime" || numberStr == "" {
			if _, err := conn.Write([]byte("🦆")); err != nil {
				if debugMode {
					fmt.Printf("Write failed (malform)")
				}
			}
			return
		}

		numericPattern, _ := regexp.Compile(`^[+-]?\d+(?:\.\d+)?$`)
		if !numericPattern.MatchString(numberStr) {
			if _, err := conn.Write([]byte("🦆")); err != nil {
				if debugMode {
					fmt.Printf("Write failed (malform)")
				}
			}
		}

		var res []byte
		if IsPrime(numberStr) {
			res = primeRes
		} else {
			res = notPrimeRes
		}
		res = append(res, byte('\n'))
		if debugMode {
			fmt.Println(string(res))
		}
		if _, err := conn.Write(res); err != nil {
			if debugMode {
				fmt.Printf("Write failed (res)")
			}
			return
		}
	}
}
