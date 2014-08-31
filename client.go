package netmock

import "log"
import "io"
import "net/http"
import "net/url"
import "strings"
import "github.com/nowk/go-methods"

// Client is a mock struct that maps to http.DefaultClient
type Client struct {
	mock *Mock
}

// NewClient returns a client that impelments the HTTPClient interface which
// maps to the net/http's http.Client
func NewClient(mock *Mock) (client *Client) {
	client = &Client{
		mock: mock,
	}

	return
}

// responseTo looks up the MockRequest and returns or calls t.Error if a
// MockRequest is not found
func (c Client) responseTo(method, urlStr string, req *http.Request) (*http.Response, error) {
	mr := c.mock.FindMock(method, urlStr)
	if mr != nil {
		// ensure Response, when returning an mocked Error, there is no response
		// avail
		if res := mr.Response; res != nil {
			res.Request = req
		}

		mr.Increment()

		return mr.Response, mr.Err
	}

	// return a blank response with the original request for unmocked URLs
	blankResp := &http.Response{
		Request: req,
	}
	c.mock.t.Errorf("Called to unmocked URL: [%s] %s", method, urlStr)

	return blankResp, nil
}

// maps methods in the http.Client

func (c Client) Do(req *http.Request) (*http.Response, error) {
	return c.responseTo(req.Method, req.URL.String(), req)
}

func (c Client) Get(urlStr string) (*http.Response, error) {
	req, err := http.NewRequest(methods.GET, urlStr, nil)
	check(err)

	return c.responseTo(methods.GET, urlStr, req)
}

func (c Client) Head(urlStr string) (*http.Response, error) {
	req, err := http.NewRequest(methods.HEAD, urlStr, nil)
	check(err)

	return c.responseTo(methods.HEAD, urlStr, req)
}

func (c Client) Post(urlStr string, btype string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(methods.POST, urlStr, body)
	check(err)
	req.Header.Set("Content-Type", btype)

	return c.responseTo(methods.POST, urlStr, req)
}

func (c *Client) PostForm(urlStr string, data url.Values) (*http.Response, error) {
	return c.Post(urlStr, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func check(err error) {
	if err != nil {
		log.Panicf("mock error: %s", err)
	}
}
