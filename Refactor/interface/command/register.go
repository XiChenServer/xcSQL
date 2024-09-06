package command

import (
	"SQL/Refactor/interface/wal"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
)

// RegisterCommands 注册所有的命令
func RegisterCommands(c *StringCommand) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "kvdb",
		Short: "一个简单的键值数据库 CLI",
	}

	// 创建 WAL 文件
	wa, err := wal.CreateWalFile("test.txt", wal.Every)
	if err != nil {
		fmt.Println("创建 WAL 文件失败:", err)
		return nil
	}
	// 设置命令
	var setCmd = &cobra.Command{
		Use:   "set [key] [value] [ttl(optional)]",
		Short: "设置键值对，可以选择性指定过期时间",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			value := args[1]
			ttl, _ := cmd.Flags().GetInt("ttl")
			duration := time.Duration(ttl) * time.Second

			// 将命令写入 WAL 文件
			walCommand := fmt.Sprintf("set %s %s %d", key, value, ttl)
			if err := wa.AddCommandToFile(walCommand); err != nil {
				fmt.Println("写入 WAL 文件失败:", err)
				return
			}

			// 执行设置操作
			if err := c.Set(key, value, duration); err != nil {
				fmt.Println("设置失败:", err)
				return
			}
			fmt.Println("设置成功")
		},
	}
	setCmd.Flags().IntP("ttl", "t", 0, "过期时间（秒）")

	// 删除命令
	var delCmd = &cobra.Command{
		Use:   "del [key]",
		Short: "删除键值对",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]

			// 将命令写入 WAL 文件
			walCommand := fmt.Sprintf("del %s", key)
			if err := wa.AddCommandToFile(walCommand); err != nil {
				fmt.Println("写入 WAL 文件失败:", err)
				return
			}

			// 执行删除操作
			if err := c.Del(key); err != nil {
				fmt.Println("删除失败:", err)
				return
			}
			fmt.Println("删除成功")
		},
	}

	// 获取命令
	var getCmd = &cobra.Command{
		Use:   "get [key]",
		Short: "获取键值对",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			value, err := c.Get(key)
			if err != nil {
				fmt.Println("获取失败:", err)
				return
			}
			fmt.Println("值:", value)
		},
	}

	// 退出命令
	var exitCmd = &cobra.Command{
		Use:   "exit",
		Short: "退出 CLI 并关闭 WAL 文件",
		Run: func(cmd *cobra.Command, args []string) {
			// 关闭 WAL 文件
			err := wa.Close()
			if err != nil {
				fmt.Println("关闭 WAL 文件失败:", err)
			} else {
				fmt.Println("WAL 文件关闭成功，退出程序")
			}
			// 退出程序
			// 这里可以选择调用 os.Exit() 直接退出程序
			os.Exit(0)
		},
	}

	rootCmd.AddCommand(setCmd, delCmd, getCmd, exitCmd)
	for {
		select {
		case v, ok := <-wa.CommandChan:
			if !ok {
				return rootCmd
			}
			args := strings.Fields(v)
			//fmt.Println(args)
			rootCmd.SetArgs(args)
			if err := rootCmd.Execute(); err != nil {
				panic(err)
			}
		}
	}
}

// 启动wal
func WalStart() {

}
