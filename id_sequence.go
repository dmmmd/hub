package main

import "sync"

type IdSequenceInterface interface {
	NextId() uint64
}

func newIdSequence() *IdSequence {
	return &IdSequence{}
}

type IdSequence struct {
	nextId  uint64
	idMutex sync.Mutex
}

func (r *IdSequence) NextId() uint64 {
	r.idMutex.Lock()
	defer r.idMutex.Unlock()

	r.nextId++
	return r.nextId
}
