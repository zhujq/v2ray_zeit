package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Welcome you,your info:\r\n")
	fmt.Println(w, "METHOD:"+r.Method+"\r\n")
	fmt.Println(w, "URL:\r\n")
	fmt.Println(w, "PATH:"+r.URL.Path+"\r\n")
	fmt.Println(w, "SCHEME:"+r.URL.Scheme+"\r\n")
	fmt.Println(w, "HOST:"+r.URL.Host+"URL-End\r\n")
    fmt.Println(w, "Proto:"+r.Proto+"\r\n")
	fmt.Println(w, "HOST:"+r.Host+"\r\n")
	fmt.Println(w, "RequestUrl:"+r.RequestURI+"\r\n")
	if r.URL.Path == "/dw"{
        resp, err := http.Get("http://www.google.com")
        if err != nil {
            panic(err)
        }

        defer resp.Body.Close()
        fmt.Println(w, "Proxying...")
        body, err := ioutil.ReadAll(resp.Body)
    	if err != nil {
          	panic(err)
    	}
 
    	fmt.Println(string(body))    

	}
}
