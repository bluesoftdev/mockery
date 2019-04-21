package httpmock

import (
	"github.com/bluesoftdev/go-http-matchers/predicate"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"testing"
)

var bodyTests = []struct {
	Name string
	Pred predicate.Predicate
	ExpectedResult bool
} {
	{"Equals Match", BodyXPathEquals("/snafu/foo", "bar"), true},
	{"Equals No Match", BodyXPathEquals("/snafu/foo", "Bar"), false},
	{"EqualsIgnoreCase Match", BodyXPathEqualsIgnoreCase("/snafu/foo", "Bar"), true},
	{"EqualsIgnoreCase No Match", BodyXPathEqualsIgnoreCase("/snafu/foo", "Baz"), false},
	{"Matches Match", BodyXPathMatches("/snafu/foo", regexp.MustCompile("b[aeiou]r")), true},
	{"Matches No Match", BodyXPathMatches("/snafu/foo", regexp.MustCompile("b[aeiou]z")), false},
}

func TestBodyXPath(t *testing.T) {
	for _, tst := range bodyTests {
		t.Run(tst.Name, func(t *testing.T) {
			testURL, _ := url.ParseRequestURI("http://localhost/foo")
			body, err := os.Open("testdata/response.xml")
			if assert.NoError(t, err) {
				request := &http.Request{
					Method: "GET",
					Header: http.Header{},
					URL:    testURL,
					Body:   body,
				}
				assert.Equal(t, tst.ExpectedResult, tst.Pred.Accept(request))
			}
		})
	}
}
