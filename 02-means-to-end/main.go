package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"net"
	"os"
)

var debugMode bool

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

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("connection from ", conn.RemoteAddr())
		go handleConn(conn)
	}
}

type packet interface {
	process() (bool, error)
}

type insert struct {
	timestamp int
	price     int
}

type query struct {
	mintime int
	maxtime int
}

func (pkt insert) process() (bool, error) {

}

func (pkt query) process() (bool, error) {

}

func handleConn(conn net.Conn) {
	defer conn.Close()

	var b = make([]byte, 9)
	for {
		if _, err := io.ReadAtLeast(conn, b, 9); err != nil {
			return
		}

		pkt, err := parsePacket(b)
		if err != nil {
			return
		}
		fmt.Println(pkt)
	}
}

func parsePacket(pkt []byte) (packet, error) {
	// length check

	switch pkt[1] {
	case 'I':
		return insert{
			timestamp: 1,
			price:     2,
		}, nil
	case 'Q':
		return query{
			mintime: math.MinInt,
			maxtime: math.MaxInt,
		}, nil
	default:
		return nil, errors.New("Unsupported header")
	}
}
