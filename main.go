package main

import (
	"log"
	"net/http"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	configFile := "./config.json"
	cfg := NewConfigFromFile(configFile)

	mgr := NewIDPoolMgr()
	mgr.LoadConfig(cfg)

	http.Handle("/", mgr)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
