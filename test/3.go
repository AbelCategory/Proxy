package main

import (
    "fmt"
    "net"
)

func main() {
	cidr := "192.168.0.0/24"

	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Println("解析失败:", err)
		return
	}

	minIP := ip.Mask(ipNet.Mask)
	maxIP := make(net.IP, len(minIP))
	copy(maxIP, minIP)

	for i := range maxIP {
		maxIP[i] |= ^ipNet.Mask[i]
	}
	cc := ip.To4()
	fmt.Println(cc[0])

	fmt.Println("最小 IP 地址:", minIP)
	fmt.Println("最大 IP 地址:", maxIP)
}