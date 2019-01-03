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

// Defines a fixed delay for the response.  The duration string should be formatted as expected by time.ParseDuration
func FixedDelay(d string) {
	dd, err := time.ParseDuration(d)
	if err != nil {
		panic(fmt.Sprintf("Parsing time for FixedDelay. error = %s", err.Error()))
	}
	fd := fixedDelay{delayBase: delayBase{}, delay: dd}
	fd.waiter = &fd
	DecorateHandler(&fd, NoopHandler)
}

// Defines a delay that is uniformly distributed between a minimum and a maximum.  The min and max parameters are
// expected to conform the the format expected by time.ParseDuration
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

func nextWaitTimeNormal(m, u, s float64) func() time.Duration {
	return func() time.Duration {
		return time.Duration(math.Min(m, math.Exp(float64(u)+float64(s)*rand.NormFloat64())) * float64(time.Second))
	}
}

// Defines a delay whose distribution conforms to a Normal Distribution with the given mean, standard deviation, and
// maximum.  All the durations are expressed in a format compatible with time.ParseDuration
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
	a := 1 + (stdDevF*stdDevF)/math.Pow(meanF, 2)
	// u := math.Log(meanF) - a/2
	u := math.Log(meanF / math.Sqrt(a))
	s := math.Sqrt(math.Log(a))
	DecorateHandler(Waiter(nextWaitTimeNormal(maxF, u, s)), NoopHandler)
}

func nextWaitTimeSmoothedNormal(max, u, s float64) func() time.Duration {
	next := make(chan time.Duration)

	getNext := func() float64 {
		return math.Exp(float64(u) + float64(s)*rand.NormFloat64())
	}
	a1 := 0.3
	//a2 := 0.5
	go func() {
		var last1 float64 = math.Exp(u)
		//var last2 float64 = math.Exp(u)
		for {
			x1 := getNext()
			n := a1*x1 + (1.0-a1)*last1
			last1 = n
			//n := a2*x2 + (1.0-a2)*last2
			//last2 = n
			next <- time.Duration(math.Min(max, n) * float64(time.Second))
		}
	}()
	return func() time.Duration {
		return <-next
	}
}

// WIP, This is an incomplete API and should not be used.
func SmoothedNormalDelay(mean, stdDev, max string) {
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
	a := 1 + (stdDevF*stdDevF)/math.Pow(meanF, 2)
	// u := math.Log(meanF) - a/2
	u := math.Log(meanF / math.Sqrt(a))
	s := math.Sqrt(math.Log(a))
	DecorateHandler(Waiter(nextWaitTimeSmoothedNormal(maxF, u, s)), NoopHandler)
}

// Defines a generic waiter that will use the provided waitTime function to acquire the duration to wait.
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
