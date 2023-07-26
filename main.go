package main

import (
    "flag"
    "fmt"
)

var is_TLS = false
var is_divide = true

var ADDR []string

func main() {
    // conn, err := net.Listen("tcp", ":8080")
    // if err != nil {
    //     fmt.Println("linsten_err:", err)
    //     return
    // }
    // for {
    //     client, err := conn.Accept()
    //     if err != nil {
    //         fmt.Println("accept_err:", err)
    //         continue
    //     }
    //     go process(client)
    // }
    var t string
    var port string
    flag.StringVar(&t, "type", "server", "get proxy type")
    flag.StringVar(&port, "port", "8080", "get proxy port")
    flag.BoolVar(&is_TLS, "TLS", false, "enable TLS hijacking")
    flag.Parse()
    fmt.Println(t)
    switch t {
    case "server": 
        Socks5_Proxy("127.0.0.1:" + port)
    case "client":
        flag.BoolVar(&is_divide, "div", false, "enable client rules proxy")
        flag.Parse()
        var cli string
        flag.StringVar(&cli, "client", "127.0.0.1:1926", "get client address")
        go Socks5_Proxy("127.0.0.1:" + port)
        c, _ := gen_proxy("127.0.0.1:8080")
        transfer(cli, c)
    case "multi":
        var addr string
        flag.StringVar(&addr, "addr", "127.0.0.1:1926|127.0.0.1:7777", "the multi proxy address")
        flag.Parse()
        addr = addr + "|"
        fmt.Println(addr)
        lst := 0
        n := len(addr)
        for i := 0; i < n; i++ {
            if addr[i] == '|' {
                ADDR = append(ADDR, addr[lst : i])
                lst = i + 1
            }
        }
        ADDR = append(ADDR, "127.0.0.1:" + port)
        multi_proxy()
    default:
        fmt.Println("type not supported")
    }
}
