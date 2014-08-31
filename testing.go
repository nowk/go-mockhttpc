package netmock

type Testing interface {
	Error(...interface{})
	Errorf(string, ...interface{})
}
