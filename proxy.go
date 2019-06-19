package main

import(
        "log"
        "net/url"
        "net/http"
        "net/http/httputil"
)

func Handler(w http.ResponseWriter, r *http.Request) {
        remote, err := url.Parse("http://google.com")
        if err != nil {
                panic(err)
        }

        proxy := httputil.NewSingleHostReverseProxy(remote)
        http.HandleFunc("/dw", handlerwww(proxy))
        if err != nil {
                panic(err)
        }
}

func handlerwww(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
        return func(w http.ResponseWriter, r *http.Request) {
                log.Println(r.URL)
                w.Header().Set("X-Ben", "Rad")
                p.ServeHTTP(w, r)
        }
}