package main

import (
	"flag"
	"fmt"
	"github.com/yushi/gother/handler"
	"log"
	"net/http"
)

const VERSION = "0.1.0"

func listHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
<html>
  <body>
    <div>
      <a href="/proc/mem">/proc/mem</a>
    </div>
    <div>
    <a href="/proc/load">/proc/load</a>
    </div>
    <div>
    <a href="/gh/notification">/gh/notification</a>?user=XXX&apiurl=YYY
    </div>
  </body>
</html>
`)
}

func startService(port int) {
	procHandler := new(handler.ProcHandler)
	procHandler.Start()
	githubHandler := new(handler.GithubHandler)
	ekikaraHandler := new(handler.EkikaraHandler)
	ekikaraHandler.Init()

	http.HandleFunc("/", listHandler)
	http.HandleFunc("/proc/mem", procHandler.HandleMemory)
	http.HandleFunc("/proc/load", procHandler.HandleLoadavg)
	http.HandleFunc("/gh/notification", githubHandler.HandleNotification)
	http.HandleFunc("/ekikara", ekikaraHandler.HandleEkikara)

	log.Printf("About to listen on %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	port := flag.Int("p", 9090, "listen port")
	printVersion := flag.Bool("v", false, "print version")
	flag.Parse()

	if *printVersion {
		fmt.Println(VERSION)
		return
	}

	startService(*port)
}
