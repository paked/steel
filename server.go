package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	http.Handle("/", r)

	log.Println("Serving on localhost:8080")
	log.Println(http.ListenAndServe("localhost:8080", nil))
}
