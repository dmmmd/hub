package main

import "testing"

func TestRelayCommandReturnsRelay(t *testing.T) {
	m := NewRelayMessage(9001, []int64{100500, 42}, "irrelevant")
	if MessageTypeRelay != m.Command() {
		t.Fail()
	}
}

func TestRelayBodyReturnsBody(t *testing.T) {
	body := "testBody\nline2 foobar"
	m := NewRelayMessage(9001, []int64{100500, 42}, body)
	if body != m.Body() {
		t.Fail()
	}
}

func TestRelayReceiversReturnsReceivers(t *testing.T) {
	receivers := []int64{100500, 42}
	m := NewRelayMessage(9001, receivers, "irrelevant")
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
	m := NewIdentityMessage(42)
	if MessageTypeIdentity != m.Command() {
		t.Fail()
	}
}

func TestIdentityBodyReturnsIsEmpty(t *testing.T) {
	m := NewIdentityMessage(42)
	if "" != m.Body() {
		t.Fail()
	}
}

func TestIdentityReceiversAreEmpty(t *testing.T) {
	m := NewIdentityMessage(42)
	r := m.Receivers()

	if 0 != len(r) {
		t.Fail()
	}
}

func TestListCommandReturnsList(t *testing.T) {
	m := NewListMessage(42)
	if MessageTypeList != m.Command() {
		t.Fail()
	}
}

func TestListBodyReturnsIsEmpty(t *testing.T) {
	m := NewListMessage(42)
	if "" != m.Body() {
		t.Fail()
	}
}

func TestListReceiversAreEmpty(t *testing.T) {
	m := NewListMessage(42)
	r := m.Receivers()

	if 0 != len(r) {
		t.Fail()
	}
}
