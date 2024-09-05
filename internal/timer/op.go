package timer

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Timer struct {
	levels  int
	refresh []uint32
}

// NewTimer 初始化一个 Timer 实例，指定要管理的 LSM 树层级数。
func NewTimer(levels int) *Timer {
	return &Timer{
		levels:  levels,
		refresh: make([]uint32, levels),
	}
}

// loadConfig 从配置文件中加载定时器设置。
func (t *Timer) loadConfig() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}

	v := viper.New()
	v.SetConfigFile(filepath.Join(path, "config", "config.yaml"))

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return err
	}

	for i := 0; i < t.levels; i++ {
		key := fmt.Sprintf("refresh.level%d", i+1)
		t.refresh[i] = v.GetUint32(key)
	}

	return nil
}

// startCleanup 启动清理任务，并为每个层级启动一个协程。
func (t *Timer) startCleanup(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < t.levels; i++ {
		wg.Add(1)
		go t.cleanupLevel(i, wg)
	}
}

// cleanupLevel 针对指定层级进行清理操作。
func (t *Timer) cleanupLevel(level int, wg *sync.WaitGroup) {
	defer wg.Done()

	interval := time.Duration(t.refresh[level]) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 清理逻辑
			fmt.Printf("Cleaning up LSM tree level %d at interval %d seconds\n", level+1, t.refresh[level])
			// 在此处添加层级的清理逻辑
		}
	}
}

// Run 启动定时器并运行清理任务。
func (t *Timer) Run() error {
	if err := t.loadConfig(); err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go t.startCleanup(&wg)

	wg.Wait() // 等待所有协程完成
	return nil
}
