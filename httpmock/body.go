package httpmock

import (
	"github.com/bluesoftdev/go-http-matchers/extractor"
	"github.com/bluesoftdev/go-http-matchers/predicate"
	"regexp"
	"strings"
)

// BodyXPathEquals checks to see if the result of the xpath expression, matches the string supplied in the 'value'
// parameter.
func BodyXPathEquals(xpath, value string) predicate.Predicate {
	return predicate.ExtractedValueAccepted(extractor.ExtractXPathString(xpath), predicate.StringEquals(value))
}

// BodyXPathEqualsIgnoreCase similar to BodyXPathEquals but ignores case when comparing the strings.
func BodyXPathEqualsIgnoreCase(xpath, value string) predicate.Predicate {
	return predicate.ExtractedValueAccepted(extractor.UpperCaseExtractor(extractor.ExtractXPathString(xpath)),
		predicate.StringEquals(strings.ToUpper(value)))
}

// BodyXPathMatches checks to see if the result of the xpath expression, matches the regular expression given in the
// 'pattern' parameter.
func BodyXPathMatches(xpath string, pattern *regexp.Regexp) predicate.Predicate {
	return predicate.ExtractedValueAccepted(extractor.ExtractXPathString(xpath), predicate.StringMatches(pattern))
}
