package netmock

import "io/ioutil"
import "net/http"
import "regexp"
import "strings"

// MockRequest represents a request to be mocked
type MockRequest struct {
	URL      interface{}
	Method   string
	Response *http.Response
	Err      error

	t         Testing
	CallCount int
}

func NewMockRequest(method string, urlStr interface{}) (mr *MockRequest) {
	mr = &MockRequest{
		Method: method,
	}
	mr.SetURL(urlStr)

	return
}

// Respond sets the response for the mocked request
func (m *MockRequest) Respond(statusCode int, body string) (*http.Response, *MockRequest) {
	r := strings.NewReader(body)
	res := &http.Response{
		StatusCode: statusCode,
		Body:       ioutil.NopCloser(r),
		Header:     http.Header{},
	}

	m.Response = res

	return res, m
}

// Fail sets the Error for the mocked request. This is a go error and not to be
// confused with 4xx/5xx server errors
func (m *MockRequest) Fail(err error) *MockRequest {
	m.Err = err

	return m
}

// SetURL sets the URL, but removes trailing backslash if url is a string
func (m *MockRequest) SetURL(urlStr interface{}) {
	m.URL = urlStr
	v, ok := urlStr.(string)
	if ok {
		m.URL = removeTrailingbackslash(v)
	}
}

// Increment increments the CallCount value
func (m *MockRequest) Increment() int {
	m.CallCount++

	return m.CallCount
}

// removeTrailingbackslash removes the trailing slash
func removeTrailingbackslash(urlStr string) string {
	reg := regexp.MustCompile(`\/$`)
	str := reg.ReplaceAllString(urlStr, "")

	return str
}
