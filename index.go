package main

import (
	"fmt"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!\r\n")
	fmt.Fprintf(w, "PATH:"+r.URL.Path+"\r\n")
	fmt.Fprintf(w, "SCHEME:"+r.URL.Scheme+"\r\n")
	fmt.Fprintf(w, "METHOD:"+r.Method+"\r\n")
	fmt.Fprintf(w, "HOST:"+r.Host+"\r\n")

}