package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func main() {
	port := ":56998"
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("listen: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("listening on port ", port)

	var failResponse = &Response{
		Method: "isPrime",
		Prime:  false,
	}

	failres, err := json.Marshal(failResponse)
	if err != nil {
		panic(err)
	}

	var successResponse = &Response{
		Method: "isPrime",
		Prime:  true,
	}

	sucessres, err := json.Marshal(successResponse)
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
		go handle(conn, sucessres, failres)
	}
}

type Request struct {
	Method string `json:"method"`
	Number int    `json:"prime"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func IsPrime(n int) bool {
	if n < 4 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	for i := 3; i < n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func handle(conn net.Conn, successRes []byte, failRes []byte) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		in := scanner.Bytes()

		req := Request{}
		if err := json.Unmarshal(in, &req); err != nil {
			if _, err := conn.Write(in); err != nil {
				fmt.Printf("Write failed (malform)")
			}
			return
		}

		var res []byte
		if !IsPrime(req.Number) {
			res = failRes
		} else {
			res = successRes
		}
		res = append(res, byte('\n'))
		fmt.Println(string(res))
		if _, err := conn.Write(res); err != nil {
			fmt.Printf("Write failed (res)")
			return
		}
	}
}
