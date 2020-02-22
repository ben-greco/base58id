package base58id

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const maxRetries = 3
const expireDuration = 1050 * time.Millisecond
const maxNumChar = 0

// TODO: Find some way of making these non-global. It will 'contaminate' tests. (Potentially make the fields on the broker)
var (
	idLength  = 1
	uniqueMap = make(map[string]*candidate)
	mapLock   sync.Mutex
)

type candidate struct {
	ID         string
	expiration time.Time
}

func addUnique(s string) {
	mapLock.Lock()
	defer mapLock.Unlock()
	uniqueMap[s] = &candidate{ID: s, expiration: time.Now().Add(expireDuration)}
}

func purgeFromUnique() {
	mapLock.Lock()
	defer mapLock.Unlock()
	if len(uniqueMap) == 0 {
		return
	}
	for k, v := range uniqueMap {
		if !v.expiration.After(time.Now()) {
			delete(uniqueMap, k)
		}
	}
}

func isUnique(s string) bool {
	mapLock.Lock()
	defer mapLock.Unlock()
	return uniqueMap[s] == nil
}

func getNumericChars(num int) string {
	var s = ""
	for i := 0; i < num; i++ {
		s += fmt.Sprint(rand.Intn(maxNumChar + 1))
	}
	return s
}
