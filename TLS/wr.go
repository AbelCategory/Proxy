package TLS

import (
    "crypto/tls"
    "fmt"
    "io"
    "log"
    "net"
    "net/http"
)

func copyHeader(w, h http.Header) {
    for keys, vals := range h {
        for _, val := range vals {
            w.Add(keys, val)
        }
    }
}

// func forward(dst, src net.Conn) (written int64, err error){
//     buf := make([]byte, 32 * 1024)
//     for {
//         n, err := src.Read(buf)
//         if n > 0 {
//             written, err := dst.Write(buf[:n])
//             if err != nil {
//                 return int64(written), err
//             }
//         }
//         if err != nil {
//             if err == EOF {
//                 break
//             }
//             return written, err
//         }
//     }
// }

func transmit(dest, client net.Conn) {
    defer dest.Close()
    defer client.Close()
    go io.Copy(client, dest)
    io.Copy(dest, client)
}

func handlehttps(w http.ResponseWriter, r *http.Request) {
    dest, err := tls.Dial("tcp", r.Host, &tls.Config{
        InsecureSkipVerify: true,
    })
    fmt.Println("1")
    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
    hij, ok := w.(http.Hijacker)
    if !ok {
        http.Error(w, "hijacking not supported", http.StatusInternalServerError)
        return
    }
    fmt.Println("2")
    client, _, err := hij.Hijack()
    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
    }
    fmt.Println("3")
    cli:= tls.Server(client, func() *tls.Config{
        cert, err := tls.LoadX509KeyPair("ca.crt", "ca.key")
        if err != nil {
            log.Panicln("Load_certification_err:", err)
        }
        return &tls.Config{
            Certificates: []tls.Certificate{cert},
            InsecureSkipVerify: true,
        }
    }())
    transmit(dest, cli)
}

func handlehttp(w http.ResponseWriter, r *http.Request) {
    res, err := http.DefaultTransport.RoundTrip(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
        return
    }
    defer res.Body.Close()
    copyHeader(w.Header(), res.Header)
    io.Copy(w, res.Body)
    w.WriteHeader(res.StatusCode)
}

func handle(w http.ResponseWriter, r *http.Request) {
    // fmt.Println("????")
    // log.Fatalln("???")
    if r.Method == http.MethodConnect {
        fmt.Println("https!")
        handlehttp(w, r)
    } else {
        fmt.Println("http!")
        handlehttps(w, r)
    }
}

func tls_hijacker() {
    server := &http.Server {
        Addr: ":2333",
        TLSConfig: func() *tls.Config {
            cert, err := tls.LoadX509KeyPair("ca.crt", "ca.key")
            if err != nil {
                log.Panicln("Load_certification_err:", err)
            }
            return &tls.Config{
                Certificates: []tls.Certificate{cert},
                InsecureSkipVerify: true,
            }
        }(),
        Handler: http.HandlerFunc(handle),
    }
    fmt.Println("ok")
    go log.Fatal(server.ListenAndServeTLS("ca.crt", "ca.key"))
}

