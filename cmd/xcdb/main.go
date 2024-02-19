package main

import "SQL/logs"
import _ "SQL/logs"

func main() {
	logs.InitLogger()
	defer logs.SugarLogger.Sync()
	logs.SugarLogger.Debug("This is a test")
}
