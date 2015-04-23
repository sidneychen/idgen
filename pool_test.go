package main

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestNewIDPool(t *testing.T) {
	dbDsn := "xxj:123@tcp(192.168.94.26:3316)/service_api?charset=utf8"
	dbType := "mysql"
	db, err := sql.Open(dbType, dbDsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	pool := NewIDPool(db, "id_gen", "goods_id")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 300; i++ {
		id := pool.GetID()
		t.Log(id)
	}
}