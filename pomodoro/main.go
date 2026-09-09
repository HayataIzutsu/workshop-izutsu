package main

import (
	"log"
	"net/http"
)

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("static")))
	return mux
}

func main() {
	log.Println("pomodoro app listening on :8080")
	if err := http.ListenAndServe(":8080", newMux()); err != nil {
		log.Fatal(err)
	}
}
