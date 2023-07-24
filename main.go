package main

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
    Socks5_Proxy("127.0.0.1:8080")
}
