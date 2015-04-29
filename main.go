package main

import (
	"fmt"
	//	"log"
	"net/http"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	configFile := "./config.json"
	cfg := NewConfigFromFile(configFile)
	bk, err := NewBackendFromConfig(cfg.DB)
	if err != nil {
		panic(err)
	}
	mgr := NewIDServiceMgr(bk, cfg)
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "it works")
	})
	http.HandleFunc("/get", mgr.Get)
	http.HandleFunc("/addservice", mgr.AddService)
	http.ListenAndServe(":8080", nil)
	//	log.Fatal(http.ListenAndServe(":8080", nil))
}
