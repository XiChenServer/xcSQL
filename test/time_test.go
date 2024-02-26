package test

import (
	"log"
	"testing"
	"time"
)

func TestGetNowTime(T *testing.T) {
	// 获取当前时间
	currentTime := time.Now()

	log.Println(currentTime)
}
