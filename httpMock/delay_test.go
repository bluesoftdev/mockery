package httpMock

import (
	"github.com/stretchr/testify/assert"
	"github.com/montanaflynn/stats"
	"net/http"
	"testing"
	"time"
)

func TestNormalDelay(t *testing.T) {
	handler := Mockery(func() {
		Endpoint("/foo/bar", func() {
			Method("GET", func() {
				RespondWithFile(200, "./ok.json")
				Header("Content-Type", "application/json")
				NormalDelay("10ms", "1ms", "500ms")
			})
		})
	})
	assert.NotNil(t, handler, "handler is nil")
	if assert.IsType(t, &http.ServeMux{}, handler, "mockery is not an http.ServeMux") {
		serveMux, _ := handler.(*http.ServeMux)
		testReq, err := http.NewRequest("GET", "http://localhost/foo/bar", nil)
		assert.NoError(t, err)
		pathHandler, pattern := serveMux.Handler(testReq)
		assert.NotEmpty(t, pattern, "pattern should not be empty: %s", pattern)
		assert.NotNil(t, pathHandler, "path handler should be defined")
		if assert.IsType(t, &mock{}, pathHandler, "path handler is not a mock") {
			pathMock, _ := pathHandler.(*mock)
			getHandler, ok := pathMock.methods["GET"]
			assert.True(t, ok, "No GET method found")
			if assert.IsType(t, &normalDelay{}, getHandler, "handler is not a normalDelay") {
				getMock, _ := getHandler.(*normalDelay)
				assert.Equal(t, time.Millisecond*10, getMock.mean)
				assert.Equal(t, time.Millisecond, getMock.stdDev)
				assert.Equal(t, time.Millisecond*500, getMock.max)
			}
		}
	}
}

func TestNormalDelay2(t *testing.T) {
	nd := normalDelay{mean: time.Millisecond, stdDev: 10*time.Microsecond, max: 2 * time.Millisecond}
	samples := make([]time.Duration,0,1000)
	for i := 0; i < 1000; i++ {
		start := time.Now()
		nd.Wait()
		end := time.Now()
		samples := append(samples,end.Sub(start))
	}
	population := stats.LoadRawData(samples)

	mean, err := population.Mean()
	assert.NoError(t,err)
	assert.Equal(t, float64(time.Millisecond), mean)

}
