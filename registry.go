package main

type Registry struct {
	nextId int64
}

func (r *Registry) NextId() int64 {
	r.nextId++
	return r.nextId
}