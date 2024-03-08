package command

import (
	"SQL/internal/database"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "mydb",
	Short: "A simple database CLI",
}

var (
	db  *database.XcDB
	ttl uint64
)

//
//func main() {
//	db = database.NewXcDB()
//
//	rootCmd.AddCommand(setCmd)
//	rootCmd.AddCommand(getCmd)
//	rootCmd.AddCommand(strlenCmd)
//	rootCmd.AddCommand(appendCmd)
//
//	if err := rootCmd.Execute(); err != nil {
//		fmt.Println("Error:", err)
//		os.Exit(1)
//	}
//}

var setCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a key-value pair in the database",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := []byte(args[0])
		value := []byte(args[1])
		err := db.Set(key, value, ttl)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Set successfully.")
	},
}

var getCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get the value for a given key from the database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := []byte(args[0])
		value, err := db.Get(key)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Value:", string(value))
	},
}

var strlenCmd = &cobra.Command{
	Use:   "strlen [key]",
	Short: "Get the length of the value for a given key from the database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := []byte(args[0])
		length, err := db.Strlen(key)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Length:", length)
	},
}

var appendCmd = &cobra.Command{
	Use:   "append [key] [value]",
	Short: "Append value to the existing value for a given key in the database",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := []byte(args[0])
		value := []byte(args[1])
		err := db.Append(key, value)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Append successfully.")
	},
}

func init() {
	rootCmd.PersistentFlags().Uint64Var(&ttl, "ttl", 0, "Time to live (in seconds)")
}
