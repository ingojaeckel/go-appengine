package main

import (
		"fmt"
		"net/http"
       )

type Hello struct{}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "root");
	});
	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "foo");
	});

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "bar");
	});

	http.ListenAndServe("localhost:4000", nil)
}
