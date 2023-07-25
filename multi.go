package main

func multi_proxy() {
    L := len(ADDR) - 1
    Socks5_Proxy(ADDR[L])
    for i := 0; i < L; i++ {
        c, _ := gen_proxy(ADDR[i + 1])
        go transfer(ADDR[i], c)
    }
}