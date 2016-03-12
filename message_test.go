package main

import "testing"

func TestRelayCommandReturnsRelay(t *testing.T) {
	m := NewRelayMessage([]int64{100500, 42}, "irrelevant")
	if MessageTypeRelay != m.Command() {
		t.Fail()
	}
}

func TestRelayBodyReturnsBody(t *testing.T) {
	body := "testBody\nline2 foobar"
	m := NewRelayMessage([]int64{100500, 42}, body)
	if body != m.Body() {
		t.Fail()
	}
}

func TestRelayReceiversReturnsReceivers(t *testing.T) {
	receivers := []int64{100500, 42}
	m := NewRelayMessage(receivers, "irrelevant")
	r := m.Receivers()

	switch {
	case len(r) != 2:
		t.Fail()
	case r[0] != 100500:
		t.Fail()
	case r[1] != 42:
		t.Fail()
	}
}

func TestIdentityCommandReturnsIdentity(t *testing.T) {
	m := NewIdentityMessage()
	if MessageTypeIdentity != m.Command() {
		t.Fail()
	}
}

func TestIdentityBodyReturnsIsEmpty(t *testing.T) {
	m := NewIdentityMessage()
	if "" != m.Body() {
		t.Fail()
	}
}

func TestIdentityReceiversAreEmpty(t *testing.T) {
	m := NewIdentityMessage()
	r := m.Receivers()

	if 0 != len(r) {
		t.Fail()
	}
}

func TestListCommandReturnsList(t *testing.T) {
	m := NewListMessage()
	if MessageTypeList != m.Command() {
		t.Fail()
	}
}

func TestListBodyReturnsIsEmpty(t *testing.T) {
	m := NewListMessage()
	if "" != m.Body() {
		t.Fail()
	}
}

func TestListReceiversAreEmpty(t *testing.T) {
	m := NewListMessage()
	r := m.Receivers()

	if 0 != len(r) {
		t.Fail()
	}
}
