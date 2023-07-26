package main

import "fmt"

func multi_proxy() {
    L := len(ADDR) - 1
    fmt.Println(L)
    for i := 0; i < L; i++ {
        c, _ := gen_proxy(ADDR[i + 1])
        go transfer(ADDR[i], c)
    }
    Socks5_Proxy(ADDR[L])
}