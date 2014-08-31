package netmock

import "bytes"

// AssertCountCalled allows assertions for how many times a MockRequest
// gets called
type AssertCountCalled struct {
	mock *MockRequest
	exp  int
}

func (a *AssertCountCalled) Times() {
	v := a.mock.CallCount
	if v != a.exp {
		a.mock.t.Errorf("Expected [%s] %s to be called %d times, but was called %d times",
			a.mock.Method,
			a.mock.URL,
			a.exp,
			v)
	}

	a.mock.CallCount = 0 // reset count
}

func (m *MockRequest) Called(exp int) *AssertCountCalled {
	return &AssertCountCalled{
		mock: m,
		exp:  exp,
	}
}

// AssertBody allows assertions for the Request body associated with a
// MockRequest
type AssertBody struct {
	mock *MockRequest
}

func (a AssertBody) Equals(exp string) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(a.mock.Response.Request.Body)

	v := buf.String()
	if v != exp {
		a.mock.t.Errorf("Expected Request Body %s, got %s", exp, v)
	}
}

func (m *MockRequest) Body() *AssertBody {
	return &AssertBody{
		mock: m,
	}
}

// AssertHeader allows assertions for the Request header associated with a
// MockRequest
type AssertHeader struct {
	mock *MockRequest
	key  string
}

func (a AssertHeader) Equals(exp string) {
	header := a.mock.Response.Request.Header

	v := header.Get(a.key)
	if v != exp {
		a.mock.t.Errorf("Expected Request Header %s to be %s, got %s", a.key, exp, v)
	}
}

func (m *MockRequest) Header(key string) *AssertHeader {
	return &AssertHeader{
		mock: m,
		key:  key,
	}
}

// AssertForm which allows you to assert Form inputs via ParseForm()
type AssertForm struct {
	mock *MockRequest
	key  string
}

func (a AssertForm) Equals(exp string) {
	a.mock.Response.Request.ParseForm()
	values := a.mock.Response.Request.Form

	v := values.Get(a.key)
	if v != exp {
		a.mock.t.Errorf("Expected Form %s to be %s, got %s", a.key, exp, v)
	}
}

func (m *MockRequest) Form(key string) *AssertForm {
	return &AssertForm{
		mock: m,
		key:  key,
	}
}

// AssertPostForm which allows you to assert PostForm inputs via ParseForm()
type AssertPostForm struct {
	mock *MockRequest
	key  string
}

func (a AssertPostForm) Equals(exp string) {
	a.mock.Response.Request.ParseForm()
	values := a.mock.Response.Request.PostForm

	v := values.Get(a.key)
	if v != exp {
		a.mock.t.Errorf("Expected PostForm %s to be %s, got %s", a.key, exp, v)
	}
}

func (m *MockRequest) PostForm(key string) *AssertPostForm {
	return &AssertPostForm{
		mock: m,
		key:  key,
	}
}
