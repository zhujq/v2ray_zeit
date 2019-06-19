package main

import(
        "fmt"
        "log"
        "net/url"
        "net/http"
        "net/http/httputil"
        "raymain"
)

func main() {
        int status :=0

        remote, err := url.Parse("http://v2ray.14065567.now.sh/dw")
        if err != nil {
                panic(err)
        }

        proxy := httputil.NewSingleHostReverseProxy(remote)
        http.HandleFunc("/dw", handler(proxy))
        http.HandleFunc("/",handlerwww)
        
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
        return func(w http.ResponseWriter, r *http.Request) {
                log.Println(r.URL)
                p.ServeHTTP(w, r)
        }
}

func handlerwww(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "text/html")
        fmt.Fprintf(w, "<br /><h3>hello world </h3>")
}