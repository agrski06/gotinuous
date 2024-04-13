package main

import "github.com/agrski06/gotinuous/internal"

func main() {
	tool := internal.InitTool()
	tool.InitRepository()
	tool.ExecStages()
}
