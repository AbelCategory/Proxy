package main

import (
    "encoding/binary"
    "errors"
    "fmt"
    "io"
    "net"
    "time"
)

func Auth(client net.Conn) error {
    var A [260]byte
    n, err := client.Read(A[:])
    if err != nil {
        return err
    }
    ver, nmethods := A[0], A[1]
    if ver != 0x05 {
        return errors.New("wrong version")
    }
    if int(nmethods) + 2 != n {
        return errors.New("wrong format")
    }
    var B = [2]byte{0x05, 0x00}
    n, err = client.Write(B[:])
    if n != 2 {
        return errors.New("write error")
    }
    if err != nil {
        return errors.New("write error:" + err.Error())
    }
    return nil
}

func UDP_Proxy(client net.UDPConn, dest string) error{
    var data [270]byte
    n, addr, err := client.ReadFromUDP(data[:])
    if err != nil {
        return errors.New("UDP_Proxy_error" + err.Error())
    }
    if dest != addr.String(){
        return nil
    }
}

func Connect(client net.Conn) (net.Conn, error) {
    var A [270]byte
    n, err := client.Read(A[:])
    if err != nil {
        return nil, err
    }
    // fmt.Println(n)
    // fmt.Println(A)
    ver, cmd, atyp := A[0], A[1], A[3]
    if ver != 0x05 {
        return nil, errors.New("wrong version")
    }
    if cmd != 0x01 && cmd != 0x02 && cmd != 0x03 {
        var B = [](byte){0x05, 0x07, 0x00, 0x01, 0, 0, 0, 0, 0, 0}
        client.Write(B[:])
        return nil, errors.New("wrong cmd")
    }
    if A[2] != 0 {
        return nil, errors.New("wrong format")
    }
    var addr string
    switch atyp {
    case 0x01:
        if n != 10 {
            return nil, errors.New("wrong format")
        }
        addr = fmt.Sprintf("%d.%d.%d.%d", A[4], A[5], A[6], A[7])
    case 0x03:
        addrl := int(A[4])
        if n != 7 + addrl {
            return nil, errors.New("wrong format")
        }
        addr = string(A[5 : 5 + addrl])
    case 0x04:
        if n != 22 {
            return nil, errors.New("wrong format")
        }
        addr = fmt.Sprintf("%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x", 
                           A[4], A[5], A[6], A[7], A[8], A[9], A[10], A[11], 
                           A[12], A[13], A[14], A[15], A[16], A[17], A[18], A[19])
    default:
        var B = [](byte){0x05, 0x08, 0x00, 0x01, 0, 0, 0, 0, 0, 0}
        client.Write((B[:]))
        return nil, errors.New("invalid type")
    }
    port := binary.BigEndian.Uint16(A[n - 2 : n])
    dest := fmt.Sprintf("%s:%d", addr, port)
    var B []byte = make([]byte, 270)
    B[0] = 0x05
    B[1] = 0x00
    B[2] = 0x00
    B[3] = atyp
    fmt.Println(n)
    copy(B[4:n], A[4:n])
    fmt.Println(B[:n])
    // fmt.Println(dest)
    if cmd == 1 {
        fmt.Println("port=",port)
        conn, er := net.DialTimeout("tcp", dest, time.Second)
        if er != nil {
            B[1] = 0x05
            client.Write(B[:n])
            return nil, errors.New("tcp error" + er.Error())
        }
        // var B = [](byte){0x05, 0, 0, 0x01, 0, 0, 0, 0, 0, 0}
        _, err = client.Write(B[:n])
        if err != nil {
            conn.Close()
            return nil, errors.New("write error" + er.Error())
        }
        return conn, nil
    } else if cmd == 3 {
        conn, err := net.ListenUDP("udp", &net.UDPAddr{
            Port : 0x0786,
            IP : net.IPv4(0, 0, 0, 0),
        })
        if err != nil {
            return nil, errors.New("upd_listen_error" + err.Error())
        }
        var B = [](byte){0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x07, 0x86}
        _, err = client.Write(B[:])
        if err != nil {
            return nil, errors.New("write error" + err.Error())
        }
        err = UDP_Proxy(conn, dest)
        if err != nil {

        }
        return conn, nil
    }
    return nil, nil 
}

func request(client, dest net.Conn) {
    go io.Copy(dest, client)
    // fmt.Println("????")
    io.Copy(client, dest)
}

func process(client net.Conn) {
    err := Auth(client)
    if err != nil {
        fmt.Println("authentication_err: ", err)
        return
    }
    fmt.Println("authentication ok!")
    dest, err := Connect(client)
    if err != nil {
        fmt.Println("connect_err: ", err)
        return
    }
    fmt.Println("connection ok!")
    request(client, dest)
    defer dest.Close()
    defer client.Close()
}

func main() {
    conn, err := net.Listen("tcp", ":2333")
    if err != nil {
        fmt.Println("linsten_err: ", err)
        return
    }
    for {
        client, err := conn.Accept()
        if err != nil {
            fmt.Println("accept_err: ", err)
            continue
        }
        go process(client)
    }

}
