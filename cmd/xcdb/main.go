package main

import (
	"SQL/internal/command"
	"SQL/internal/database"
	"SQL/internal/wal"
	"SQL/logs"
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var dbName string

var CliDB = &cobra.Command{
	Use:   "DB [name]",
	Short: "Connect to the database",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please provide the database name")
			os.Exit(1)
		}
		dbName = args[0]
		db := database.DBConnect(dbName)
		logs.SugarLogger.Infof("Connected to database:", dbName)
		fmt.Println("Connected to database:", dbName)
		// 启动守护协程
		go daemon(db.Wal, db)
		// 处理用户命令
		handleCommands(db, db.Wal)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 这里什么都不做，所有的逻辑在 PersistentPreRun 钩子中处理
	},
}

func main() {
	rootCmd := &cobra.Command{Use: "app"}
	rootCmd.AddCommand(CliDB)
	// 添加其他命令...
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func handleCommands(db *database.XcDB, wal *wal.WAL) {
	// 循环接受用户输入的命令
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter command: ")
		scanner.Scan()
		input := scanner.Text()
		if strings.ToLower(input) == "exit" { // 将输入转换为小写，以便匹配退出命令
			fmt.Println("Exiting...")
			database.DBExit(db)
			os.Exit(0)
		}
		// 将命令记录到 WAL 中
		if err := wal.Write(input); err != nil {
			fmt.Println("Error writing to WAL:", err)
			continue
		}
		// 执行命令
		if err := handleCommand(input, db); err != nil {
			fmt.Println("Error:", err)
		}
	}
}

// 启动一个守护协程
func daemon(wal *wal.WAL, db *database.XcDB) {
	for {
		// 从 WAL 日志文件中读取新的命令
		cmd, err := wal.ReadNextCommand()
		if err != nil {
			log.Println("Error reading command from WAL:", err)
			continue
		}
		fmt.Println(cmd)
		//// 执行命令
		//if err := executeCommand(cmd, db); err != nil {
		//	log.Println("Error executing command:", err)
		//	continue
		//}
		//
		//// 更新数据库状态
		//if err := updateDatabaseState(db); err != nil {
		//	log.Println("Error updating database state:", err)
		//	continue
		//}
	}
}

func handleCommand(input string, db *database.XcDB) error {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}
	cmd := strings.ToLower(parts[0]) // 将命令转换为小写
	switch cmd {
	case "set":
		return command.HandleSetCommand(parts, db)
	case "get":
		return command.HandleGetCommand(parts, db)
	case "append":
		return command.HandleAppendCommand(parts, db)
	case "strlen":
		return command.HandleStrlenCommand(parts, db)
	// 添加其他命令的处理逻辑...
	case "exit":
		// 在主函数中处理退出逻辑，这里不再需要处理
		return nil
	default:
		return fmt.Errorf("Unknown command: %s", cmd)
	}
}
