package base58id

import (
	"fmt"
	"testing"
	"time"
)

func TestInvalidInstanceID(t *testing.T) {
	_, err := New(3, 101)
	if err == nil {
		t.Error("this should have been an error")
		t.Fail()
		return
	}
}

func TestMachineIDUniqueness(t *testing.T) {
	var b, c, d, e *ShortIDServer
	a, err := New(3, 1)
	if err != nil {
		t.Errorf("error creating new server: %v", err)
		t.Fail()
		return
	}
	b, err = New(3, 2)
	if err != nil {
		t.Errorf("error creating new server: %v", err)
		t.Fail()
		return
	}
	c, err = New(3, 3)
	if err != nil {
		t.Errorf("error creating new server: %v", err)
		t.Fail()
		return
	}
	d, err = New(3, 3943)
	if err != nil {
		t.Errorf("error creating new server: %v", err)
		t.Fail()
		return
	}
	e, err = New(3, 3944)
	if err != nil {
		t.Errorf("error creating new server: %v", err)
		t.Fail()
		return
	}
	m := make(map[string]bool)
	for i := 0; i < 2000000; i++ {
		id := a.Get()
		if m[id] {
			t.Errorf("this id from a was already found: %v", id)
			t.Fail()
			return
		} else {
			m[id] = true
		}
		id = b.Get()
		if m[id] {
			t.Errorf("this id from b was already found: %v", id)
			t.Fail()
			return
		} else {
			m[id] = true
		}
		id = c.Get()
		if m[id] {
			t.Errorf("this id from c was already found: %v", id)
			t.Fail()
			return
		} else {
			m[id] = true
		}
		id = d.Get()
		if m[id] {
			t.Errorf("this id from d was already found: %v", id)
			t.Fail()
			return
		} else {
			m[id] = true
		}
		id = e.Get()
		if m[id] {
			t.Errorf("this id from e was already found: %v", id)
			t.Fail()
			return
		} else {
			m[id] = true
		}
	}
}

func TestMaxLength(t *testing.T) {
	s, err := New(1, 9999)
	if err != nil {
		t.Errorf("error creating new server: %v", err)
		t.Fail()
		return
	}
	var id string
	maxLength := 0
	longestID := ""
	for i := 0; i < 5000000; i++ {
		id = s.Get()
		if len(id) > maxLength {
			longestID = id
			maxLength = len(id)
		}
	}
	if maxLength > 15 {
		fmt.Println("This was the longest id: ", longestID, " it has length: ", maxLength)
		t.Fail()
	}
}

func TestSpeedAndUniquenessSingle(t *testing.T) {
	s, err := New(100)
	if err != nil {
		t.Errorf("error creating new server: %v", err)
		t.Fail()
		return
	}
	m := make(map[string]bool)
	start := time.Now()
	for i := 0; i < 1800000; i++ {
		id := s.Get()
		if m[id] {
			t.Fail()
		} else {
			m[id] = true
		}
	}
	if time.Since(start) > 10*time.Second {
		t.Log("Did not average 180k per second or more!")
		t.Fail()
		return
	}
}
