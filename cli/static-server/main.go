package main

import (
	"log"
	"net/http"
)

const Addr = ":3000"

func main() {
	fs := http.FileServer(http.Dir("testdata/static"))
	http.Handle("/", fs)

	log.Println("Static server listen on " + Addr)
	err := http.ListenAndServe(Addr, nil)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
