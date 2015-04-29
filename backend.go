package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	ErrDBConn = errors.New("get sql db error")
)

const (
	STATUS_ON  = 1
	STATUS_OFF = 0
)

const (
	UPDATE_SQL    = "UPDATE `id_gen` SET `position`=`position`+?, `update_time`=? WHERE `service`=?"
	QUERYROW_SQL  = "SELECT `position` FROM `id_gen` WHERE service=? FOR UPDATE"
	QUERYLIST_SQL = "SELECT `service`, `position` FROM `id_gen` WHERE `status`=?"
	INSERT_SQL    = "INSERT INTO `id_gen` (`service`, `position`, `update_time`, `status`) VALUES (?, ?, ?, ?)"
)

type Backend struct {
	db    *sql.DB
	table string
}

func getDBDsn(cfg *DBConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
}

func NewBackendFromConfig(cfg *DBConfig) (*Backend, error) {
	dsn := getDBDsn(cfg)
	db, err := sql.Open(cfg.Type, dsn)
	if err != nil {
		return nil, ErrDBConn
	}
	return &Backend{db, cfg.TableName}, nil
}

// get all service
func (self *Backend) GetAllService(srvs *[]*IDServiceGenerator) error {
	rows, err := self.db.Query(QUERYLIST_SQL, STATUS_ON)
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
			log.Printf("Scan a row fail, error=%v", err)
			continue
		}
		srv := NewIDServiceGenerator(service, self)
		*srvs = append(*srvs, srv)
	}
	return nil
}

// create a new service
func (self *Backend) CreateService(service string) (*IDServiceGenerator, error) {

	_, err := self.db.Exec(INSERT_SQL, service, 1, time.Now().Unix(), STATUS_ON)
	if err != nil {
		return nil, err
	}
	return NewIDServiceGenerator(service, self), nil
}

// get a range of id
func (self *Backend) GetRange(service string, step uint64) (lastID uint64, err error) {
	tx, err := self.db.Begin()
	if err != nil {
		return ERR_ID, err
	}

	tx.Exec(UPDATE_SQL, step, time.Now().Unix(), service)

	err = tx.QueryRow(QUERYROW_SQL, service).Scan(&lastID)
	if err != nil {
		return ERR_ID, err
	}

	err = tx.Commit()
	if err != nil {
		return ERR_ID, err
	}

	return lastID, nil
}
