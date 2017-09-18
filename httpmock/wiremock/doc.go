// This package provides a mechanism for importing wiremock Json API Mappings into
// mockery/httpmock.  The implementation, at this time, is incomplete.  It handles
// file, string, and json based response bodies.  It does not handle any templating.
// It can handle request matching conditions that include "EqualsTo", "Contains",
// "Matches" and "DoesNotMatch" for Headers and Query Parameters.  It also supports
// all types of url matching.
package wiremock
