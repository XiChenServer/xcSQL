package test

import (
	"fmt"
	"testing"
)

func TestCode(t *testing.T) {
	str := "烫"
	fmt.Printf("字符串 \"%s\" 的二进制表示为: ", str)
	for _, char := range str {
		fmt.Printf("%08b ", char)
	}
	fmt.Println()
}
