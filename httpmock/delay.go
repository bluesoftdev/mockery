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

func nextWaitTimeNormal(m, u,s float64) func() time.Duration {
	return func() time.Duration {
		return time.Duration(math.Min(m, math.Exp(float64(u)+float64(s)*rand.NormFloat64())) * float64(time.Second))
	}
}

func NormalDelay(mean, stdDev, max string) {
	var err error
	meanDuration, err := time.ParseDuration(mean)
	if err != nil {
		panic(fmt.Sprintf("Parsing mean for NormalDelay in error = %s", err.Error()))
	}
	meanF := float64(meanDuration) / float64(time.Second)
	stdDevDuration, err := time.ParseDuration(stdDev)
	if err != nil {
		panic(fmt.Sprintf("Parsing stdDev for NormalDelay in error = %s", err.Error()))
	}
	stdDevF := float64(stdDevDuration) / float64(time.Second)
	maxDuration, err := time.ParseDuration(max)
	if err != nil {
		panic(fmt.Sprintf("Parsing max for NormalDelay in error = %s", err.Error()))
	}
	maxF := float64(maxDuration) / float64(time.Second)

	// Calculate mu & sigma
	a := math.Log(1 + math.Pow(stdDevF/meanF, 2))
	u := math.Log(meanF) - a/2
	s := math.Sqrt(a)
	DecorateHandler(Waiter(nextWaitTimeNormal(maxF,u,s)), NoopHandler)
}

func Waiter(waitTime func() time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(waitTime())
	})
}

func (fd *fixedDelay) Wait() error {
	time.Sleep(fd.delay)
	return nil
}

func (ud *uniformDelay) Wait() error {
	time.Sleep(ud.min + time.Duration(rand.Intn(int(ud.max-ud.min))))
	return nil
}

func (d *delayBase) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	d.waiter.Wait()
}
