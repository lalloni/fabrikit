package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReflectionShallowCopy(t *testing.T) {
	a := assert.New(t)
	s1 := &struct {
		Name   string
		Age    int
		Phones []string
		Xs     []struct {
			A string
			B int
		}
	}{Name: "pepe", Age: 20, Phones: []string{"1", "2"}, Xs: []struct {
		A string
		B int
	}{{A: "x1", B: 1}, {A: "x2", B: 2}}}
	s2 := reflectionShallowCopy(s1)
	a.EqualValues(s1, s2)
	t.Logf("s1: %[1]p %+[1]v", s1)
	t.Logf("s2: %[1]p %+[1]v", s2)
}
