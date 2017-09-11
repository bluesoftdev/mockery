package httpMock

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type waiter interface {
	Wait() error
}

type delayBase struct {
	waiter waiter
}

type fixedDelay struct {
	delayBase
	delay time.Duration
}

type normalDelay struct {
	delayBase
	mean   time.Duration
	stdDev time.Duration
	max    time.Duration
}

func FixedDelay(d string) {
	dd, err := time.ParseDuration(d)
	if err != nil {
		panic(fmt.Sprintf("Parsing time for FixedDelay in %s:%s error = %s", currentMock.url, currentMockMethod.method, err.Error()))
	}
	fd := fixedDelay{delayBase: delayBase{}, delay: dd}
	fd.waiter = &fd
	DecorateHandler(&fd, NoopHandler)
}

func NormalDelay(mean, stdDev, max string) {
	nd := normalDelay{delayBase: delayBase{}, mean: 0, stdDev: 0, max: 0}
	nd.waiter = &nd
	var err error
	nd.mean, err = time.ParseDuration(mean)
	if err != nil {
		panic(fmt.Sprintf("Parsing mean for NormalDelay in %s:%s error = %s", currentMock.url, currentMockMethod.method, err.Error()))
	}
	nd.stdDev, err = time.ParseDuration(stdDev)
	if err != nil {
		panic(fmt.Sprintf("Parsing stdDev for NormalDelay in %s:%s error = %s", currentMock.url, currentMockMethod.method, err.Error()))
	}
	nd.max, err = time.ParseDuration(max)
	if err != nil {
		panic(fmt.Sprintf("Parsing max for NormalDelay in %s:%s error = %s", currentMock.url, currentMockMethod.method, err.Error()))
	}
	DecorateHandler(&nd, NoopHandler)
}

func (fd *fixedDelay) Wait() error {
	time.Sleep(fd.delay)
	return nil
}

func (nd *normalDelay) NextWaitTime() time.Duration {
	seed := rand.NormFloat64()

	if seed < 0 {
		if nd.stdDev*5 > nd.mean {
			seed = seed * float64(nd.mean) / 5.0
		} else {
			seed = seed * float64(nd.stdDev)
		}
	} else {
		if nd.mean+nd.stdDev*5 > nd.max {
			seed = seed * float64(nd.max-nd.mean) / 5.0
		} else {
			seed = seed * float64(nd.stdDev)
		}
	}

	return time.Duration(float64(nd.mean) + seed)
}

func (nd *normalDelay) Wait() error {

	time.Sleep(nd.NextWaitTime())
	return nil
}

func (d *delayBase) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	d.waiter.Wait()
}
