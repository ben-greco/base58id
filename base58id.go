package base58id

import (
	"bytes"
	"errors"
	"log"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/akamensky/base58"
)

var (
	numericChars   = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	uniqueMap      = make(map[string]*candidate)
	idLength       = 1
	maxRetries     = 3
	expireDuration = 1050 * time.Millisecond
	mapLock        sync.Mutex
)

type candidate struct {
	ID         string
	expiration time.Time
}

type ShortIDServer struct {
	shortest    bool
	capacity    int
	instanceID  string
	createChan  chan string
	requestChan chan bool
	receiveChan chan string
	ids         []string
}

// Get retrieves a single ID from a short ID server
func (s *ShortIDServer) Get() string {
	s.requestChan <- true
	t := <-s.receiveChan

	return t
}

// Get many retrieves n IDs from a short ID server
func (s *ShortIDServer) GetMany(n int) []string {
	var send []string
	for i := 0; i < n; i++ {
		send = append(send, s.Get())
	}

	return send
}

func (s *ShortIDServer) initialize() {
	s.createChan = make(chan string)
	s.receiveChan = make(chan string)
	s.requestChan = make(chan bool)

	go s.generator()

	go s.distributor()
}

func (s *ShortIDServer) generator() {
	for {
		s.createChan <- s.getUniqueID()
	}
}

func (s *ShortIDServer) distributor() {
	ticker := time.NewTicker(expireDuration)

	for {
		switch len(s.ids) {
		case 0:
			t := <-s.createChan
			s.ids = append(s.ids, t)
		case s.capacity:
			<-s.requestChan

			var x string

			x, s.ids = s.ids[0], s.ids[1:]

			s.receiveChan <- x
		default:
			select {
			case <-s.requestChan:
				var x string
				x, s.ids = s.ids[0], s.ids[1:]
				s.receiveChan <- x
			case t := <-s.createChan:
				s.ids = append(s.ids, t)
			case <-ticker.C:
				purgeFromUnique()

				idLength = 1
			}
		}
	}
}

// New creates a short ID server. Any two short ID servers created without an instance ID or with the same instance ID
// *WILL* generate duplicate IDs in short order. Any two short ID servers with different instance IDs *WILL NOT* create
// duplicates IDs. The capacity parameter determines how many IDs are kept in memory at any time. Capacity is always at
// least 1 no matter what input is given. See parameters section of the README to learn more about tradeoffs related to
// speed versus in memory storage in your application.
func New(capacity int, instanceID ...int) (*ShortIDServer, error) {
	if capacity < 1 {
		capacity = 1
	}

	s := ShortIDServer{
		capacity: capacity,
	}

	if len(instanceID) > 1 {
		return nil, errors.New("two or more machine IDs were given, only zero or one ID is allowed")
	}

	if len(instanceID) > 0 {
		// validate their id does not have a zero init
		id := strconv.Itoa(instanceID[0])
		if strings.Contains(id, "0") {
			return nil, errors.New("your instance ID contained a zero which is not allowed")
		}

		s.instanceID = id
		s.shortest = false
	} else {
		s.shortest = true
	}

	s.initialize()

	return &s, nil
}

func (s *ShortIDServer) getUniqueID() string {
	t := s.newSingleID(idLength)

	retries := 0

	for !isUnique(t) {
		retries++
		if retries > maxRetries {
			idLength++

			retries = 0
		}

		t = s.newSingleID(idLength)
	}

	addUnique(t)

	return t
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

func (s *ShortIDServer) newSingleID(length int) string {
	num := getNumericChars(length)

	secondsString := strconv.FormatInt(time.Now().Unix(), 10)

	var uniqString string
	if s.shortest {
		uniqString = num + "0" + secondsString
	} else {
		uniqString = num + "0" + s.instanceID + "0" + secondsString
	}

	my128 := big.NewInt(0)

	myInt, ok := my128.SetString(uniqString, 10)
	if !ok {
		log.Println("base58id: problem in set string!")
	}

	return base58.Encode(myInt.Bytes())
}

func getNumericChars(num int) string {
	var b bytes.Buffer

	for i := 0; i < num; i++ {
		b.WriteString(numericChars[rand.Intn(len(numericChars))])
	}

	return b.String()
}
