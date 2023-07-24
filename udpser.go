package main

import (
    "encoding/binary"
    "errors"
    "fmt"
    "net"
)

func UDP_Proxy(A, B *net.UDPConn, dest string) error {
    fmt.Println("dest:", dest)
    Abuf := make([]byte, 777)
    An, Aaddr, err := A.ReadFromUDP(Abuf)
    fmt.Println("????????")
    fmt.Println("Abuf:", Abuf)
    fmt.Println("Aaddr", Aaddr)
    var Bn int
    var addr string
    var port uint16
    if err != nil {
        return errors.New("UDP_get_err" + err.Error())
    }
    if Aaddr.String() == dest {
        fmt.Println("ok")
        atyp := Abuf[3]
        switch atyp {
        case 0x01:
            addr = net.IPv4(Abuf[4], Abuf[5], Abuf[6], Abuf[7]).String()
            port = binary.BigEndian.Uint16(Abuf[8 : 10])
        case 0x03:
            addrl := int(Abuf[4])
            addr = string(Abuf[5 : 5 + addrl])
            port = binary.BigEndian.Uint16(Abuf[5 + addrl : 7 + addrl])
        case 0x04:
            addr = "[" + net.IP(Abuf[4 : 20]).String() + "]"
            port = binary.BigEndian.Uint16(Abuf[20 : 22])
        default:
            return errors.New("wrong type")
        }
        add := fmt.Sprintf("%s.%d", addr, port)
        Baddr, err := net.ResolveUDPAddr("udp", add)
        if err != nil{
            return errors.New("resolve error" + err.Error())
        }
        Bbuf := make([]byte, 777)
        for {
            B.WriteToUDP(Abuf[0 : An], Baddr)
            Bn, Baddr, err = B.ReadFromUDP(Bbuf)
            if err != nil {
                fmt.Println("read error of B", err)
            }
            A.WriteToUDP(Bbuf[0 : Bn], Aaddr)
            An, Aaddr, err = A.ReadFromUDP(Abuf)
            if err != nil{
                fmt.Println("read error of A", err)
            }
        }
    }
    return nil
}