package set

import "fmt"

type Rep int32

func (r Rep) Int32() int32 {
	return int32(r)
}

func NewRep(n int32) (Rep, error) {
	if n <= 0 {
		return Rep(0), fmt.Errorf("invalid rep: %d", n)
	}
	return Rep(n), nil
}
