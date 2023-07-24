package main

import (
    "fmt"
    "net"
)

func process(conn net.Conn) {
    defer conn.Close()
    var data [2333]byte
    for {
        n, err := conn.Read(data[:])
        if err != nil {
            fmt.Println("read_err", err)
            break
        }
        fmt.Println(string(data[:n]))
        fmt.Println(data[:n])
        conn.Write(data[:n])
    }
}

func tcp() {
    listen, err := net.Listen("tcp", ":2333")
    if err != nil {
        fmt.Println("listen_error:", err)
        return
    }
    fmt.Println("listen ok")
    for {
        conn, err := listen.Accept()
        if err != nil {
            fmt.Println("Accept failed:", err)
            continue
        }
        fmt.Println("accept ok")
        go process(conn)
    }
}

func udp() {
    conn, err := net.ListenUDP("udp", &net.UDPAddr{
        Port : 2333,
        IP : net.IPv4(0, 0, 0, 0),
    })
    if err != nil {
        fmt.Println("udp_listen_failed", err)
        return
    }
    for {
        var data [2333]byte
        n, addr, err := conn.ReadFromUDP(data[:])
        if err != nil {
            fmt.Println("udp_read_error", err)
            return
        }
        fmt.Println(string(data[:]))
        _, err = conn.WriteToUDP(data[: n], addr)
        if err != nil {
            fmt.Println("udp_write_error", err)
            return
        }
    }
}

func main() {
    udp()
    
}