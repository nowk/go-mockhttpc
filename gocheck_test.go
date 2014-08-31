package netmock_test

import "testing"
import . "github.com/nowk/go-netmock"
import . "gopkg.in/check.v1"

func TestSuite(t *testing.T) {
	TestingT(t)
}

type tSuite struct {
	mock *Mock
}

func (s *tSuite) SetUpTest(c *C) {
	s.mock = NewMock(c)
	s.mock.Expect("GET", "http://example.com/bad").Respond(404, "")
	s.mock.Expect("GET", "http://example.com/good").Respond(200, "")
}

var _ = Suite(&tSuite{})

func (s *tSuite) TestMock(c *C) {
	res, _ := s.mock.HTTPClient.Get("http://example.com/good")
	if res.StatusCode != 200 {
		c.Errorf("Expected 200, go %d", res.StatusCode)
	}

	res, _ = s.mock.HTTPClient.Get("http://example.com/bad")
	if res.StatusCode != 404 {
		c.Errorf("Expected 404, go %d", res.StatusCode)
	}

	// res, _ = s.mock.HTTPClient.Get("http://example.com/not/mocked")
	// if res.StatusCode != 200 {
	// 	c.Errorf("Expected 200, go %d", res.StatusCode)
	// }
}
