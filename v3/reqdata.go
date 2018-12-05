package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func reqDataHandler(w http.ResponseWriter, r *http.Request) {
	resp := reqservice.Stat()
	json, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}
