package simpleserver

import (
	"fmt"
	"net/http"
)

func Run() {
	http.HandleFunc("/", HelloHandler)
	http.ListenAndServe("localhost:8080", nil)
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world \n")
}
