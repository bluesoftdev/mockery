package httpmock

import (
	"regexp"
	"strings"
	. "github.com/bluesoftdev/go-http-matchers/predicate"
	. "github.com/bluesoftdev/go-http-matchers/extractor"
)

// Checks to see if the result of the xpath expression, matches the string supplied in the 'value' parameter.
func BodyXPathEquals(xpath, value string) Predicate {
	return ExtractedValueAccepted(ExtractXPathString(xpath),StringEquals(value))
}

// Similar to BodyXPathEquals but ignores case when comparing the strings.
func BodyXPathEqualsIgnoreCase(xpath, value string) Predicate {
	return ExtractedValueAccepted(UpperCaseExtractor(ExtractXPathString(xpath)),StringEquals(strings.ToUpper(value)))
}

// Checks to see if the result of the xpath expression, matches the regular expression given in the 'pattern' parameter.
func BodyXPathMatches(xpath string, pattern *regexp.Regexp) Predicate {
	return ExtractedValueAccepted(ExtractXPathString(xpath),StringMatches(pattern))
}
