package main

import (
	"SQL/internal/database"
	"SQL/logs"
	"bufio"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
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
		handleCommands(db)
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

func handleCommands(db *database.XcDB) {
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
		db.Wal.Write(input)
		err := handleCommand(input, db)
		if err != nil {
			fmt.Println("Error:", err)
		}
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
		return handleSetCommand(parts, db)
	case "get":
		return handleGetCommand(parts, db)
	case "append":
		return handleAppendCommand(parts, db)
	case "strlen":
		return handleStrlenCommand(parts, db)
	// 添加其他命令的处理逻辑...
	case "exit":
		// 在主函数中处理退出逻辑，这里不再需要处理
		return nil
	default:
		return fmt.Errorf("Unknown command: %s", cmd)
	}
}

func handleSetCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 3 && len(parts) != 4 {
		return errors.New("Usage: set [key] [value] [ttl]...")
	}
	key := []byte(parts[1])
	value := []byte(parts[2])
	ttl := uint64(0)
	if len(parts) == 4 {
		ttl, _ = strconv.ParseUint(parts[3], 10, 64)
	}

	return db.Set(key, value, ttl)
}

func handleGetCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 2 {
		return errors.New("Usage: get [key]")
	}
	key := []byte(parts[1])
	value, err := db.Get(key)
	if err != nil {
		return err
	}
	fmt.Println("Value:", string(value))
	return nil
}

func handleAppendCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 3 {
		return errors.New("Usage: append [key] [value]")
	}
	key := []byte(parts[1])
	value := []byte(parts[2])
	return db.Append(key, value)
}

func handleStrlenCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 2 {
		return errors.New("Usage: strlen [key]")
	}
	key := []byte(parts[1])
	valueLen, err := db.Strlen(key)
	if err != nil {
		return err
	}
	fmt.Println("Value length:", valueLen)
	return nil
}
