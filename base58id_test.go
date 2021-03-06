package base58id

import (
	"fmt"
	"testing"
	"time"
)

func TestFastSingle(t *testing.T) {
	s, err := New(1)
	if err != nil {
		t.Fail()
	}
	for i := 0; i < 400; i++ {
		time.Sleep(300 * time.Millisecond)
		fmt.Println(s.Get())
	}
}

func TestInvalidInstanceID(t *testing.T) {
	_, err := New(3, 101)
	if err == nil {
		t.Error("this should have been an error")
		t.Fail()

		return
	}
}

func TestMachineIDUniqueness(t *testing.T) {
	m := make(map[string]bool)

	var b, d, e *ShortIDServer

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

	for i := 0; i < 2000000; i++ {
		id := a.Get()
		if m[id] {
			t.Errorf("this id from a was already found: %v", id)
			t.Fail()

			return
		}

		m[id] = true

		id = b.Get()
		if m[id] {
			t.Errorf("this id from b was already found: %v", id)
			t.Fail()

			return
		}

		m[id] = true

		id = d.Get()
		if m[id] {
			t.Errorf("this id from d was already found: %v", id)
			t.Fail()

			return
		}

		m[id] = true

		id = e.Get()
		if m[id] {
			t.Errorf("this id from e was already found: %v", id)
			t.Fail()

			return
		}

		m[id] = true
	}
}

func TestMaxLength(t *testing.T) {
	var id string

	s, err := New(1, 9999)
	if err != nil {
		t.Errorf("error creating new server: %v", err)
		t.Fail()

		return
	}

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

	for i := 0; i < 2200000; i++ {
		id := s.Get()
		if m[id] {
			t.Fail()
		} else {
			m[id] = true
		}
	}

	if time.Since(start) > 10*time.Second {
		t.Log("Did not average 220k per second or more!")
		t.Fail()

		return
	}
}
