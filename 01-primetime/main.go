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

	res, err := json.Marshal(failResponse)
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
		go handle(conn, res)
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

func isPrime(n int) bool {
	return true
}

func handle(conn net.Conn, failResponse []byte) {
	defer conn.Close()

	buf := bytes.NewBuffer([]byte{})
	io.Copy(buf, conn)
	fmt.Println(buf.String())

	req := request{}
	if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
		if _, err := conn.Write(failResponse); err != nil {
			fmt.Printf("Write failed")
		}
		return
	}

	fmt.Println(req)

	if !isPrime(req.Number) {
		if _, err := conn.Write(failResponse); err != nil {
			fmt.Printf("Write failed")
		}
		return
	}

	// io.Copy(conn, buf)
}
