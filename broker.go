package base58id

import (
	"github.com/akamensky/base58"
	"log"
	"math/big"
	"strconv"
	"time"
)

// TODO: comment broker
// TODO: comment fields in broker
// Broker is a ...
type Broker struct {
	shortest    bool
	capacity    int
	instanceID  string
	createChan  chan string
	requestChan chan bool
	receiveChan chan string
	ids         []string
}

// TODO: use capacity to make a channel with capacity rather than tracking it manually
// TODO: implement some kind of destructor that stops id generation and frees the extra ids
// TODO: allow a capacity of 0, that only generates an id on command
/*
New creates a short ID broker. Any two ID brokers created with without an instance ID WILL generate duplicate IDs in
short order. Any two short ID brokers with different instance IDs will not create duplicates for hundreds of years.
The capacity parameter determines how many IDs are kept in memory at any time. Capacity is always at least 1 no
matter what input is given. For applications that could generate
*/
func New(options ...BrokerOption) (*Broker, error) {
	// Initialize with defaults
	b := &Broker{
		capacity:    1,
		shortest:    true,
		createChan:  make(chan string),
		receiveChan: make(chan string),
		requestChan: make(chan bool),
	}

	// Alter struct based on given options
	for _, option := range options {
		err := option(b)
		if err != nil {
			return nil, err
		}
	}

	go b.generator()
	go b.distributor()

	return b, nil
}

// Next retrieves a single ID from a short ID broker
func (b *Broker) Next() string {
	b.requestChan <- true
	t := <-b.receiveChan
	return t
}

// Many retrieves n IDs from a short ID broker
func (b *Broker) Many(n int) []string {
	var send []string
	for i := 0; i < n; i++ {
		send = append(send, b.Next())
	}
	return send
}

func (b *Broker) generator() {
	for {
		b.createChan <- b.getUniqueID()
	}
}

// TODO: comment cases (what each one does and why)
// TODO: reflect changes from making createChan a channel with capacity (select statement with <-requestChan and <-ticker.C)
func (b *Broker) distributor() {
	ticker := time.NewTicker(expireDuration)
	for {
		switch len(b.ids) {
		case 0:
			t := <-b.createChan
			b.ids = append(b.ids, t)
		case b.capacity:
			<-b.requestChan
			var x string
			x, b.ids = b.ids[0], b.ids[1:]
			b.receiveChan <- x
		default:
			select {
			case <-b.requestChan:
				var x string
				x, b.ids = b.ids[0], b.ids[1:]
				b.receiveChan <- x
			case t := <-b.createChan:
				b.ids = append(b.ids, t)
			case <-ticker.C:
				purgeFromUnique()
				idLength = 1
			}
		}
	}
}

func (b *Broker) getUniqueID() string {
	t := b.newSingleID(idLength)
	retries := 0
	for !isUnique(t) {
		retries++
		if retries > maxRetries {
			idLength++
			retries = 0
		}
		t = b.newSingleID(idLength)
	}
	addUnique(t)
	return t
}

func (b *Broker) newSingleID(length int) string {

	num := getNumericChars(length)
	secondsString := strconv.FormatInt(time.Now().Unix(), 10)
	var uniqString string
	if b.shortest {
		uniqString = num + "0" + secondsString
	} else {
		uniqString = num + "0" + b.instanceID + "0" + secondsString
	}
	my128 := big.NewInt(0)
	myInt, ok := my128.SetString(uniqString, 10)
	if !ok {
		log.Println("base58id: problem in set string!")
	}
	return base58.Encode(myInt.Bytes())
}
