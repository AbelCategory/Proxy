package main

// import (
//     "fmt"
//     "log"
//     "net/http"
//     "net/http/httputil"
//     "time"
// )

// func http_proxy(w http.ResponseWriter, r *http.Request) {
//     reqDump, err := httputil.DumpRequest(r, true)
//     if err != nil {
//         log.Println("http_request_err:", err)
//         return
//     }
//     fmt.Println(string(reqDump))
//     client := &http.Client{
//         Timeout: 10 * time.Second,
//     }
//     res, err := client.Do(r)
//     if err != nil {
//         log.Println("client_err:", err)
//         return
//     }
//     defer res.Body.Close()
//     resDump, err := httputil.DumpResponse(res, true)
//     if err != nil {
//         log.Println("http_response_err:", err)
//         return
//     }
//     fmt.Println(string(resDump))
//     res.Write(w)
// }

// func main() {
//     server := &http.Server{
//         ReadTimeout: 10 * time.Second,
//         WriteTimeout: 10 * time.Second,
//         Addr: ":8080",
//         Handler: http.HandlerFunc(http_proxy),
//     }
//     go log.Fatal(server.ListenAndServe())
// }
