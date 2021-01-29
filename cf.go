package ttr

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func CFHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("failed reading request body: %v", err)
	}

	fmt.Println(r)
	fmt.Println(string(body))
}
