package main

type IdSequenceInterface interface {
	NextId() int64
}

func newIdSequence() *IdSequence {
	return &IdSequence{}
}

type IdSequence struct {
	nextId int64
}

func (r *IdSequence) NextId() int64 {
	r.nextId++
	return r.nextId
}
