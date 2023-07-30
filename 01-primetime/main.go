package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	var failResponse = &response{
		Method: "isPrime",
		Prime:  false,
	}

	failres, err := json.Marshal(failResponse)
	if err != nil {
		panic(err)
	}

	var successResponse = &response{
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

type request struct {
	Method string
	Number int
}

type response struct {
	Method string
	Prime  bool
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

func handle(conn net.Conn, success []byte, failres []byte) {
	defer conn.Close()

	for {
		buf := bytes.NewBuffer([]byte{})
		io.Copy(buf, conn)
		fmt.Println(buf.String())

		req := request{}
		if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
			if _, err := conn.Write(buf.Bytes()); err != nil {
				fmt.Printf("Write failed")
			}
			return
		}

		fmt.Println(req)

		if !IsPrime(req.Number) {
			if _, err := conn.Write(failres); err != nil {
				fmt.Printf("Write failed (fail res)")
				return
			}
		}

		if _, err := conn.Write(success); err != nil {
			fmt.Printf("Write failed (success res)")
			return
		}
	}
}
