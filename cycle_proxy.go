package cycle_proxy

import (
	"fmt"
	"math/rand"
	"net/url"
	"sync"
)

type CycleProxy struct {
	Mutex   *sync.Mutex
	Proxies []string
	Index   int
}

func NewCycleProxy(proxies []string) *CycleProxy {
	return &CycleProxy{
		Mutex:   &sync.Mutex{},
		Proxies: proxies,
		Index:   0,
	}
}

func (cycle *CycleProxy) Shuffle() {
	//cycle.Index = rand.Intn(cycle.Length() - 1)
	rand.Shuffle(cycle.Length(), func(i, j int) {
		cycle.Proxies[i], cycle.Proxies[j] = cycle.Proxies[j], cycle.Proxies[i]
	})
}

func (cycle *CycleProxy) Next() (*url.URL, error) {
	cycle.Mutex.Lock()
	defer cycle.Mutex.Unlock()

	if cycle.Length() == 0 {
		return nil, fmt.Errorf("proxies is empty")
	}

	cycle.Index = cycle.Index + 1
	if cycle.Index >= cycle.Length() {
		cycle.Index = 0
	}

	return url.Parse(cycle.Proxies[cycle.Index])
}

func (cycle *CycleProxy) Now() (*url.URL, error) {
	cycle.Mutex.Lock()
	defer cycle.Mutex.Unlock()
	return url.Parse(cycle.Proxies[cycle.Index])
}

func (cycle *CycleProxy) Length() int {
	return len(cycle.Proxies)
}
