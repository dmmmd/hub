package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdSequenceImplementsIdSequenceInterface(t *testing.T) {
	var _ IdSequenceInterface = (*IdSequence)(nil)
}

func TestNextIdAlwaysIncrements(t *testing.T) {
	s := newIdSequence()

	id1 := s.NextId()
	id2 := s.NextId()
	id3 := s.NextId()
	id4 := s.NextId()

	assert.True(t, id1 < id2)
	assert.True(t, id2 < id3)
	assert.True(t, id3 < id4)
}
