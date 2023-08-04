package main

import (
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
    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("accept: ", err.Error())
            os.Exit(1)
        }
        fmt.Println("connection from ", conn.RemoteAddr())
        go handle(conn)
    }
}

func handle(conn net.Conn) {
    defer conn.Close()

    if _, err := io.Copy(conn, conn); err != nil {
        fmt.Println("copy: ", err.Error())
    }
}
