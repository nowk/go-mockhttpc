package netmock_test

import "net/url"
import "strings"
import "testing"
import . "github.com/nowk/go-netmock"

func TestSpiesCalledTimes(t *testing.T) {
	mock := NewMock(t)
	_, mres := mock.Expect("GET", "http://example.com").Respond(_OK, "")

	mock.HTTPClient.Get("http://example.com")
	mres.Called(1).Times()

	mock.HTTPClient.Get("http://example.com")
	mock.HTTPClient.Get("http://example.com")
	mock.HTTPClient.Get("http://example.com")
	mres.Called(3).Times()
}

func TestSpiesRequestBody(t *testing.T) {
	mock := NewMock(t)
	_, mres := mock.Expect("POST", "http://example.com").Respond(_OK, "")

	body := strings.NewReader("Hello World!")
	mock.HTTPClient.Post("http://example.com", "", body)

	mres.Body().Equals("Hello World!")
}

func TestSpiesRequestHeader(t *testing.T) {
	mock := NewMock(t)
	_, mres := mock.Expect("POST", "http://example.com").Respond(_OK, "")

	mock.HTTPClient.PostForm("http://example.com", url.Values{
		"foo": {"Bar"},
	})

	mres.Header("Content-Type").Equals("application/x-www-form-urlencoded")
}

func TestSpiesForm(t *testing.T) {
	mock := NewMock(t)
	_, mres := mock.Expect("POST", "http://example.com?foo=bar").Respond(_OK, "")

	mock.HTTPClient.PostForm("http://example.com?foo=bar", url.Values{
		"baz": {"qux"},
	})

	mres.Form("baz").Equals("qux")
	mres.Form("foo").Equals("bar")
}

func TestSpiesPostForm(t *testing.T) {
	mock := NewMock(t)
	_, mres := mock.Expect("POST", "http://example.com?foo=bar").Respond(_OK, "")

	mock.HTTPClient.PostForm("http://example.com?foo=bar", url.Values{
		"baz": {"qux"},
	})

	mres.PostForm("baz").Equals("qux")
	mres.PostForm("foo").Equals("")
}
