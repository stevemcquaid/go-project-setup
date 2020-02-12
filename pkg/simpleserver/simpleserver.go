package simpleserver

import (
	"fmt"
	"net/http"
)

func Run() {
	http.HandleFunc("/", HelloHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		panic(err)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world \n")
}
