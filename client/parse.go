package client

import (
    "encoding/binary"
    "errors"
    "fmt"
    "io"
    "net"
)

type authpac struct {
    nmethods byte
    methods []byte
}

func (t* authpac) data() []byte{
    return append([]byte{0x05, t.nmethods}, t.methods...)
}

func (t* authpac) meth(method... byte) {
    t.methods = method
    t.nmethods = byte(len(method))
}

func resolve_addr(host string, port uint16) []byte{
    if host[0] >= '0' && host[0] <= '9' {
        ip := net.ParseIP(host).To4()
        return binary.BigEndian.AppendUint16(append([]byte{0x01}, ip...), port)

    } else if host[0] == '[' {
        ip := net.ParseIP(host)
        return binary.BigEndian.AppendUint16(append([]byte{0x04}, ip...), port)
    } else {
        str := []byte(host)
        return binary.BigEndian.AppendUint16(append([]byte{0x03, byte(len(str))}, str...), port)
    }
}

func tran_addr(atyp byte, ser net.Conn) (string, error) {
    var A [777]byte
    var addr string
    var port uint16
    switch atyp {
    case 0x01:
        n, err := io.ReadFull(ser, A[0 : 6])
        if n != 6 {
            return "", errors.New("wrong format")
        }
        if err != nil {
            return "", errors.New("wrong_read" + err.Error())
        }
        addr = net.IPv4(A[0], A[1], A[2], A[3]).String()
        port = binary.BigEndian.Uint16(A[4 : 6])
    case 0x03:
        _, err := io.ReadFull(ser, A[0 : 1])
        addrl := int(A[4])
        if err != nil {
            return "", errors.New("wrong_read" + err.Error())
        }
        n, err := io.ReadFull(ser, A[0 : 2 + addrl])
        if err != nil {
            return "", errors.New("wrong_read" + err.Error())
        }
        if n != 2 + addrl {
            return "", errors.New("wrong format")
        }
        addr = string(A[0 : 2 + addrl])
        port = binary.BigEndian.Uint16(A[addrl : addrl + 2])
    case 0x04:
        n, err := io.ReadFull(ser, A[0 : 18])
        if n != 18 {
            return "", errors.New("wrong format")
        }
        if err != nil {
            return "", errors.New("wrong_read" + err.Error())
        }
        addr = "[" + net.IP(A[0 : 16]).String() + "]"
        port = binary.BigEndian.Uint16(A[16 : 18])
    default:
        return "", errors.New("connect_bind_error")
    }
    return fmt.Sprintf("%s:%d", addr, port), nil
}