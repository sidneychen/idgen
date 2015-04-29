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
	log.Printf("Finish all service loaded")
	if err != nil {
		panic(err)
	}
	return mgr
}

func (self *IDServiceMgr) RegisterIDService(ider *IDServiceGenerator) {
	self.data[ider.Service] = ider
}

func (self *IDServiceMgr) NewId(service string) uint64 {

	ider, ok := self.data[service]
	if !ok {
		return ERR_ID
	}
	return ider.Gen()
}

func (self *IDServiceMgr) loadAllService() error {
	srvs := []*IDServiceGenerator{}
	err := self.bk.GetAllService(&srvs)
	if err != nil {
		return err
	}
	log.Print(srvs)

	self.locker.Lock()
	defer self.locker.Unlock()

	for _, srv := range srvs {
		self.RegisterIDService(srv)
	}
	return nil
}

func (self *IDServiceMgr) addService(service string) error {
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

func (self *IDServiceMgr) AddService(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	service := r.Form.Get("service")
	if service == "" {
		fmt.Fprint(w, "invalid service")
		return
	}
	err := self.addService(service)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, "success")
}

func (self *IDServiceMgr) Get(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	service := r.Form.Get("service")
	log.Printf("get a request, method=get, service=%v", service)
	id := self.NewId(service)
	fmt.Fprint(w, id)
}
