package main

import (
	"fmt"
	"net/http"
	"strings"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()  // parse arguments, you have to call this by yourself
	fmt.Println("form:", r.Form)  // print form information in server side
	fmt.Println("path:", r.URL.Path)
	fmt.Println("scheme:", r.URL.Scheme)
	fmt.Println("url_long:", r.Form["url_long"])

	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	fmt.Fprintf(w, "returning string\n") // send data to client side
}

func main() {
	http.HandleFunc("/", sayhelloName) // set router

	err := http.ListenAndServe(":8080", nil) // set listen port
	if err != nil {
		panic(err)
	}
	fmt.Println("Server is running on port 8080")
}
