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

func SlienceRise(s []int) {
	s = append(s, 0)
	for i := range s {
		s[i]++
	}

}
func Test_SclicePrint(t *testing.T) {
	s1 := []int{1, 2}
	s2 := s1
	fmt.Println(&s1[0])
	s2 = append(s2, 3)
	SlienceRise(s1)
	fmt.Println(&s1[0])
	SlienceRise(s2)

	fmt.Println(s1, s2)
}
func Test_Extend(t *testing.T) {
	var slice []int
	s1 := append(slice, 1, 2, 3)
	s2 := append(s1, 4)
	fmt.Println(&s1[0])
	fmt.Println(&s2[0])
	fmt.Println(&s1[0] == &s2[0])
}
