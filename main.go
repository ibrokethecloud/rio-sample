package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var bgColor string
	var ok bool

	if bgColor, ok = os.LookupEnv("COLOR"); !ok {
		bgColor = "powderblue"
	}

	fmt.Fprintf(w, "<html><body style=\"background-color:%s;\"><h1>Rio Demo</h1><p>Sample application to show case Rio code workflows</p></body></html>", bgColor)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
