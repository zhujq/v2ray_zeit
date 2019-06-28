package main

import (
	"fmt"
	"net/http"
)


func Handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Welcome you,your info:\r\n")
	return
}
