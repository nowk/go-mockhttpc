package netmock

import "fmt"
import "net/url"
import "reflect"
import "regexp"
import . "github.com/nowk/go-httpclienti"

const VERSION = "0.0.1"

// Mock
type Mock struct {
	t            Testing
	HTTPClient   HTTPClient
	MockRequests []*MockRequest
}

func NewMock(t Testing) (mock *Mock) {
	mock = &Mock{
		t: t,
	}
	mock.HTTPClient = NewClient(mock)

	return
}

func (m *Mock) Reset() {
	m.MockRequests = nil
}

// Register registers a MockRequest to a specific http method + url (or regexp)
func (m *Mock) Expect(method string, urlStr interface{}) *MockRequest {
	mres := NewMockRequest(method, urlStr)
	mres.t = m.t
	m.register(mres)

	return mres
}

// register appends MockRequests
func (m *Mock) register(mres *MockRequest) *MockRequest {
	if dup := m.checkDup(mres); dup {
		panic(fmt.Sprintf("%s %s is already registered", mres.Method, mres.URL))
	}

	m.MockRequests = append(m.MockRequests, mres)

	return mres
}

// checkDup checks for a duplicate MockRequest
func (m Mock) checkDup(mres *MockRequest) bool {
	for _, r := range m.MockRequests {
		if r.Method == mres.Method && r.URL == mres.URL {
			return true
		}
	}

	return false
}

// FindMock finds a MockRequest by url and method
func (m Mock) FindMock(method, urlStr string) *MockRequest {
	for _, r := range m.MockRequests {
		if r.Method == method && urlsMatch(r.URL, urlize(urlStr)) {
			return r
		}
	}

	return nil
}

// urlize takes a string and runs it through url.Parse
func urlize(urlStr string) *url.URL {
	str := removeTrailingbackslash(urlStr)
	u, err := url.Parse(str)
	if err != nil {
		panic(fmt.Sprintf("Url Parse error: %s", err))
	}

	return u
}

// urlsMatch matches an interface (regexp|string) against a url.URL
func urlsMatch(a interface{}, b *url.URL) bool {
	switch v := a.(type) {
	case *regexp.Regexp:
		return v.MatchString(b.String())
	case string:
		url := urlize(v)
		return reflect.DeepEqual(url.Query(), b.Query()) &&
			url.Scheme == b.Scheme &&
			url.Host == b.Host &&
			url.Path == b.Path
	}

	return false
}
