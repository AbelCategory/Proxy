package main

import (
    "fmt"
    "net"
)

func main(){
    // addr := net.UDPAddr{
    //     IP : net.IPv4(1,2,3,4),
    //     Port : 8080,
    // }
    // buf := make([]byte, 64)
    // copy(buf[0:4], addr.IP)
    // fmt.Println(addr.IP, addr.Port)
    // fmt.Println(buf)
    // A, err := net.ListenUDP("udp", nil)
    // if err != nil {
    //     fmt.Println(err)
    //     return
    // }
    // ProxyAddr := A.LocalAddr().(*net.UDPAddr)
    // buf := make([]byte, 64)
    // copy(buf[:], ProxyAddr.IP)
    // fmt.Println(ProxyAddr.IP, ProxyAddr.Port)
    // fmt.Println(buf)
    // ipBytes := []byte{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 13} // 要分配的IPv6地址，以16个字节表示
	// port := 8080                                                        // 要分配的端口号

	// ip := net.IP(ipBytes[1 : 17])
	// if ip == nil {
	// 	fmt.Println("Invalid IP address")
	// 	return
	// }
    // fmt.Println(ip)

	// addr := &net.UDPAddr{
	// 	IP:   ip,
	// 	Port: port,
	// }

	// fmt.Println("UDP Address:", addr)
    // buf := make([]byte, 64)
    // copy(buf[:], addr.IP)
    // fmt.Println(buf)
    buf := []byte{1,2,3,4}
    ip := net.IPv4(buf[0], buf[1], buf[2], buf[3])
    cmp := make([]byte, 16)
    copy(cmp, ip)
    fmt.Println(ip)
    fmt.Println(cmp)
    jsl := make([]byte, 16)
    ip4 := ip.To4()
    copy(jsl, ip4)
    fmt.Println(jsl)
}