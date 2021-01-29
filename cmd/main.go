package main

import (
	"fmt"
	"net/http"

	ttr "github.com/saifulwebid/telegram-to-raindrop"
)

func main() {
	http.HandleFunc("/", ttr.CFHandler)

	fmt.Println("Listening on 127.0.0.1:8080...")
	http.ListenAndServe("127.0.0.1:8080", http.DefaultServeMux)
}
