package httpmock

import (
	"fmt"
	"math"
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

type uniformDelay struct {
	delayBase
	min time.Duration
	max time.Duration
}

type normalDelay struct {
	delayBase
	mean, stdDev, max float64
}

func FixedDelay(d string) {
	dd, err := time.ParseDuration(d)
	if err != nil {
		panic(fmt.Sprintf("Parsing time for FixedDelay. error = %s", err.Error()))
	}
	fd := fixedDelay{delayBase: delayBase{}, delay: dd}
	fd.waiter = &fd
	DecorateHandler(&fd, NoopHandler)
}

func UniformDelay(min, max string) {
	var ud uniformDelay
	var err error
	ud.min, err = time.ParseDuration(min)
	if err != nil {
		panic(fmt.Sprintf("Parsing min for UniformDelay in error = %s", err.Error()))
	}
	ud.max, err = time.ParseDuration(max)
	if err != nil {
		panic(fmt.Sprintf("Parsing max for UniformDelay in error = %s", err.Error()))
	}
	ud.waiter = &ud
	DecorateHandler(&ud, NoopHandler)
}

func NormalDelay(mean, stdDev, max string) {
	nd := normalDelay{delayBase: delayBase{}, mean: 0, stdDev: 0, max: 0}
	nd.waiter = &nd
	var err error
	meanDuration, err := time.ParseDuration(mean)
	if err != nil {
		panic(fmt.Sprintf("Parsing mean for NormalDelay in error = %s", err.Error()))
	}
	nd.mean = float64(meanDuration) / float64(time.Second)
	stdDevDuration, err := time.ParseDuration(stdDev)
	if err != nil {
		panic(fmt.Sprintf("Parsing stdDev for NormalDelay in error = %s", err.Error()))
	}
	nd.stdDev = float64(stdDevDuration) / float64(time.Second)
	maxDuration, err := time.ParseDuration(max)
	if err != nil {
		panic(fmt.Sprintf("Parsing max for NormalDelay in error = %s", err.Error()))
	}
	nd.max = float64(maxDuration) / float64(time.Second)
	DecorateHandler(&nd, NoopHandler)
}

func (fd *fixedDelay) Wait() error {
	time.Sleep(fd.delay)
	return nil
}

func (ud *uniformDelay) Wait() error {
	time.Sleep(ud.min + time.Duration(rand.Intn(int(ud.max-ud.min))))
	return nil
}

func (nd *normalDelay) NextWaitTime() time.Duration {
	a := math.Log(1 + math.Pow(nd.stdDev/nd.mean, 2))
	u := math.Log(nd.mean) - a/2
	s := math.Sqrt(a)

	seed := math.Exp(float64(u)+float64(s)*rand.NormFloat64())
	return time.Duration(seed * float64(time.Second))
	//if seed < 0 {
	//	if nd.stdDev*5 > nd.mean {
	//		seed = seed * float64(nd.mean) / 5.0
	//	} else {
	//		seed = seed * float64(nd.stdDev)
	//	}
	//} else {
	//	if nd.mean+nd.stdDev*5 > nd.max {
	//		seed = seed * float64(nd.max-nd.mean) / 5.0
	//	} else {
	//		seed = seed * float64(nd.stdDev)
	//	}
	//}
	//
	//return time.Duration(float64(nd.mean) + seed)
}

func (nd *normalDelay) Wait() error {

	time.Sleep(nd.NextWaitTime())
	return nil
}

func (d *delayBase) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	d.waiter.Wait()
}
