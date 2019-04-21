package httpmock

import (
	"bytes"
	"fmt"
	"github.com/montanaflynn/stats"
	"github.com/stretchr/testify/assert"
	"github.com/wcharczuk/go-chart"
	"io"
	"log"
	"math"
	"os"
	"testing"
	"time"
)

func TestFixedDelay(t *testing.T) {
	currentMockHandler = NoopHandler
	FixedDelay("100ms")

	_, samples := runSamples()
	population := stats.LoadRawData(samples)

	mean, err := population.Mean()
	assert.NoError(t, err)
	assert.InDelta(t, 100*float64(time.Millisecond), mean, 10*float64(time.Millisecond))

	median, err := population.Median()
	assert.NoError(t, err)
	assert.InDelta(t, 100*float64(time.Millisecond), median, 10*float64(time.Millisecond))

	p95, err := population.Percentile(85.0)
	assert.NoError(t, err)
	assert.InDelta(t, 100*float64(time.Millisecond), p95, 10*float64(time.Millisecond))

	p975, err := population.Percentile(97.5)
	assert.NoError(t, err)
	assert.InDelta(t, 100*float64(time.Millisecond), p975, 10*float64(time.Millisecond))

	p9985, err := population.Percentile(99.85)
	assert.NoError(t, err)
	assert.InDelta(t, 100*float64(time.Millisecond), p9985, 10*float64(time.Millisecond))
}

func TestUniformDelay(t *testing.T) {
	currentMockHandler = NoopHandler
	UniformDelay("100ms", "200ms")

	_, samples := runSamples()
	population := stats.LoadRawData(samples)

	mean, err := population.Mean()
	assert.NoError(t, err)
	assert.InDelta(t, 150*float64(time.Millisecond), mean, 10*float64(time.Millisecond))

	median, err := population.Median()
	assert.NoError(t, err)
	assert.InDelta(t, 150*float64(time.Millisecond), median, 10*float64(time.Millisecond))

	p95, err := population.Percentile(85.0)
	assert.NoError(t, err)
	assert.InDelta(t, 185*float64(time.Millisecond), p95, 10*float64(time.Millisecond))

	p975, err := population.Percentile(97.5)
	assert.NoError(t, err)
	assert.InDelta(t, 197*float64(time.Millisecond), p975, 10*float64(time.Millisecond))

	p9985, err := population.Percentile(99.85)
	assert.NoError(t, err)
	assert.InDelta(t, 199*float64(time.Millisecond), p9985, 10*float64(time.Millisecond))
}

func TestNormalDelay(t *testing.T) {
	currentMockHandler = NoopHandler
	NormalDelay("100ms", "20ms", "200ms")

	timeSamples, samples := runSamples()
	population := stats.LoadRawData(samples)

	mean, err := population.Mean()
	assert.NoError(t, err)
	assert.InDelta(t, 100*float64(time.Millisecond), mean, 10*float64(time.Millisecond))

	median, err := population.Median()
	assert.NoError(t, err)
	assert.InDelta(t, 100*float64(time.Millisecond), median, 10*float64(time.Millisecond))

	p95, err := population.Percentile(85.0)
	assert.NoError(t, err)
	assert.InDelta(t, 120*float64(time.Millisecond), p95, 10*float64(time.Millisecond))

	p975, err := population.Percentile(97.5)
	assert.NoError(t, err)
	assert.InDelta(t, 140*float64(time.Millisecond), p975, 15*float64(time.Millisecond))

	p9985, err := population.Percentile(99.85)
	assert.NoError(t, err)
	assert.InDelta(t, 160*float64(time.Millisecond), p9985, 25*float64(time.Millisecond))

	buckets := makeBuckets(50, 250*time.Millisecond, 50*time.Millisecond)
	histogram := makeHistogram(buckets, samples)
	sw := bytes.NewBuffer(make([]byte, 0, 512))
	printHistogram(sw, buckets, histogram, 40)
	t.Logf("\n%s", sw.String())
	renderTimeSeries(timeSamples, samples, "./testdata/NormalDelayTest1.png")
}

func TestNormalSmoothedDelay(t *testing.T) {
	currentMockHandler = NoopHandler
	SmoothedNormalDelay("100ms", "20ms", "200ms")

	timeSamples, durationSamples := runSamples()

	population := stats.LoadRawData(durationSamples)

	mean, err := population.Mean()
	assert.NoError(t, err)
	assert.InDelta(t, float64(100*time.Millisecond), mean, 1*float64(time.Millisecond))

	median, err := population.Median()
	assert.NoError(t, err)
	assert.InDelta(t, float64(100*time.Millisecond), median, 2*float64(time.Millisecond))

	p95, err := population.Percentile(85.0)
	assert.NoError(t, err)
	assert.InDelta(t, float64(120*time.Millisecond), p95, 15*float64(time.Millisecond))

	p975, err := population.Percentile(97.5)
	assert.NoError(t, err)
	assert.InDelta(t, float64(140*time.Millisecond), p975, 25*float64(time.Millisecond))

	p9985, err := population.Percentile(99.85)
	assert.NoError(t, err)
	assert.InDelta(t, float64(160*time.Millisecond), p9985, 35*float64(time.Millisecond))

	buckets := makeBuckets(50, 150*time.Millisecond, 50*time.Millisecond)
	histogram := makeHistogram(buckets, durationSamples)
	sw := bytes.NewBuffer(make([]byte, 0, 512))
	printHistogram(sw, buckets, histogram, 40)
	t.Logf("\n%s", sw.String())
	renderTimeSeries(timeSamples, durationSamples, "./testdata/SmooothedNormalDelayTest1.png")
}

func renderTimeSeries(times []time.Time, durations []time.Duration, fileName string) {
	durationFloats := make([]float64, len(durations))
	for i, d := range durations {
		durationFloats[i] = float64(d) / float64(time.Second)
	}
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: times,
				YValues: durationFloats,
			},
		},
	}

	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	graph.Render(chart.PNG, f)
}

const parallel = 1000
const samples = 10000

type timedSample struct {
	timestamp time.Time
	duration  time.Duration
}

func runSamples() ([]time.Time, []time.Duration) {
	numSamples := samples
	timeSamples := make([]time.Time, 0, numSamples)
	durationSamples := make([]time.Duration, 0, numSamples)
	sampleChan := make(chan timedSample, parallel)
	doneChan := make(chan int, parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			for s := 0; s < numSamples/parallel; s++ {
				duration := timeAction(func() {
					currentMockHandler.ServeHTTP(nil, nil)
				})
				sampleChan <- timedSample{time.Now(), duration}
			}
			doneChan <- i
		}()
	}
	done := 0
	for done < parallel {
		select {
		case sample := <-sampleChan:
			timeSamples = append(timeSamples, sample.timestamp)
			durationSamples = append(durationSamples, sample.duration)
		case <-doneChan:
			done++
		}
	}
	close(sampleChan)
	for sample := range sampleChan {
		timeSamples = append(timeSamples, sample.timestamp)
		durationSamples = append(durationSamples, sample.duration)
	}
	return timeSamples, durationSamples
}

func makeHistogram(buckets, samples []time.Duration) []int {
	histo := make([]int, len(buckets)+1, len(buckets)+1)
	for _, s := range samples {
		counted := false
		for i := range buckets {
			if s < buckets[i] {
				histo[i]++
				counted = true
				break
			}
		}
		if !counted {
			histo[len(histo)-1]++
		}
	}
	return histo
}

func makeBuckets(size int, max, start time.Duration) []time.Duration {
	interval := time.Duration(math.Floor(float64(max-start)/float64(size) + 0.5))
	size = int((max - start) / interval)
	buckets := make([]time.Duration, size, size)
	v := start
	for i := 0; i < size; i++ {
		v += interval
		buckets[i] = v
	}
	return buckets
}

var fractionalBlocks = []rune(" \u258F\u258E\u258D\u258C\u258B\u258A\u2589\u2588")

func printHistogram(w io.Writer, buckets []time.Duration, histo []int, maxWidth int) {
	var max int
	for _, h := range histo {
		if h > max {
			max = h
		}
	}
	scale := float64(maxWidth) / float64(max)
	for i := range histo {
		if i < len(buckets) {
			io.WriteString(w, fmt.Sprintf(" %-8s: ", buckets[i].String()))
		} else {
			io.WriteString(w, fmt.Sprintf(">%-8s: ", buckets[len(buckets)-1].String()))
		}
		blocksCalc := scale * float64(histo[i])
		blocks := int(math.Floor(blocksCalc))
		blocksFractionValue := int(math.Ceil((blocksCalc - float64(blocks)) / 0.125))
		for i = 0; i < blocks; i++ {
			io.WriteString(w, "\u2588")
		}
		io.WriteString(w, string(fractionalBlocks[blocksFractionValue]))
		io.WriteString(w, "\n")
	}
}

func timeAction(f func()) time.Duration {
	start := time.Now()
	f()
	return time.Since(start)
}
