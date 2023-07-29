package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":56998")
	// conn, err := net.Dial("tcp", "43.206.142.153:56998")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	for i := 50; i < 500; i++ {
		_, err1 := conn.Write([]byte(fmt.Sprintf("%d", i)))
		if err1 != nil {
			fmt.Println("Client write failed: ", err.Error())
		}
	}
}
