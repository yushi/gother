package main

import (
	"flag"
	"fmt"
	"github.com/yushi/gother/handler"
	"log"
	"net/http"
)

func hello_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Go!")
}

func main() {
	port := flag.Int("p", 9090, "listen port")
	flag.Parse()

	procHandler := new(handler.ProcHandler)
	procHandler.Start()
	http.HandleFunc("/hello", hello_handler)
	http.HandleFunc("/proc/mem", procHandler.HandleMemory)
	http.HandleFunc("/proc/load", procHandler.HandleLoadavg)

	log.Printf("About to listen on %d", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
