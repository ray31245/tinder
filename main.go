package main

import (
	"fmt"
	"log"
	"net/http"
)

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "welcome to my tinder")
}

func main() {
	http.HandleFunc("/", welcome)

	port := "8089"

	log.Println("Tinder match server is running on port " + port)
	http.ListenAndServe(":"+port, nil)
}
