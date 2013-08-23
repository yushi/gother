package main

import (
	"flag"
	"fmt"
	"github.com/yushi/gother/statusboard"
	"github.com/yushi/gother/system"
	"log"
	"net/http"
	"time"
)

func hello_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Go!")
}

func getTimeStr() string {
	// for test
	//return time.Now().Format("15:04:05")

	return time.Now().Format("15:04")
}

func get_proc_load_handler() func(w http.ResponseWriter, r *http.Request) {
	stats := make(map[string]*system.SystemStat)

	return func(w http.ResponseWriter, r *http.Request) {
		stats[getTimeStr()] = system.GetSystemStat()
		fmt.Fprintf(w, "%s", statusboard.LoadavgGraph(stats))
	}
}

func get_proc_mem_handler() func(w http.ResponseWriter, r *http.Request) {
	stats := make(map[string]*system.SystemStat)

	return func(w http.ResponseWriter, r *http.Request) {
		stats[getTimeStr()] = system.GetSystemStat()
		fmt.Fprintf(w, "%s", statusboard.MemoryGraph(stats))
	}
}

func main() {
	port := flag.Int("p", 9090, "listen port")
	flag.Parse()

	http.HandleFunc("/hello", hello_handler)
	http.HandleFunc("/proc/mem", get_proc_mem_handler())
	http.HandleFunc("/proc/load", get_proc_load_handler())

	log.Printf("About to listen on %d", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
