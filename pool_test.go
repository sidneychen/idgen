package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	//	_ "github.com/go-sql-driver/mysql"
)

func TestConfig(t *testing.T) {
	//	cfg := Config{}

	data, _ := ioutil.ReadFile("./config.json")

	var cfg Config
	json.Unmarshal(data, &cfg)

	t.Logf("step is %v", cfg.DB)

}

func TestBackend(t *testing.T) {
	cfg := NewConfigFromFile("./config.json")
	bk, err := NewBackendFromConfig(cfg.DB)
	if err != nil {
		t.Logf("create backend, err=%v", err)
	}
	t.Logf("bk=%v", bk)

}

func TestService(t *testing.T) {
	cfg := NewConfigFromFile("./config.json")
	bk, err := NewBackendFromConfig(cfg.DB)
	if err != nil {
		t.Logf("create backend, err=%v", err)
	}
	s := NewIDServiceGenerator("goods_id", bk)
	for i := 0; i < 100; i++ {
		id := s.Gen()
		t.Logf("get an id, id=%v", id)

	}

}

//func TestNew(t *testing.T) {
//	dbDsn := "xxj:123@tcp(192.168.94.26:3316)/service_api?charset=utf8"
//	dbType := "mysql"
//	db, err := sql.Open(dbType, dbDsn)
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer db.Close()

//	pool := NewIDPool(db, "id_gen", "goods_id")
//	if err != nil {
//		t.Fatal(err)
//	}

//	for i := 0; i < 300; i++ {
//		id := pool.GetID()
//		t.Log(id)
//	}
//}
