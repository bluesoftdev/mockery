package httpmock

import (
	"regexp"
	"strings"
	. "github.com/bluesoftdev/go-http-matchers/predicate"
	. "github.com/bluesoftdev/go-http-matchers/extractor"
)

func BodyXPathEquals(xpath, value string) Predicate {
	return ExtractedValueAccepted(ExtractXPathString(xpath),StringEquals(value))
}

func BodyXPathEqualsIgnoreCase(xpath, value string) Predicate {
	return ExtractedValueAccepted(UpperCaseExtractor(ExtractXPathString(xpath)),StringEquals(strings.ToUpper(value)))
}

func BodyXPathMatches(xpath string, pattern *regexp.Regexp) Predicate {
	return ExtractedValueAccepted(ExtractXPathString(xpath),StringMatches(pattern))
}
