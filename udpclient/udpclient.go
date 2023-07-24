package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)

func tcp() {
    conn, err := net.Dial("tcp", ":2333")
    if err != nil {
        fmt.Println("err : ", err)
        return
    }
    defer conn.Close() // 关闭TCP连接
    inputReader := bufio.NewReader(os.Stdin)
    for {
        input, _ := inputReader.ReadString('\n') // 读取用户输入
        inputInfo := strings.Trim(input, "\n")
        fmt.Println(inputInfo)
        if strings.ToUpper(inputInfo) == "Q" { // 如果输入q就退出
            return
        }
        _, err := conn.Write([]byte(inputInfo)) // 发送数据
        if err != nil {
            return
        }
        buf := [512]byte{}
        n, err := conn.Read(buf[:])
        if err != nil {
            fmt.Println("recv failed, err:", err)
            return
        }
        fmt.Println(string(buf[:n]))
    }
}

func udp() {
    client, err := net.DialUDP("udp", nil, &net.UDPAddr{
        Port: 2333, 
        IP : net.IPv4(0, 0, 0, 0),
    }) 
    if err != nil {
        fmt.Println("listen_udp_error", err)
        return
    }

    client.Write([]byte("233"))
    var a = make([]byte, 270)
    n, _, _ := client.ReadFromUDP(a[:])
    fmt.Println(string(a[:n]))
}

func main() {
    udp()
}