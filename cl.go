package main

import (
    "errors"
    "fmt"
    "io"
    "log"
    "net"
    "strconv"
)

func auth(ser net.Conn) error {
    au := authpac{}
    au.meth(0)
    _, err := ser.Write(au.data())
    if err != nil {
        return errors.New("auth_write_error: " + err.Error())
    }
    var data [100]byte
    n, err := io.ReadFull(ser, data[:2])
    if n != 2 {
        return errors.New("auth_wrong_len: " + err.Error())
    }
    if err != nil {
        return errors.New("auth_read_error: " + err.Error())
    }
    if data[0] != 0x05 {
        return errors.New("server_wrong_version")
    }
    if data[1] == 0xff {
        return errors.New("not_supported_auth")
    }
    return nil
}

func try_connect(ser net.Conn, host string, port uint16, cmd byte) (*net.UDPConn, error) {
    // data[0] = 0x05
    // data[1] = cmd
    // data[2] = 0
    data := append([](byte){0x05, cmd, 0}, resolve_addr(host, port)...)
    fmt.Println("data: ", data)
    _, err := ser.Write(data)
    if err != nil {
        return nil, errors.New("connect_write_error" + err.Error())
    }
    var A [777]byte
    n, err := io.ReadFull(ser, A[:4])
    if n < 4 {
        return nil, errors.New("connect_wrong_format")
    }
    if err != nil {
        return nil, errors.New("connect_read_error" + err.Error())
    }
    if A[0] != 0x05 {
        return nil, errors.New("connect_wrong_version")
    }
    if A[1] != 0x00 {
        return nil, errors.New("connect_error")
    }
    if cmd == 0x01 {
        // var buf [23333]byte
        _, err := io.ReadFull(ser, A[4 : 10])
        if err != nil {
            return nil, errors.New("connect_wrong_format")
        }
    } else if cmd == 0x03 {
        addr, err := tran_addr(A[3], ser)
        if err != nil {
            return nil, err
        }
        // fmt.Println("true_addr: ", addr)
        upd_addr, err := net.ResolveUDPAddr("udp", addr)
        fmt.Println("bind_addr: ", upd_addr) 
        if err != nil {
            return nil, errors.New("resolve_address_error:" + err.Error())
        }
        ser, err := net.DialUDP("udp", nil, upd_addr)
        if err != nil {
            ser.Close()
            return nil, err
        }
        return ser, nil
    }
    return nil, nil
}

type Proxy_Client struct {
    host string
    port uint16
}

func gen_proxy(addr string) (Proxy_Client, error){
    h, p, err := net.SplitHostPort(addr)
    if err != nil {
        return Proxy_Client{}, err
    }
    po, _ := strconv.Atoi(p)
    return Proxy_Client{
        host: h,
        port: uint16(po),
    }, nil
}

func Do(host string, port uint16) (net.Conn, error) {
    fmt.Println("ip:",fmt.Sprintf("%s:%d", host, port))
    conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
    if err != nil {
        return nil, err
    }
    err = auth(conn)
    if err != nil {
        conn.Close()
        return nil, err
    }
    return conn, nil
}

func (c *Proxy_Client) TCP_Proxy(host string, port uint16) (net.Conn, error) {
    conn, err := Do(c.host, c.port)
    fmt.Println("so??")
    if err != nil {
        return nil, err
    }
    fmt.Println("Do_ok")
    _, err = try_connect(conn, host, port, 0x01)
    return conn, err
}

func (c* Proxy_Client) UDP_Proxy(host string, port uint16) (*net.UDPConn, error) {
    conn ,err := Do(c.host, c.port)
    if err != nil {
        return nil, err
    }
    udp, err := try_connect(conn, host, port, 0x03)
    if err != nil {
        // fmt.Println(err, "?????!!!!")
        return nil, err
    }
    return udp, nil
}

func tcp() {
    client := Proxy_Client{
        host: "127.0.0.1",
        port: 8080,
    }
    proxy, err := client.TCP_Proxy("127.0.0.1", 2333)
    if err != nil {
        log.Panicln("error:", err)
    }
    if err == nil {
        fmt.Println(proxy)
        proxy.Write([]byte("smarthehe"))
        buf := make([]byte, 777)
        n, err := proxy.Read(buf[:])
        fmt.Println(n)
        fmt.Println(buf[:n])
        if err == nil {
            fmt.Println(string(buf[:n]))
        } else {
            log.Panicln("err: ", err)
        }
    }
}

func udp() {
    client := Proxy_Client {
        host: "127.0.0.1",
        port: 8080,
    }
    // fmt.Println("????")
    addr, err  := net.ResolveUDPAddr("udp", "127.0.0.1:2333")
    if err != nil {
        log.Panicln("resolve_address_error", addr)
    }
    // fmt.Println("QWQ!!!")
    fmt.Println("????")
    proxy, err := client.UDP_Proxy("127.0.0.1", 2333)
    fmt.Println("!!!!!")
    if err != nil {
        log.Panicln("error:", err)
    }
    fmt.Println("addr", addr)
    // go func() {
        buf := make([]byte, 777)
        fmt.Println("????")
        proxy.WriteToUDP([]byte("smarthehe"), addr)
        fmt.Println(addr)
        for {
            n, addr, err := proxy.ReadFromUDP(buf[:])
            if err != nil {
                log.Println("err!:", err)
                break
            }
            fmt.Println(addr, string(buf[:n]))
        }
    // }()
}