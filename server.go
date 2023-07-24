package main

import (
    "bufio"
    "encoding/binary"
    "errors"
    "fmt"
    "io"
    "net"
    "strings"
    "time"
)

var is_TLS = true

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
    if err != nil {
        return err
    }
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

func request(client, dest net.Conn) {
    go io.Copy(client, dest)
    io.Copy(dest, client)
    fmt.Println("????")
}

func Connect(client net.Conn) error {
    var A [270]byte
    n, err := io.ReadFull(client, A[: 4])
    // fmt.Println(n, A)
    if err != nil {
        return err
    }
    // fmt.Println(n)
    // fmt.Println(A)
    ver, cmd, atyp := A[0], A[1], A[3]
    if ver != 0x05 {
        return errors.New("wrong version")
    }
    if cmd != 0x01 && cmd != 0x02 && cmd != 0x03 {
        var B = [](byte){0x05, 0x07, 0x00, 0x01, 0, 0, 0, 0, 0, 0}
        client.Write(B[:])
        return errors.New("wrong cmd")
    }
    if A[2] != 0 {
        return errors.New("wrong format")
    }
    var addr string
    switch atyp {
    case 0x01:
        n, err = io.ReadFull(client, A[4 : 10])
        if n != 6 {
            return errors.New("wrong format")
        }
        if err != nil {
            return errors.New("wrong format" + err.Error())
        }
        addr = net.IPv4(A[4], A[5], A[6], A[7]).String()
        n = 10
    case 0x03:
        n, err = io.ReadFull(client, A[4 : 5])
        addrl := int(A[4])
        if n != 1 || err != nil {
            return errors.New("errrrrr")
        }
        n, err = io.ReadFull(client, A[5 : 5 + addrl + 2])
        if n != addrl + 2{
            return errors.New("wrong format")
        }
        if err != nil {
            return errors.New("wrong format" + err.Error())
        }
        addr = string(A[5 : 5 + addrl])
        n = 5 + addrl + 2
    case 0x04:
        n, err = io.ReadFull(client, A[4 : 22])
        if n != 18 {
            return errors.New("wrong format")
        }
        if err != nil {
            return errors.New("wrong format" + err.Error())
        }
        // addr = fmt.Sprintf("[%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x]", 
        //                    A[4], A[5], A[6], A[7], A[8], A[9], A[10], A[11], 
        //                    A[12], A[13], A[14], A[15], A[16], A[17], A[18], A[19])
        addr = "[" + net.IP(A[4 : 20]).String() + "]"
        n = 22
    default:
        var B = [](byte){0x05, 0x08, 0x00, 0x01, 0, 0, 0, 0, 0, 0}
        _, _ = client.Write(B[:])
        return errors.New("invalid type")
    }
    fmt.Println()
    fmt.Println(A)
    port := binary.BigEndian.Uint16(A[n - 2 : n])
    fmt.Println(port)
    dest := fmt.Sprintf("%s:%d", addr, port)
    fmt.Println(dest)
    if cmd == 1 {
        conn, er := net.DialTimeout("tcp", dest, time.Second * 3)
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
            return errors.New("tcp error" + er.Error())
        }
        var B = [](byte){0x05, 0, 0, 0x01, 0, 0, 0, 0, 0, 0}
        _, err = client.Write(B[:])
        if err != nil {
            conn.Close()
            return errors.New("write error" + er.Error())
        }
        
        if is_TLS {
            defer client.Close()
            rd := bufio.NewReader(client)
            buf, err := rd.Peek(8)
            fmt.Println("buf:", buf)
            if err != nil {
                return errors.New("read_client_err:" + err.Error())
            }
            if buf[0] == 0x16 && buf[1] == 0x03 && buf[2] == 0x01 {
                conn.Close()
                lis, err := net.Listen("tcp", ":0")
                if err != nil {
                    return errors.New("listener_err_:" + err.Error())
                }
                Addr := lis.Addr().String()
                fmt.Println("listen TLS:", Addr)
                go handleTLS(lis, addr, dest)
                prox, err := net.DialTimeout("tcp", Addr, time.Second * 3)
                if err != nil {
                    return errors.New("t..dial_error: " + err.Error())
                }
                defer prox.Close()
                var buf[32 * 1024]byte
                go io.Copy(client, prox)
                n, err := rd.Read(buf[:])
                if err != nil {
                    return err
                }
                prox.Write(buf[:n])
                request(client, prox)
            }
        } else{
            request(client, conn)
        }
        // conn.Write([]byte("???"))
        return nil
    } else if cmd == 3 {
        // fmt.Println("??????!!!!")
        adr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
        if err != nil {
            return errors.New("resolve_adr_err:" + err.Error())
        }
        A, err := net.ListenUDP("udp", adr)
        if err != nil {
            return errors.New("listen_udp_error:" + err.Error())
        }
        B, err := net.ListenUDP("udp", nil)
        if err != nil {
            return errors.New("listen_udp_error:" + err.Error())
        }
        ad := A.LocalAddr().(*net.UDPAddr)
        fmt.Println("bind_addr:", ad)
        // var buf[777]byte = [](byte){0x05, 0x00, 0x00, 0x01, 0x00,}
        buf := make([]byte, 333)
        buf[0] = 0x05
        buf[3] = 0x01
        copy(buf[4 : 8], ad.IP)
        // fmt.Println(ad.IP)
        binary.BigEndian.PutUint16(buf[8 : 10], uint16(ad.Port))
        fmt.Println("buf:", buf)
        _, err = client.Write(buf[: 22])

        if err != nil {
            return errors.New("write_client_err:" + err.Error())
        }
        defer A.Close()
        defer B.Close()
        fmt.Println("????")
        err = UDP_Proxy(A, B, dest)
        return err
    }
    return nil
}

func process(client net.Conn) {
    fmt.Println("ok")
    err := Auth(client)
    if err != nil {
        fmt.Println("authentication_err:", err)
        return
    }
    fmt.Println("authentication ok!")
    err = Connect(client)
    if err != nil {
        fmt.Println("connect_err:", err)
        return
    }
    fmt.Println("connection ok!")
    // request(client, dest)
    // defer dest.Close()
    // defer client.Close()
}

func Socks5_Proxy(s string) error {
    conn, err := net.Listen("tcp", s)
    fmt.Println("errrrrr")
    if err != nil {
        // fmt.Println("linsten_err:", err)
        return errors.New("listen_err: " + err.Error())
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