package wiremock

import (
	"github.com/bluesoftdev/go-http-matchers/predicate"
	"github.com/bluesoftdev/mockery/httpmock"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"time"
)

type wireMockValueCondition struct {
	EqualTo         string
	CaseInsensitive *bool
	BinaryEqualTo   string
	Contains        string
	Matches         string
	DoesNotMatch    string
}

type wireMockRequest struct {
	Method          string
	Url             string `json:"url"`
	UrlPattern      string `json:"urlPattern"`
	UrlPath         string `json:"urlPath"`
	UrlPathPattern  string `json:"urlPathPattern"`
	Headers         map[string]wireMockValueCondition
	QueryParameters map[string]wireMockValueCondition
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
	Priority *int
	Request  wireMockRequest
	Response wireMockResponse
}

var mappingFilePattern = regexp.MustCompile("^.*\\.json$")

// WireMockEndpoints takes the dirName and looks for .json files in a subdirectory named "mappings"
// any files named in the mappings are looked for in the __files subdirectory of the base dir name.
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

// WireMockEndpoint takes the name of the base dir the files are expected in and the filename of a
// wiremock .json mapping file.
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
	predicates := make([]predicate.Predicate, 0, 10)
	if wm.Request.Url != "" {
		predicates = append(predicates, predicate.RequestURIEquals(wm.Request.Url))
	} else if wm.Request.UrlPattern != "" {
		predicates = append(predicates, predicate.RequestURIMatches(regexp.MustCompile(wm.Request.UrlPattern)))
	} else if wm.Request.UrlPath != "" {
		predicates = append(predicates, predicate.PathEquals(wm.Request.UrlPath))
	} else if wm.Request.UrlPathPattern != "" {
		predicates = append(predicates, predicate.PathMatches(regexp.MustCompile(wm.Request.UrlPathPattern)))
	}
	if wm.Request.Method != "" {
		predicates = append(predicates, predicate.MethodIs(wm.Request.Method))
	}
	if wm.Request.Headers != nil {
		for header, condition := range wm.Request.Headers {
			if condition.EqualTo != "" {
				if condition.CaseInsensitive != nil && *condition.CaseInsensitive {
					predicates = append(predicates, predicate.HeaderEqualsIgnoreCase(header, condition.EqualTo))
				} else {
					predicates = append(predicates, predicate.HeaderEquals(header, condition.EqualTo))
				}
			} else if condition.Contains != "" {
				if condition.CaseInsensitive != nil && *condition.CaseInsensitive {
					predicates = append(predicates, predicate.HeaderContainsIgnoreCase(header, condition.Contains))
				} else {
					predicates = append(predicates, predicate.HeaderContains(header, condition.Contains))
				}
			} else if condition.Matches != "" {
				predicates = append(predicates, predicate.HeaderMatches(header, regexp.MustCompile(condition.Matches)))
			} else if condition.DoesNotMatch != "" {
				predicates = append(predicates, predicate.Not(predicate.HeaderMatches(header, regexp.MustCompile(condition.Matches))))
			}
		}
	}
	if wm.Request.QueryParameters != nil {
		for query, condition := range wm.Request.QueryParameters {
			if condition.EqualTo != "" {
				if condition.CaseInsensitive != nil && *condition.CaseInsensitive {
					predicates = append(predicates, predicate.QueryParamEqualsIgnoreCase(query, condition.EqualTo))
				} else {
					predicates = append(predicates, predicate.QueryParamEquals(query, condition.EqualTo))
				}
			} else if condition.Contains != "" {
				if condition.CaseInsensitive != nil && *condition.CaseInsensitive {
					predicates = append(predicates, predicate.QueryParamContainsIgnoreCase(query, condition.Contains))
				} else {
					predicates = append(predicates, predicate.QueryParamContains(query, condition.Contains))
				}
			} else if condition.Matches != "" {
				predicates = append(predicates, predicate.QueryParamMatches(query, regexp.MustCompile(condition.Matches)))
			} else if condition.DoesNotMatch != "" {
				predicates = append(predicates, predicate.Not(predicate.QueryParamMatches(query, regexp.MustCompile(condition.Matches))))
			}
		}
	}
	priority := httpmock.DefaultPriority
	if wm.Priority != nil {
		priority = *wm.Priority
	}
	httpmock.EndpointForConditionWithPriority(priority, predicate.And(predicates...), func() {
		if wm.Response.Headers != nil {
			for k, v := range wm.Response.Headers {
				httpmock.Header(k, v.(string))
			}
		}
		if wm.Response.BodyFileName != "" {
			httpmock.LogLocation(fmt.Sprintf("Responding with %s", wm.Response.BodyFileName))
			httpmock.RespondWithFile(wm.Response.Status, dataDirName+"/"+wm.Response.BodyFileName)
		} else if wm.Response.Body != "" {
			httpmock.RespondWithString(wm.Response.Status, wm.Response.Body)
		} else if wm.Response.JsonBody != nil {
			httpmock.RespondWithJson(wm.Response.Status, wm.Response.JsonBody)
		} else {
			httpmock.Respond(wm.Response.Status)
		}
		if wm.Response.DelayDistribution != nil {
			if wm.Response.DelayDistribution.Algorithm == "lognormal" {
				s := wm.Response.DelayDistribution.Sigma
				mean := float64(wm.Response.DelayDistribution.Median) * float64(time.Millisecond) / float64(time.Second)
				u := math.Log(mean)
				stddev := math.Sqrt(math.Exp(2*u+s*s) * (math.Exp(s*s) - 1))
				httpmock.NormalDelay(
					fmt.Sprintf("%dns", int(mean*float64(time.Second))),
					fmt.Sprintf("%dns", int(stddev*float64(time.Second))),
					fmt.Sprintf("%dns", wm.Response.DelayDistribution.Median+int(stddev*5.0*float64(time.Second))))
			} else if wm.Response.DelayDistribution.Algorithm == "uniform" {
				httpmock.UniformDelay(
					fmt.Sprintf("%dms", wm.Response.DelayDistribution.Lower),
					fmt.Sprintf("%dms", wm.Response.DelayDistribution.Upper))
			}
		} else if wm.Response.FixedDelayMilliseconds != nil {
			httpmock.FixedDelay(fmt.Sprintf("%dms", *wm.Response.FixedDelayMilliseconds))
		}
	})
}
