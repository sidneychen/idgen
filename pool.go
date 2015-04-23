package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	STEP     = 1000
	CAPACITY = 100
)

type IDPool struct {
	Service   string
	Start     uint64
	End       uint64
	Step      uint64
	Hits      uint64
	position  uint64
	db        *sql.DB
	tablename string
	ch        chan uint64
}

func NewIDPool(db *sql.DB, tablename string, service string) *IDPool {
	pool := &IDPool{
		Service:   service,
		Step:      STEP,
		db:        db,
		tablename: tablename,
		position:  0,
		ch:        make(chan uint64, CAPACITY),
	}
	go pool.generator()
	return pool
}

func (self *IDPool) generator() {
	for {
		err := self.incr()
		if err != nil {
			continue
		}
		self.ch <- self.position
	}
}

func (self *IDPool) incr() (err error) {
	self.position++
	if self.position < self.End {
		return nil
	}

	err = self.allocIdRange()
	if err != nil {
		return err
	}
	self.position = self.Start
	return nil
}

func (self *IDPool) GetID() uint64 {
	id := <-self.ch
	self.Hits++
	return id
}

func (self *IDPool) allocIdRange() (err error) {
	tx, err := self.db.Begin()
	if err != nil {
		return err
	}

	var sql string
	sql = fmt.Sprintf("UPDATE `%s` SET `position`=`position`+?, `update_time`=? WHERE `service`=?", self.tablename)
	tx.Exec(sql, self.Step, time.Now().Unix(), self.Service)

	var position uint64
	sql = fmt.Sprintf("SELECT `position` FROM `%s` WHERE service=? FOR UPDATE", self.tablename)
	err = tx.QueryRow(sql, self.Service).Scan(&position)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return err
	}
	self.Start = position - self.Step
	self.position = self.Start
	self.End = position
	log.Printf("%s: get id range (%d, %d) ", self.Service, self.Start, self.End-1)

	return err
}
