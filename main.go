package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, here is hello-go</h1>")
	} else if r.URL.Path == "/about" {
		fmt.Fprint(w, "This repository for recording learnning "+
			"<a href=\"mailto:albert@example.com\">albert@example.com</a>")
	} else {
		fmt.Fprint(w, "<h1>not found</h1>")
	}
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil)
}
