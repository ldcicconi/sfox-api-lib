package sfoxapi

import (
	"fmt"
	"sync"
	"time"

	"github.com/Xuyuanp/trie"
)

// Something I wrote to keep track of api errors in a production environment
type ErrorSourceKey string

const (
	CreateOrderKey   ErrorSourceKey = "CreateOrder"
	OrderStatusKey                  = "OrderStatus"
	GetOpenOrdersKey                = "GetOpenOrders"
	CancelOrderKey                  = "CancelOrder"
	BalanceKey                      = "Balance"
)

type Monitor struct {
	lock   sync.Mutex
	errors map[ErrorSourceKey]*trie.Trie
}

func NewMonitor() *Monitor {
	return &Monitor{
		errors: make(map[ErrorSourceKey]*trie.Trie),
	}
}

func (monitor *Monitor) RecordError(source ErrorSourceKey, err error) {
	// default async lol
	go func() {
		monitor.lock.Lock()
		defer monitor.lock.Unlock()
		if monitor.errors[source] == nil {
			monitor.errors[source] = trie.NewTrie()
		}
		monitor.errors[source].Insert(err.Error())
		return
	}()
}

func (monitor *Monitor) Start() {
	go func() {
		lastTime := time.Now()
		nextTime := lastTime.Truncate(time.Hour)
		for {
			nextTime = nextTime.Add(time.Hour)
			waitingTime := time.Until(nextTime)
			time.Sleep(waitingTime)
			// log stuff
			monitor.lock.Lock()
			fmt.Printf("SFOX-API-LIB errors since %s (%f min) :\n", lastTime.Format(time.RFC822Z), nextTime.Sub(lastTime).Minutes())
			printErrorMap(monitor.errors)
			lastTime = nextTime
			// flush the tries for garbage collection
			monitor.errors = make(map[ErrorSourceKey]*trie.Trie)
			monitor.lock.Unlock()
		}
	}()
}

func printErrorMap(m map[ErrorSourceKey]*trie.Trie) {
	for s, trie := range m {
		errors := errorsAndFrequenciesToString(trie)
		if errors == "" {
			continue
		}
		fmt.Printf("Errors from %s:\n", s)
		fmt.Print(errors)
	}
}

func errorsAndFrequenciesToString(t *trie.Trie) string {
	var str string
	m := t.Travel()
	for e, freq := range m {
		str += fmt.Sprintf("freq: %v error: %s\n", freq, e)
	}
	return str
}
