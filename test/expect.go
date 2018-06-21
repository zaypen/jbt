package test

import "testing"

type I struct {
	*testing.T
}

func (it *I) Should(should string) *Expecting {
	return &Expecting{it, should}
}

type Expecting struct {
	*I
	should string
}

func (e *Expecting) Expect(expected interface{}) *Expected {
	return &Expected{e,expected}
}

type Expected struct {
	*Expecting
	expected interface{}
}

func (e *Expected) ToBe(actual interface{}) {
	if e.expected != actual {
		e.Errorf("%s. Expected: %v, Actual: %v", e.should, e.expected, actual)
	}
}

func It(t *testing.T) *I {
	return &I{t}
}