package main

import (
	"fmt"
	"log"
	//	"log"
	"net/http"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

const (
	ERR_ID = 0
)

type IDServiceMap map[string]*IDServiceGenerator

type IDServiceMgr struct {
	//	stats       *Stats
	data   IDServiceMap
	bk     *Backend
	locker *sync.RWMutex // protects data field
}

func NewIDServiceMgr(bk *Backend, cfg *Config) *IDServiceMgr {
	mgr := &IDServiceMgr{
		//		stats:       &Stats{},
		data:   make(IDServiceMap),
		bk:     bk,
		locker: new(sync.RWMutex),
	}
	log.Printf("Start service")
	err := mgr.loadAllService()
	if err != nil {
		panic(err)
	}
	return mgr
}

func (self *IDServiceMgr) RegisterIDService(ider *IDServiceGenerator) {
	self.data[ider.Service] = ider
}

// 获取新的id
func (self *IDServiceMgr) NewId(service string) uint64 {

	ider, ok := self.data[service]
	if !ok {
		return ERR_ID
	}
	return ider.Gen()
}

func (self *IDServiceMgr) loadAllService() error {
	iders, err := self.bk.GetAllService()
	if err != nil {
		return err
	}

	self.locker.Lock()
	defer self.locker.Unlock()

	for ider := range iders {
		self.RegisterIDService(ider)
	}
	return nil
}

// 获取新的id
func (self *IDServiceMgr) AddService(service string) error {
	self.locker.Lock()
	defer self.locker.Unlock()
	if _, ok := self.data[service]; !ok {
		ider, err := self.bk.CreateService(service)
		if err != nil {
			return err
		}
		self.RegisterIDService(ider)
	}
	return nil
}

func (self *IDServiceMgr) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	switch r.URL.Path {
	case "/get":

		service := r.Form.Get("service")
		id := self.NewId(service)

		fmt.Fprint(w, id)
	case "/addservice":
		service := r.Form.Get("service")
		if service == "" {
			fmt.Fprint(w, "invalid service")
			return
		}
		err := self.AddService(service)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		fmt.Fprint(w, "success")

	}
}
