package netmock_test

import "regexp"
import "testing"
import . "github.com/nowk/go-netmock"

func TestMatchRegardlessOfTrailingBackSlash(t *testing.T) {
	mock := Mock{}

	for k, v := range map[string]string{
		"http://example.com/":      "http://example.com",
		"http://example.com/posts": "http://example.com/posts/",
	} {
		mock.Expect("GET", k)
		mres := mock.FindMock("GET", v)
		if mres == nil {
			t.Errorf("Expected %s to return a MockRequest for %s", k, v)
		}
	}
}

func TestMatchWithQuery(t *testing.T) {
	url := "http://example.com/foo?token=12345&id=67890"
	mock := Mock{}
	mock.Expect("GET", url)

	for _, v := range []string{
		"http://example.com/foo?token=12345&id=67890",
		"http://example.com/foo?id=67890&token=12345",
	} {
		mres := mock.FindMock("GET", v)
		if mres == nil {
			t.Errorf("Expected %s to return a MockRequest for %s", url, v)
		}
	}

	for _, v := range []string{
		"http://example.com/",
		"http://example.com/foo",
		"http://example.com/foo?token=",
		"http://example.com/foo?token=12345&id=67899",
		"http://example.com/foo?token=12345&id=67899&q=abcde",
	} {
		mres := mock.FindMock("GET", v)
		if mres != nil {
			t.Errorf("Expected %s to not return a MockRequest for %s", url, v)
		}
	}
}

func TestRegisterURLRegex(t *testing.T) {
	reg := regexp.MustCompile(`http:\/\/example\.com\/foo`)

	mock := Mock{}
	mock.Expect("GET", reg)

	for _, v := range []string{
		"http://example.com/foo",
		"http://example.com/foo?token=",
		"http://example.com/foo?token=12345&id=67899",
		"http://example.com/foo?id=67899&q=abcde&token=12345&",
	} {
		mres := mock.FindMock("GET", v)
		if mres == nil {
			t.Errorf("Expected %s to return a MockRequest", v)
		}
	}

	for _, v := range []string{
		"http://example.com",
		"http://example.com/fo",
		"http://example.com/bar",
	} {
		mres := mock.FindMock("GET", v)
		if mres != nil {
			t.Errorf("Expected %s to not return a MockRequest", v)
		}
	}
}

func TestRegisterMustBeUnique(t *testing.T) {
	defer func() {
		exp := "GET http://example.com is already registered"
		err := recover()
		if exp != err {
			t.Errorf("Expected panic %s, got %s", exp, err)
		}
	}()

	mock := Mock{}
	mock.Expect("GET", "http://example.com/")
	mock.Expect("GET", "http://example.com")
}

func TestClearRegistry(t *testing.T) {
	mock := Mock{}
	mock.Expect("GET", "http://example.com/1")
	mock.Expect("GET", "http://example.com/2")

	mock.Reset()
	if n := len(mock.MockRequests); n != 0 {
		t.Errorf("Expected 0, got %d", n)
	}

	mock.Expect("GET", "http://example.com/1")
	mock.Expect("GET", "http://example.com/2")
	if n := len(mock.MockRequests); n != 2 {
		t.Errorf("Expected 2, got %d", n)
	}
}
