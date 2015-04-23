package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	ERR_ID     = 0
	STATUS_ON  = 1
	STATUS_OFF = 0
)

type IDPoolMap map[string]*IDPool

type IDPoolMgr struct {
	//	stats       *Stats
	data      IDPoolMap
	db        *sql.DB
	tablename string
	locker    *sync.RWMutex
}

func NewIDPoolMgr() *IDPoolMgr {
	return &IDPoolMgr{
		//		stats:       &Stats{},
		data:   make(IDPoolMap),
		locker: new(sync.RWMutex),
	}
}

func getDBDsn(cfg *DBConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
}

func (self *IDPoolMgr) LoadConfig(cfg *Config) (err error) {
	dbDsn := getDBDsn(cfg.DB)
	db, err := sql.Open(cfg.DB.Type, dbDsn)
	if err != nil {
		panic(err)
	}

	log.Printf("connected to db, %s", dbDsn)
	self.tablename = cfg.TableName
	self.db = db
	self.loadAllIDPoolFromDb()
	return nil
}

func (self *IDPoolMgr) RegisterIDPool(pool *IDPool) {
	self.data[pool.Service] = pool
}

// 获取新的id
func (self *IDPoolMgr) NewId(service string) uint64 {

	pool, ok := self.data[service]
	if !ok {
		return ERR_ID
	}
	return pool.GetID()
}

func (self *IDPoolMgr) loadAllIDPoolFromDb() (err error) {

	sql := fmt.Sprintf("SELECT `service`, `position` FROM `%s` WHERE `status`=?", self.tablename)
	rows, err := self.db.Query(sql, STATUS_ON)
	if err != nil {
		return err
	}
	defer rows.Close()

	var (
		position uint64
		service  string
	)
	for rows.Next() {
		err = rows.Scan(&service, &position)
		if err != nil {
			return err
		}
		pool := NewIDPool(self.db, self.tablename, service)
		self.RegisterIDPool(pool)
	}
	return nil
}

// 获取新的id
func (self *IDPoolMgr) AddService(service string) (err error) {

	if _, ok := self.data[service]; !ok {
		sql := fmt.Sprintf("INSERT INTO `%s` (`service`, `position`, `update_time`, `status`) VALUES (?, ?, ?, ?)", self.tablename)
		_, err := self.db.Exec(sql, service, 1, time.Now().Unix(), STATUS_ON)
		if err != nil {
			return err
		}
		pool := NewIDPool(self.db, self.tablename, service)
		self.RegisterIDPool(pool)
	}
	return nil
}

func (self *IDPoolMgr) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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
