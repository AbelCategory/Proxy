package main

import (
    "fmt"
    "io"
    "net"
)

func forward(client, dest net.Conn) {
    go io.Copy(client, dest)
    io.Copy(dest, client)
    fmt.Println("????")
}

func client_process(conn net.Conn, c Proxy_Client) {
    buf1, err := forward_auth(conn)
    if err != nil {
        fmt.Println("wrong_auth")
        return
    }
    err, dest := forward_conn(conn, buf1[:], c)
    fmt.Println(dest)
    if err != nil {
        fmt.Println("wrong_conn")
        return
    }
    forward(conn, dest)
}

func transfer(addr string, c Proxy_Client) {
    listen, err := net.Listen("tcp", addr)
    if err != nil {
        fmt.Println("listen_error: ", err)
        return
    }
    for {
        conn, err := listen.Accept()
        if err != nil {
            fmt.Println("err: ", err)
            continue
        }
        go client_process(conn, c)
    }
}

// func main() {
//     transfer()
// }

// func main() {
//     conn, err := net.Listen("tcp", ":2333")
//     if err != nil {
//         log.Fatal("listen_err:", err)
//     }
//     defer conn.Close()
//     for {
//         client, err := conn.Accept()
//         if err != nil {
//             fmt.Println("acc_err:", err)
//             continue
//         }
//         // go Do()
//     }
// }