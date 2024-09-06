package command

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// 测试 Set 方法
func TestSet(t *testing.T) {
	cmd := NewStringCommand()

	// 测试没有过期时间的设置
	err := cmd.Set("key1", "value1", 0)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	value, err := cmd.Get("key1")
	if err != nil || value != "value1" {
		t.Errorf("Get failed: %v, got %v", err, value)
	}

	// 测试带有过期时间的设置
	err = cmd.Set("key2", "value2", 2*time.Second)
	if err != nil {
		t.Errorf("Set with ttl failed: %v", err)
	}

	// 等待过期时间
	time.Sleep(3 * time.Second)

	value, err = cmd.Get("key2")
	if err == nil || value != "" {
		t.Errorf("Get after ttl failed: %v, got %v", err, value)
	}
}

// 测试 Del 方法
func TestDel(t *testing.T) {
	cmd := NewStringCommand()
	rootCmd := RegisterCommands(cmd)

	err := cmd.Set("key1", "value1", 0)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	err = cmd.Del("key1")
	if err != nil {
		t.Errorf("Del failed: %v", err)
	}

	value, err := cmd.Get("key1")
	if err == nil || value != "" {
		t.Errorf("Get after Del failed: %v, got %v", err, value)
	}
	err = rootCmd.Execute()
	if err != nil {
		return
	}
}

// 测试 RegisterCommands
func TestRegisterCommands(t *testing.T) {
	cmd := NewStringCommand()
	rootCmd := RegisterCommands(cmd)

	var buf bytes.Buffer
	rootCmd.SetOutput(&buf)

	// 模拟 "set" 命令
	rootCmd.SetArgs([]string{"set", "key1", "value1", "10"})
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	//// 模拟 "get" 命令
	//buf.Reset()
	rootCmd.SetArgs([]string{"get", "key1"})
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("Execute failed: %v", err)
	}
	//output = buf.String()
	//if !contains(output, "值: value1") {
	//	t.Errorf("Expected output to contain '值: value1', got %s", output)
	//}
	//
	// 模拟 "del" 命令
	rootCmd.SetArgs([]string{"del", "key1"})
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("Execute failed: %v", err)
	}
	//output = buf.String()
	//if !contains(output, "删除成功") {
	//	t.Errorf("Expected output to contain '删除成功', got %s", output)
	//}
	//
	//// 模拟 "get" 命令，检查删除是否生效
	//buf.Reset()
	//rootCmd.SetArgs([]string{"get", "key1"})
	//if err := rootCmd.Execute(); err != nil {
	//	t.Errorf("Execute failed: %v", err)
	//}
	//output = buf.String()
	//if !contains(output, "获取失败: 键 key1 不存在") {
	//	t.Errorf("Expected output to contain '获取失败: 键 key1 不存在', got %s", output)
	//}
}

// 辅助函数：检查字符串是否包含子串
func contains(str, substr string) bool {
	return strings.Contains(str, substr)
}
