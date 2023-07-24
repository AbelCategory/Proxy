package main

import (
    "crypto/tls"
    "fmt"
    "io"
    "log"
    "net"
)

func handleTLS(lis net.Listener, addr string, host string) {
    conn, err := lis.Accept()
    if err != nil {
        fmt.Println("TLS! listen_err: ", err)
        return
    }
    log.Println("listen_accept!!")
    rcert, err := gen_cert(addr)
    if err != nil {
        fmt.Println("TLS! generate_cert_error: ", err)
        return
    }
    log.Println("generate_certification_ok!!")
    RemoteConfig := &tls.Config{
        Certificates: []tls.Certificate{rcert},
    }
    // tlsConfig := &tls.Config{
    //     Certificates: []tls.Certificate{cert},
    //     InsecureSkipVerify: true,
    // }
    B, err := tls.Dial("tcp", host, RemoteConfig)
    if err != nil {
        fmt.Println("TLS! dial_err: ", err)       
        return
    }
    defer B.Close()
    log.Println("TLS_Dial OK!!")
    // fmt.Println(tlsConfig)
    lcert, err := gen_cert(addr)
    if err != nil {
        fmt.Println("TLS! generate_cert_error: ",err)
    }
    LocalConfig := &tls.Config{
        Certificates: []tls.Certificate{lcert},
        // InsecureSkipVerify: true,
    }
    fmt.Println(LocalConfig)
    A := tls.Server(conn, LocalConfig)
    err = A.Handshake()
    log.Println("?????>>>>")
    if err != nil {
        fmt.Println("TLS! server_err: ", err)
        return
    }
    log.Println("TLS_Serve OK!!")
    defer B.Close()
    go io.Copy(A, B)
    io.Copy(B, A)
}