package main

import (
	"SQL/internal/database"
	"SQL/logs"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
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
		// 循环接受用户输入的命令
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Enter command: ")
			scanner.Scan()
			input := scanner.Text()
			//if input == "exit" {
			//	fmt.Println("Exiting...")
			//	db.Close()
			//	os.Exit(0)
			//}
			err := handleCommand(input, db)
			if err != nil {
				fmt.Println("Error:", err)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 这里什么都不做，所有的逻辑在 PersistentPreRun 钩子中处理
	},
}

func handleCommand(input string, db *database.XcDB) error {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}
	cmd := parts[0]
	switch cmd {
	case "set":
		if len(parts) != 3 && len(parts) != 4 {
			logs.SugarLogger.Error("Usage: set [key] [value] [ttl]...")
			return fmt.Errorf("Usage: set [key] [value] [ttl]...")

		}
		key := []byte(parts[1])
		value := []byte(parts[2])
		ttl, _ := strconv.Atoi(string(parts[3]))
		logs.SugarLogger.Info(" set [key] [value] [ttl]...")
		return db.Set(key, value, []uint64{uint64(ttl)}...)
	case "get":
		if len(parts) != 2 {
			logs.SugarLogger.Error("Usage: get [key]")
			return fmt.Errorf("Usage: get [key]")
		}
		key := []byte(parts[1])
		value, err := db.Get(key)
		if err != nil {
			return err
		}
		fmt.Println("Value:", string(value))
	case "exit":
		database.DBExit(db)
		//fmt.Println("Exiting...")

		os.Exit(0)
	default:
		return fmt.Errorf("Unknown command: %s", cmd)
	}
	return nil
}

// 添加其他命令...

func main() {
	rootCmd := &cobra.Command{Use: "app"}
	rootCmd.AddCommand(CliDB)
	// 添加其他命令...
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
