package httpMock

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

type waiter interface {
	Wait() error
}

type delayBase struct {
	handler http.Handler
	waiter  waiter
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
	fd := fixedDelay{delayBase: delayBase{handler: currentMockMethod}, delay: dd}
	fd.waiter = &fd
	currentMockMethodHandler = &fd
}

func NormalDelay(mean, stdDev, max string) {
	nd := normalDelay{delayBase: delayBase{handler: currentMockMethod}, mean: 0, stdDev: 0, max: 0}
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
	currentMockMethodHandler = &nd
}

func (fd *fixedDelay) Wait() error {
	time.Sleep(fd.delay)
	return nil
}

func (nd *normalDelay) Wait() error {
	seed := rand.NormFloat64() * float64(nd.stdDev)
	var scaled float64
	if seed < 0 {
		scaled = (float64(nd.mean) / math.MaxFloat64) * seed
	} else {
		scaled = (float64(nd.max-nd.mean) / math.MaxFloat64) * seed
	}
	waitTime := time.Duration(float64(nd.mean) + scaled)
	time.Sleep(waitTime)
	return nil
}

func (d *delayBase) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := d.waiter.Wait()
	if err != nil {
		log.Printf("ERROR: waiting for waiter to complete: %s", err.Error())
		w.WriteHeader(500)
	} else {
		d.handler.ServeHTTP(w, req)
	}
}
