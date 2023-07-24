package client

import (
    "encoding/binary"
    "errors"
    "fmt"
    "io"
    "net"
    "strings"
)

func check_ipv4(A []byte) bool {
    if A[0] == 1 && A[1] == 0 && A[2] >= 1 && A[2] <= 3 {
        return true
    }
    if A[0] == 1 && A[1] == 0 && A[2] >=8 && A[2] <= 15 {
        return true
    }
    return false
}

func check_ipv6(A []byte) bool {
    return false
}

func check_domain(addr string) bool {
    if strings.Contains(addr, ".cn") {
        return true
    }
    if strings.Contains(addr, "baidu.com") {
        return true
    }
    if strings.Contains(addr, "zhihu.com") {
        return true
    }
    return false
}

func forward_auth(client net.Conn) ([]byte, error) {
    var A [260]byte
    n, err := io.ReadFull(client, A[:2])
    if n != 2 {
        return nil, errors.New("wrong format")
    }
    if err != nil {
        return nil, err
    }
    ver, nmethods := A[0], A[1]
    n, err = io.ReadFull(client, A[2 : 2 + nmethods])
    if err != nil {
        return nil, err
    }
    if ver != 0x05 {
        return nil, errors.New("wrong version")
    }
    if n != int(nmethods) {
        return nil, errors.New("wrong version")
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
        return nil, errors.New("write error")
    }
    if err != nil {
        return nil, errors.New("write error:" + err.Error())
    }
    if B[1] == 0xff {
        return nil, errors.New("auth error")
    } else {
        return A[: 2 + int(nmethods)], nil
    }
}

func forward_conn(client net.Conn, buf []byte) (error, net.Conn) {
    var A [270]byte
    n, err := io.ReadFull(client, A[: 4])
    // fmt.Println(n, A)
    if err != nil {
        return err, nil
    }
    // fmt.Println(n)
    // fmt.Println(A)
    ver, cmd, atyp := A[0], A[1], A[3]
    if ver != 0x05 {
        return errors.New("wrong version"), nil
    }
    if cmd != 0x01 && cmd != 0x02 && cmd != 0x03 {
        var B = [](byte){0x05, 0x07, 0x00, 0x01, 0, 0, 0, 0, 0, 0}
        client.Write(B[:])
        return errors.New("wrong cmd"), nil
    }
    if A[2] != 0 {
        return errors.New("wrong format"), nil
    }
    var addr string
    var dest net.Conn
    switch atyp {
    case 0x01:
        n, err = io.ReadFull(client, A[4 : 10])
        if n != 6 {
            return errors.New("wrong format"), nil
        }
        if err != nil {
            return errors.New("wrong format" + err.Error()), nil
        }
        addr = net.IPv4(A[4], A[5], A[6], A[7]).String()
        n = 10
    case 0x03:
        n, err = io.ReadFull(client, A[4 : 5])
        addrl := int(A[4])
        if n != 1 || err != nil {
            return errors.New("errrrrr"), nil
        }
        n, err = io.ReadFull(client, A[5 : 5 + addrl + 2])
        if n != addrl + 2{
            return errors.New("wrong format"), nil
        }
        if err != nil {
            return errors.New("wrong format" + err.Error()), nil
        }
        addr = string(A[5 : 5 + addrl])
        n = 5 + addrl + 2
    case 0x04:
        n, err = io.ReadFull(client, A[4 : 22])
        if n != 18 {
            return errors.New("wrong format"), nil
        }
        if err != nil {
            return errors.New("wrong format" + err.Error()), nil
        }
        // addr = fmt.Sprintf("[%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x]", 
        //                    A[4], A[5], A[6], A[7], A[8], A[9], A[10], A[11], 
        //                    A[12], A[13], A[14], A[15], A[16], A[17], A[18], A[19])
        addr = "[" + net.IP(A[4 : 20]).String() + "]"
        n = 22
    default:
        var B = [](byte){0x05, 0x08, 0x00, 0x01, 0, 0, 0, 0, 0, 0}
        _, _ = client.Write(B[:])
        return errors.New("invalid type"), nil
    }
    fmt.Println(A)
    port := binary.BigEndian.Uint16(A[n - 2 : n])
    add := fmt.Sprintf("%s:%d", addr, port)
    fmt.Printf("dest_address: %s:%d\n", addr, port)
    if cmd == 0x03 {
        return nil, nil
    }
    // fmt.Println("??????")
    if atyp == 0x01 {
        if check_ipv4(A[4 : 8]) {
            dest, err = net.Dial("tcp", add)
            if err != nil {
                return errors.New("dial_error" + err.Error()), nil
            }
        }
    } else if atyp == 0x03 {
        if check_domain(addr) {
            dest, err = net.Dial("tcp", add)
            if err != nil {
                return errors.New("dial_error" + err.Error()), nil
            }
        }
    } else {
        if check_ipv6(A[4 : 20]) {
            dest, err = net.Dial("tcp", add)
            if err != nil {
                return errors.New("dial_error" + err.Error()), nil
            }
        }
    }
    if dest != nil {
        fmt.Println("??????")
        var data = [](byte){0x05, 0, 0, 0x01, 0, 0, 0, 0, 0, 0}
        client.Write(data[:])
        return nil, dest
    } else {
        fmt.Println("ok")
        c := Proxy_Client{
            host: "127.0.0.1",
            port: 8080,
        }
        if cmd == 0x01 {
            proxy, err := c.TCP_Proxy(addr, port)
            fmt.Println("?????:", err)
            if err != nil {
                var data = [](byte){0x05, 0, 0x01, 0x01, 0, 0, 0, 0, 0, 0}
                client.Write(data[:])
                // proxy.Close()
                return errors.New("proxy_connect_err:" + err.Error()), nil
            }
            fmt.Println("suuu")
            var data = [](byte){0x05, 0, 0, 0x01, 0, 0, 0, 0, 0, 0}
            client.Write(data[:])
            return nil, proxy 
        } else {
            var data = [](byte){0x05, 0, 0x01, 0x01, 0, 0, 0, 0, 0, 0}
            client.Write(data[:])
        }
        return nil, dest
    }
    // return nil, dest
}