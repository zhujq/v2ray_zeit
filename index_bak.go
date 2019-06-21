package main

import (
	"log"
	"fmt"
	"net/url"
	"net/http"
	"net/http/httputil"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome you,your info:\r\n")
	fmt.Fprintf(w, "METHOD:"+r.Method+"\r\n")
	fmt.Fprintf(w, "URL:\r\n")
	fmt.Fprintf(w, "PATH:"+r.URL.Path+"\r\n")
	fmt.Fprintf(w, "SCHEME:"+r.URL.Scheme+"\r\n")
	fmt.Fprintf(w, "HOST:"+r.URL.Host+"URL-End\r\n")
    fmt.Fprintf(w, "Proto:"+r.Proto+"\r\n")
	fmt.Fprintf(w, "HOST:"+r.Host+"\r\n")
	fmt.Fprintf(w, "RequestUrl:"+r.RequestURI+"\r\n")
	if r.URL.Path == "/dw"{
        remote, err := url.Parse("https://v2ray.14065567.now.sh/dw")
        if err != nil {
                panic(err)
        }

        proxy := httputil.NewSingleHostReverseProxy(remote)
        fmt.Fprintf(w, "Proxying...")
        http.HandleFunc("/dw", handlerwww(proxy))
        if err != nil {
                panic(err)
        }
       

	}
}

func handlerwww(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
        return func(w http.ResponseWriter, r *http.Request) {
                log.Println(r.URL)
                p.ServeHTTP(w, r)
        }
}