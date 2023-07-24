package TLS

import (
    "crypto/tls"
    "fmt"
    "io"
    "net"
)

func handleTLS(lis net.Listener, addr string, host string) {
    conn, err := lis.Accept()
    if err != nil {
        fmt.Println("TLS! listen_err: ", err)
        return
    }
    cert, err := gen_cert(addr)
    if err != nil {
        fmt.Println("TLS! generate_cert_error: ", err)
        return
    }
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
        InsecureSkipVerify: true,
    }
    B, err := tls.Dial("tcp", host, tlsConfig)
    if err != nil {
        fmt.Println("TLS! dial_err: ", err)       
        return
    }
    A := tls.Server(conn, tlsConfig)
    err = A.Handshake()
    if err != nil {
        fmt.Println("TLS! server_err: ", err)
        return
    }
    defer B.Close()
    go io.Copy(A, B)
    io.Copy(B, A)
}