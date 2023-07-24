package main

import (
    "encoding/binary"
    "errors"
    "fmt"
    "io"
    "net"
    "time"
    "strings"
)

func Auth(client net.Conn) error {
    var A [260]byte
    n, err := io.ReadFull(client, A[:2])
    if n != 2 {
        return errors.New("wrong format")
    }
    if err != nil {
        return err
    }
    ver, nmethods := A[0], A[1]
    n, err = io.ReadFull(client, A[2 : 2 + nmethods])
    if ver != 0x05 {
        return errors.New("wrong version")
    }
    if n != int(nmethods) {
        return errors.New("wrong version")
    }
    fmt.Println("auth:", A)
    var B = [2]byte{0x05, 0xff}
    for i := 2; i < 2 + int(nmethods); i++ {
        if A[i] == 0 {
            B[1] = 0x00
            break
        }
    }
    n, err = client.Write(B[:])
    if n != 2 {
        return errors.New("write error")
    }
    if err != nil {
        return errors.New("write error:" + err.Error())
    }
    if B[1] == 0xff {
        return errors.New("auth error")
    } else {
        return nil
    }
}

func Connect(client net.Conn) (net.Conn, error) {
    var A [270]byte
    n, err := io.ReadFull(client, A[: 4])
    // fmt.Println(n, A)
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
        n, err = io.ReadFull(client, A[4 : 10])
        if n != 6 {
            return nil, errors.New("wrong format")
        }
        if err != nil {
            return nil, errors.New("wrong format" + err.Error())
        }
        addr = fmt.Sprintf("%d.%d.%d.%d", A[4], A[5], A[6], A[7])
        n = 10
    case 0x03:
        n, err = io.ReadFull(client, A[4 : 5])
        addrl := int(A[4])
        if n != 1 || err != nil {
            return nil, errors.New("errrrrr")
        }
        n, err = io.ReadFull(client, A[5 : 5 + addrl + 2])
        if n != addrl + 2{
            return nil, errors.New("wrong format")
        }
        if err != nil {
            return nil, errors.New("wrong format" + err.Error())
        }
        addr = string(A[5 : 5 + addrl])
        n = 5 + addrl + 2
    case 0x04:
        n, err = io.ReadFull(client, A[4 : 22])
        if n != 18 {
            return nil, errors.New("wrong format")
        }
        if err != nil {
            return nil, errors.New("wrong format" + err.Error())
        }
        addr = fmt.Sprintf("[%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x]", 
                           A[4], A[5], A[6], A[7], A[8], A[9], A[10], A[11], 
                           A[12], A[13], A[14], A[15], A[16], A[17], A[18], A[19])
        n = 22
    default:
        var B = [](byte){0x05, 0x08, 0x00, 0x01, 0, 0, 0, 0, 0, 0}
        _, _ = client.Write(B[:])
        return nil, errors.New("invalid type")
    }
    fmt.Println(A)
    port := binary.BigEndian.Uint16(A[n - 2 : n])
    fmt.Println(port)
    dest := fmt.Sprintf("%s:%d", addr, port)
    fmt.Println(dest)
    if cmd == 1 {
        conn, er := net.DialTimeout("tcp", dest, time.Second)
        if er != nil {
            var B = [](byte){0x05, 0x03, 0, 0x01, 0, 0, 0, 0, 0, 0}
            // fmt.Println("???:", er)
            // fmt.Println(strings.Contains(er.Error(), "connection"))
            if strings.Contains(er.Error(), "connection") {
                B[1] = 0x05
            } else if strings.Contains(er.Error(), "lookup") {
                B[1] = 0x04
            }
            _, _ = client.Write(B[:])
            return nil, errors.New("tcp error" + er.Error())
        }
        var B = [](byte){0x05, 0, 0, 0x01, 0, 0, 0, 0, 0, 0}
        n, err = client.Write(B[:])
        if err != nil {
            conn.Close()
            return nil, errors.New("write error" + er.Error())
        }
        return conn, nil
    }
    return nil, nil
}

func request(client, dest net.Conn) {
    go io.Copy(client, dest)
    io.Copy(dest, client)
    fmt.Println("????")
}

func process(client net.Conn) {
    err := Auth(client)
    if err != nil {
        fmt.Println("authentication_err:", err)
        return
    }
    // fmt.Println("authentication ok!")
    dest, err := Connect(client)
    if err != nil {
        fmt.Println("connect_err:", err)
        return
    }
    fmt.Println("connection ok!")
    request(client, dest)
    defer dest.Close()
    defer client.Close()
}

func main() {
    conn, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("linsten_err:", err)
        return
    }
    for {
        client, err := conn.Accept()
        if err != nil {
            fmt.Println("accept_err:", err)
            continue
        }
        go process(client)
    }

}
