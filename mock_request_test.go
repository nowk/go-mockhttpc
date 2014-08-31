package netmock_test

import "testing"
import . "github.com/nowk/go-netmock"

func TestIncrement(t *testing.T) {
	mres := MockRequest{}
	mres.Increment()
	mres.Increment()
	mres.Increment()

	exp := 3
	if v := mres.CallCount; v != exp {
		t.Errorf("Expected %d, got %d", exp, v)
	}
}
