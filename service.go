package main

import (
	"log"
)

const (
	STEP     = 1000
	CAPACITY = 500
)

type IDGenerator interface {
	Gen()
}

type IDServiceGenerator struct {
	Service string
	Start   uint64
	End     uint64
	Step    uint64
	Hits    uint64

	position uint64
	ch       chan uint64
	bk       *Backend
}

func NewIDServiceGenerator(service string, bk *Backend) *IDServiceGenerator {
	idGen := &IDServiceGenerator{
		Service:  service,
		Step:     STEP,
		Start:    0,
		End:      0,
		bk:       bk,
		position: 0,
		ch:       make(chan uint64, CAPACITY),
	}
	go idGen.generator()
	log.Printf("Get a new service, service=%s", service)
	return idGen
}

func (self *IDServiceGenerator) generator() {
	for {
		err := self.incr()
		if err != nil {
			continue
		}
		self.ch <- self.position
	}
}

func (self *IDServiceGenerator) Gen() uint64 {
	id := <-self.ch
	self.Hits++
	return id
}

func (self *IDServiceGenerator) incr() (err error) {
	self.position++
	if self.position < self.End {
		return nil
	}

	lastId, err := self.bk.GetRange(self.Service, self.Step)
	if err != nil {
		return err
	}
	self.Start = lastId - self.Step
	self.End = lastId
	self.position = self.Start

	log.Printf("service=%s, get id range (%d, %d) ", self.Service, self.Start, self.End-1)
	return nil
}
