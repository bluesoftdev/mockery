package httpMock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type wireMockRequest struct {
	Method         string
	Url            string `json:"url"`
	UrlPath        string `json:"urlPath"`
	UrlPathPattern string `json:"urlPattern"`
}

type wireMockDelayDistribution struct {
	Algorithm string `json:"type"`
	Median    int
	Sigma     float64
	Lower     int
	Upper     int
}

type wireMockResponse struct {
	Status        int
	StatusMessage string
	Headers       map[string]interface{}

	Body         string
	JsonBody     interface{}
	Base64Body   string
	BodyFileName string

	FixedDelayMilliseconds *int
	DelayDistribution      *wireMockDelayDistribution
}

type wireMock struct {
	Request  wireMockRequest
	Response wireMockResponse
}

var mappingFilePattern = regexp.MustCompile("^.*\\.json$")

func WireMockEndpoints(dirName string) {
	mappingDir := dirName + string(os.PathSeparator) + "mappings"
	dataDir := dirName + string(os.PathSeparator) + "__files"
	mappingFiles, err := ioutil.ReadDir(mappingDir)
	if err != nil {
		panic("Error trying to list mapping files:" + err.Error())
	}
	for _, mappingFile := range mappingFiles {
		if mappingFilePattern.MatchString(mappingFile.Name()) {
			WireMockEndpoint(dataDir, mappingDir+"/"+mappingFile.Name())
		}
	}
}

func WireMockEndpoint(dataDirName, fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		panic("Error opening mapping file: " + err.Error())
	}
	d := json.NewDecoder(f)
	var wm wireMock
	err = d.Decode(&wm)
	if err != nil {
		panic("Error parsing mapping file: " + err.Error())
	}
	predicates := make([]Predicate, 0, 10)
	if wm.Request.Url != "" {
		predicates = append(predicates, PathEquals(wm.Request.Url))
	} else if wm.Request.UrlPath != "" {
		predicates = append(predicates, PathStartsWith(wm.Request.UrlPath))
	} else if wm.Request.UrlPathPattern != "" {
		predicates = append(predicates, PathMatches(regexp.MustCompile(wm.Request.UrlPathPattern)))
	}
	if wm.Request.Method != "" {
		predicates = append(predicates, MethodIs(wm.Request.Method))
	}
	EndpointForCondition(And(predicates...), func() {
		if wm.Response.Headers != nil {
			for k, v := range wm.Response.Headers {
				Header(k, v.(string))
			}
		}
		if wm.Response.BodyFileName != "" {
			RespondWithFile(wm.Response.Status, dataDirName+"/"+wm.Response.BodyFileName)
		} else if wm.Response.Body != "" {
			RespondWithString(wm.Response.Status, wm.Response.Body)
		} else if wm.Response.JsonBody != nil {
			RespondWithJson(wm.Response.Status, wm.Response.JsonBody)
		} else {
			Respond(wm.Response.Status)
		}
		if wm.Response.DelayDistribution != nil {
			if wm.Response.DelayDistribution.Algorithm == "lognormal" {
				stddev := wm.Response.DelayDistribution.Sigma * float64(wm.Response.DelayDistribution.Median)
				NormalDelay(
					fmt.Sprintf("%dms", wm.Response.DelayDistribution.Median),
					fmt.Sprintf("%dms", stddev),
					fmt.Sprintf("%dms", wm.Response.DelayDistribution.Median+int(stddev*5.0)))
			} else if wm.Response.DelayDistribution.Algorithm == "uniform" {
				UniformDelay(
					fmt.Sprintf("%dms", wm.Response.DelayDistribution.Lower),
					fmt.Sprintf("%dms", wm.Response.DelayDistribution.Upper))
			}
		} else if wm.Response.FixedDelayMilliseconds != nil {
			FixedDelay(fmt.Sprintf("%dms", *wm.Response.FixedDelayMilliseconds))
		}
	})
}
