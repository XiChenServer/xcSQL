package main

import (
	"SQL/Refactor/interface/command"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// 初始化命令和 WAL
	cmd := command.NewStringCommand()
	rootCmd := command.RegisterCommands(cmd)
	rootCmd.Execute()
	// 创建一个 Scanner 读取标准输入
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("欢迎使用 KV 数据库 CLI。请输入命令（set key value [ttl], del key, get key）：")

	for {
		fmt.Print("请输入命令: ")
		if !scanner.Scan() {
			fmt.Println("读取输入失败")
			break
		}
		input := scanner.Text()

		// 解析用户输入的命令
		args := strings.Fields(input)
		if len(args) == 0 {
			fmt.Println("无效的命令")
			continue
		}
		rootCmd.SetArgs(args)
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
		}
	}
}
