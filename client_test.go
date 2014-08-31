package netmock_test

import "bytes"
import "errors"
import "fmt"
import "net/http"
import "net/url"
import "strings"
import "testing"
import . "github.com/nowk/go-netmock"

const (
	_OK        = 200
	_NOT_FOUND = 404
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Assert struct{}

func (a Assert) StatusCode(t *testing.T, res *http.Response, exp int) {
	if code := res.StatusCode; code != exp {
		t.Errorf("Expected `StatusCode` %d, got %d", exp, code)
	}
}

func (a Assert) Header(t *testing.T, res interface{}, key string, exp string) {
	var val string
	switch v := res.(type) {
	case *http.Response:
		val = v.Header.Get(key)
	case *http.Request:
		val = v.Header.Get(key)
	}
	if val != exp {
		t.Errorf("Expected `Header` %s to be %s, got %s", key, exp, val)
	}
}

func (a Assert) RequestURL(t *testing.T, req *http.Request, exp string) {
	val := req.URL.String()
	if val != exp {
		t.Errorf("Expected URL %s, got %s", exp, val)
	}
}

func (a Assert) Body(t *testing.T, req *http.Request, exp string) {
	body := req.Body
	defer req.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(body)

	val := buf.String()
	if val != exp {
		t.Errorf("Expected URL %s, got %s", exp, val)
	}
}

var assert = &Assert{}

func TestGet(t *testing.T) {
	mock := NewMock(t)
	mock.Expect("GET", "http://example.com").Respond(_OK, "")

	res, err := mock.HTTPClient.Get("http://example.com")
	check(err)

	assert.StatusCode(t, res, _OK)
	assert.RequestURL(t, res.Request, "http://example.com")
}

func TestHead(t *testing.T) {
	mock := NewMock(t)
	mock.Expect("HEAD", "http://example.com").Respond(_OK, "")

	res, err := mock.HTTPClient.Head("http://example.com")
	check(err)

	assert.StatusCode(t, res, _OK)
	assert.RequestURL(t, res.Request, "http://example.com")
}

func TestPost(t *testing.T) {
	mock := NewMock(t)
	mock.Expect("POST", "http://example.com").Respond(_OK, "")

	body := strings.NewReader("Hello World!")
	res, err := mock.HTTPClient.Post("http://example.com", "", body)
	check(err)

	assert.StatusCode(t, res, _OK)
	assert.RequestURL(t, res.Request, "http://example.com")
	assert.Body(t, res.Request, "Hello World!")
}

func TestPostForm(t *testing.T) {
	mock := NewMock(t)
	mock.Expect("POST", "http://example.com").Respond(_OK, "")

	res, err := mock.HTTPClient.PostForm("http://example.com", url.Values{
		"foo": {"bar"},
	})
	check(err)

	assert.StatusCode(t, res, _OK)
	assert.RequestURL(t, res.Request, "http://example.com")
	assert.Header(t, res.Request, "Content-Type", "application/x-www-form-urlencoded")
	assert.Body(t, res.Request, "foo=bar")
}

func TestDo(t *testing.T) {
	mock := NewMock(t)
	mock.Expect("PUT", "http://example.com").Respond(_NOT_FOUND, "")

	url, err := url.Parse("http://example.com")
	check(err)

	req := &http.Request{
		Method: "PUT",
		URL:    url,
	}

	res, err := mock.HTTPClient.Do(req)
	check(err)

	assert.StatusCode(t, res, _NOT_FOUND)
	assert.RequestURL(t, res.Request, "http://example.com")
}

func TestModifyAfterExpect(t *testing.T) {
	mock := NewMock(t)
	resp, _ := mock.Expect("GET", "http://example.com").Respond(_OK, "")
	resp.StatusCode = _NOT_FOUND
	resp.Header.Add("Content-Type", "application/json")

	res, err := mock.HTTPClient.Get("http://example.com/")
	check(err)

	assert.StatusCode(t, res, _NOT_FOUND)
	assert.Header(t, res, "Content-Type", "application/json")
}

func TestError(t *testing.T) {
	mock := NewMock(t)
	mock.Expect("GET", "http://example.com").Fail(errors.New("Boom!"))

	_, err := mock.HTTPClient.Get("http://example.com")
	if err == nil {
		t.Errorf("Expected an error")
	}

	if er := err.Error(); er != "Boom!" {
		t.Errorf("Expected Boom!, got %s", er)
	}
}

type fTesting struct {
	Errormsg string
}

func (f fTesting) Error(args ...interface{}) {}

func (f *fTesting) Errorf(format string, args ...interface{}) {
	f.Errormsg = fmt.Sprintf(format, args...)
}

func TestCallToUnregisteredURL(t *testing.T) {
	test := &fTesting{}
	mock := NewMock(test)
	mock.Expect("GET", "http://example.com/foo").Respond(_OK, "")
	res, _ := mock.HTTPClient.Get("http://example.com/bar")

	exp := "Called to unmocked URL: [GET] http://example.com/bar"
	if test.Errormsg != exp {
		t.Errorf("Expected %s, got %s", exp, test.Errormsg)
	}

	assert.RequestURL(t, res.Request, "http://example.com/bar")
}

func TestFix2ndExpectIsNotLookedUp(t *testing.T) {
	mock := NewMock(t)
	mock.Expect("GET", "http://example.com/foo").Respond(_NOT_FOUND, "")
	mock.Expect("GET", "http://example.com/bar").Respond(_OK, "")

	res, err := mock.HTTPClient.Get("http://example.com/bar")
	check(err)

	assert.StatusCode(t, res, 200)
}
