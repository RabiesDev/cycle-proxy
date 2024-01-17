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

func (cycle *CycleProxy) Randomize() *CycleProxy {
	cycle.Index = rand.Intn(len(cycle.Proxies) - 1)
	return cycle
}

func (cycle *CycleProxy) NextProxy() (*url.URL, error) {
	nextProxy, err := cycle.Next()
	if err != nil {
		return nil, err
	}
	return url.Parse(nextProxy)
}

func (cycle *CycleProxy) Next() (string, error) {
	cycle.Mutex.Lock()
	defer cycle.Mutex.Unlock()

	if len(cycle.Proxies) == 0 {
		return "", fmt.Errorf("proxies is empty")
	}

	cycle.Index = cycle.Index + 1
	if cycle.Index >= len(cycle.Proxies) {
		cycle.Index = 0
	}

	return cycle.Proxies[cycle.Index], nil
}

func (cycle *CycleProxy) CurrentProxy() (*url.URL, error) {
	return url.Parse(cycle.Current())
}

func (cycle *CycleProxy) Current() string {
	cycle.Mutex.Lock()
	defer cycle.Mutex.Unlock()
	return cycle.Proxies[cycle.Index]
}

func (cycle *CycleProxy) Length() int {
	return len(cycle.Proxies)
}
