package base58id

import (
	"fmt"
	"testing"
	"time"
)

func TestInvalidInstanceID(t *testing.T) {
	_, err := New(WithCapacity(3), WithInstanceID(101))
	if err == nil {
		t.Error("this should have been an error")
		t.Fail()
		return
	}
}

// TODO: Streamline test if possible. Took over 10 minutes to run on my laptop
func TestMachineIDUniqueness(t *testing.T) {
	var b, c, d, e *Broker
	a, err := New(WithCapacity(3), WithInstanceID(1))
	if err != nil {
		t.Errorf("error creating new server: %v", err)
		t.Fail()
		return
	}
	b, err = New(WithCapacity(3), WithInstanceID(2))
	if err != nil {
		t.Errorf("error creating new server: %v", err)
		t.Fail()
		return
	}
	c, err = New(WithCapacity(3), WithInstanceID(3))
	if err != nil {
		t.Errorf("error creating new server: %v", err)
		t.Fail()
		return
	}
	d, err = New(WithCapacity(3), WithInstanceID(3943))
	if err != nil {
		t.Errorf("error creating new server: %v", err)
		t.Fail()
		return
	}
	e, err = New(WithCapacity(3), WithInstanceID(3944))
	if err != nil {
		t.Errorf("error creating new server: %v", err)
		t.Fail()
		return
	}
	m := make(map[string]bool)
	for i := 0; i < 2000000; i++ {
		id := a.Next()
		if m[id] {
			t.Errorf("this id from a was already found: %v", id)
			t.Fail()
			return
		} else {
			m[id] = true
		}
		id = b.Next()
		if m[id] {
			t.Errorf("this id from b was already found: %v", id)
			t.Fail()
			return
		} else {
			m[id] = true
		}
		id = c.Next()
		if m[id] {
			t.Errorf("this id from c was already found: %v", id)
			t.Fail()
			return
		} else {
			m[id] = true
		}
		id = d.Next()
		if m[id] {
			t.Errorf("this id from d was already found: %v", id)
			t.Fail()
			return
		} else {
			m[id] = true
		}
		id = e.Next()
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
	s, err := New(WithCapacity(3), WithInstanceID(9999))
	if err != nil {
		t.Errorf("error creating new server: %v", err)
		t.Fail()
		return
	}
	var id string
	maxLength := 0
	longestID := ""
	for i := 0; i < 5000000; i++ {
		id = s.Next()
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
	s, err := New(WithCapacity(100))
	if err != nil {
		t.Errorf("error creating new server: %v", err)
		t.Fail()
		return
	}
	m := make(map[string]bool)
	start := time.Now()
	for i := 0; i < 1800000; i++ {
		id := s.Next()
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
