package main

import (
	"errors"
	"fmt"
	"io"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":56998")
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	result := make([]byte, 0)
	for {
		buf := make([]byte, 128)
		n, err := conn.Read(buf)
		if errors.Is(err, io.EOF) {
			fmt.Println("IO EOF", err)
			result = append(result, buf[:n]...)
			_, err1 := conn.Write(result)
			if err1 != nil {
				println("Server write failed: ", err.Error())
				return
			}
			fmt.Println("[W]", string(result))
			return
		}

		if err != nil {
			println("Server read failed: ", err.Error())
			return
		}

		result = append(result, buf[:n]...)

		fmt.Println("[R]", string(buf[:n]))
	}

}
