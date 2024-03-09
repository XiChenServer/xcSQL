package command

import (
	"SQL/internal/database"
	"fmt"
	"github.com/spf13/cobra"
)

var CliDB = &cobra.Command{
	Use:   "DB [name]",
	Short: "Connect to the database",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please provide the database name")
			return
		}
		name := args[0]
		db := database.DBConnect(name)
		// 可以在这里做更多与数据库相关的操作
		fmt.Println("Connected to database:", name)
		defer db.Close() // 确保在程序退出时关闭数据库连接
	},
}
