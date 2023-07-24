// package main

// import (
//     "fmt"
//     "io"
//     "log"
//     "net/http"
//     "net/http/httputil"
//     "time"
// )

// type Proxy struct {
// 	sessions map[string][]*http.Request
// }

// func NewProxy() *Proxy {
// 	return &Proxy{
// 		sessions: make(map[string][]*http.Request),
// 	}
// }

// func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	sessionID := r.Header.Get("Session-ID")
// 	if sessionID == "" {
// 		sessionID = generateSessionID()
// 	}
// 	p.sessions[sessionID] = append(p.sessions[sessionID], r)

// 	if r.Method == "REPLAY" {
// 		p.ReplaySession(w, r, sessionID)
// 		return
// 	}

// 	r.Header.Set("Session-ID", sessionID)
// 	r.RequestURI = ""
// 	resp, err := http.DefaultTransport.RoundTrip(r)
// 	if err != nil {
// 		http.Error(w, "Error processing request", http.StatusInternalServerError)
// 		return
// 	}
// 	defer resp.Body.Close()
// 	copyHeader(w.Header(), resp.Header)
// 	w.WriteHeader(resp.StatusCode)
// 	io.Copy(w, resp.Body)
// }

// func (p *Proxy) ReplaySession(w http.ResponseWriter, r *http.Request, sessionID string) {
// 	requests, ok := p.sessions[sessionID]
// 	if !ok {
// 		http.Error(w, "Session not found", http.StatusNotFound)
// 		return
// 	}

// 	for _, req := range requests {
// 		reqDump, _ := httputil.DumpRequest(req, true)
// 		fmt.Fprintf(w, "Replaying request:\n%s\n", reqDump)

// 		resp, err := http.DefaultTransport.RoundTrip(req)
// 		if err != nil {
// 			http.Error(w, "Error processing request", http.StatusInternalServerError)
// 			return
// 		}
// 		defer resp.Body.Close()

// 		respDump, _ := httputil.DumpResponse(resp, true)
// 		fmt.Fprintf(w, "Response:\n%s\n", respDump)
// 	}
// }

// func generateSessionID() string {
// 	return fmt.Sprintf("%d", time.Now().UnixNano())
// }

// func copyHeader(dst, src http.Header) {
// 	for k, vv := range src {
// 		for _, v := range vv {
// 			dst.Add(k, v)
// 		}
// 	}
// }

// func main() {
// 	proxy := NewProxy()
// 	server := http.Server{
// 		Addr:    ":8080",
// 		Handler: proxy,
// 	}
// 	log.Fatal(server.ListenAndServe())
// }

package main

import (
    "fmt"
    "net"
)

func main() {
	ip := net.ParseIP("192.168.0.1")
	ip4 := net.IP.To4(ip)
	buf := make([]byte, 4)
	copy(buf[:], ip4[:])
	fmt.Println(buf)
}