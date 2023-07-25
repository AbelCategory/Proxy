package main

import (
    "bufio"
    "bytes"
    "net"
    "os"
    "strings"
)

func CIDR_resolve(cidr string) ([]byte, []byte, error) {
    ip, ipNet, err := net.ParseCIDR(cidr)
    if err != nil {
        return nil, nil, err
    }
    minIP := ip.Mask(ipNet.Mask)
    maxIP := make(net.IP, len(minIP))
    copy(maxIP, minIP)
    for i := range maxIP {
        maxIP[i] |= ^ipNet.Mask[i]
    }
    return minIP, maxIP, nil
}

func check_ipv4(A []byte) bool {
    // if A[0] == 1 && A[1] == 0 && A[2] >= 1 && A[2] <= 3 {
    //     return true
    // }
    // if A[0] == 1 && A[1] == 0 && A[2] >=8 && A[2] <= 15 {
    //     return true
    // }
    // return false
    file, err := os.Open("divide/ipv4.txt")
    if err != nil{
        return false
    }
    defer file.Close()
    IP := net.IPv4(A[0], A[1], A[2], A[3])
    sc := bufio.NewScanner(file)
    for sc.Scan() {
        s := sc.Text()
        // ones, _ := ipNet.Mask.Size()
        // l := binary.LittleEndian.Uint32(ipd[: 4])
        // l = l & (cc << ones)
        // r := l | (cc >> (32 - ones))
        // if l <= ip && ip <= r {
        //     return true
        // }
        minIP, maxIP, err := CIDR_resolve(s)
        if err != nil {
            continue
        }
        if bytes.Compare(IP, minIP) >= 0 && bytes.Compare(IP, maxIP) <= 0{
            return true
        }
    }
    return false
}

func check_ipv6(A []byte) bool {
    file, err := os.Open("divide/ipv6.txt")
    if err != nil {
        return false
    }
    defer file.Close()
    IP := net.IP(A[0 : 16])
    sc := bufio.NewScanner(file)
    for sc.Scan() {
        s := sc.Text()
        minIP, maxIP, err := CIDR_resolve(s)
        if err != nil {
            continue
        }
        if bytes.Compare(IP, minIP) >= 0 && bytes.Compare(IP, maxIP) <= 0 {
            return true
        }
    }
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