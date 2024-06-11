package test_1

import (
	"fmt"
	"testing"
)

func TestFor(t *testing.T) {
	slice := []int{0, 1, 2, 3}
	m := make(map[int]*int)

	for key, val := range slice {
		fmt.Println(&val)
		v := val
		m[key] = &v
	}

	for k, v := range m {
		fmt.Println(k, "->", *v)
	}
}
