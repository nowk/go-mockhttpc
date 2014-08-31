# go-netmock

[![Build Status](https://travis-ci.org/nowk/go-netmock.svg?branch=master)](https://travis-ci.org/nowk/go-netmock)
[![GoDoc](https://godoc.org/gopkg.in/nowk/go-netmock.v0?status.svg)](http://godoc.org/gopkg.in/nowk/go-netmock.v0)

`http.Client` mocking when you can't control the URL

## Install

    go get gopkg.in/nowk/go-netmock.v0

## Example

You use a `http.DefaultClient` mapped interface.

    type SomeApi struct {
      HTTPClient HTTPClient
    }

    func (s SomeApi) GetSomeApiRequest() (*http.Response, error) {
      return s.HTTPClient.Get("http://example.com")
    }

Pass your mock in your tests.

    func TestRequest(t *testing.T) {
      mock := netmock.NewMock(t)
      mock.Expect("GET", "http://example.com").Respond(200, "Hello World!")

      api := &SomeApi{mock.HTTPClient}
      res, _ := api.GetSomeApiRequest()

      if code := res.StatusCode; code != 200 {
        t.Errorf("Expected 200, got %d", code)
      }
    }

---

#### Register routes

    mock.Expect("GET", "http://example.com").Respond(200, "Hello World!")

With a `regexp`

    reg := regexp.MustCompile(`http:\/\/example\.com\/foo`)
    mock.Expect("GET", reg).Respond(200, "Hello World!")

---

#### Modify the response

    res, _ := mock.Expect("GET", "http://example.com").Respond(200, `{"foo": "bar"}`)
    res.Header.Add("Content-Type", "application/json")

---

#### Spy on the Request

    _, mr := mock.Expect("GET", "http://example.com").Respond(200, "Hello World!")

Number of times called

    mr.Called(1).Times()

Request Body

    mr.Body().Equals("Hello World!")

Request Header

    mr.Header("Content-Type").Equals("application/json")

Form values

    mr.Form("foo").Equals("bar")

PostForm values

    mr.PostForm("baz").Equals("qux")

## License

MIT
