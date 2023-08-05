package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
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
	process(map[int32]int32) (interface{}, error)
}

type insert struct {
	timestamp int32
	price     int32
}

type query struct {
	mintime int32
	maxtime int32
}

func (pkt insert) process(m map[int32]int32) (interface{}, error) {
	m[pkt.timestamp] = pkt.price
	return m, nil
}

func (pkt query) process(m map[int32]int32) (interface{}, error) {
	count, sum := int32(0), int32(0)
	for timestamp, price := range m {
		if pkt.mintime <= timestamp && timestamp >= pkt.maxtime {
			count += 1
			sum += price
		}
	}
	if count == 0 {
		return 0, nil
	}
	return sum / count, nil
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	pricedb := make(map[int32]int32)
	for {
		var b = make([]byte, 9)
		if _, err := io.ReadAtLeast(conn, b, 9); err != nil {
			return
		}

		pkt, err := parsePacket(b)
		if err != nil {
			if debugMode {
				fmt.Println("Error: ", err)
			}
			return
		}
		result, err := pkt.process(pricedb)
		if err != nil {
			if debugMode {
				fmt.Println("Error: ", err)
			}
			return
		}

		switch value := result.(type) {
		case int32:
			buf := new(bytes.Buffer)
			err := binary.Write(buf, binary.BigEndian, value)
			if err != nil {
				if debugMode {
					fmt.Printf("Write binary failed")
				}
			}
			if _, err := conn.Write(buf.Bytes()); err != nil {
				if debugMode {
					fmt.Printf("Write response failed")
				}
			}
		case map[int32]int32:
			pricedb = value
		default:
			if debugMode {
				fmt.Printf("Unrecognised response")
			}
		}
	}
}

type PacketStruct struct {
	Header byte
	P1     int32
	P2     int32
}

func parsePacket(pkt []byte) (packet, error) {
	if len(pkt) < 9 {
		return nil, errors.New("packet too small")
	}

	data := PacketStruct{}
	err := binary.Read(bytes.NewBuffer(pkt[:]), binary.BigEndian, &data)
	if err != nil {
		if debugMode {
			fmt.Println("Error: ", err)
		}
	}

	h := strconv.Itoa(int(data.Header))
	switch h {
	case "I":
		return insert{
			timestamp: data.P1,
			price:     data.P2,
		}, nil
	case "Q":
		return query{
			mintime: data.P1,
			maxtime: data.P2,
		}, nil
	default:
		return nil, errors.New("unsupported header")
	}
}
