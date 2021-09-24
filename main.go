package main

import (
	"k8s.io/component-base/logs"
)

func main() {
	//command := app.NewK8SClusterCommand()

	logs.InitLogs()
	defer logs.FlushLogs()

	//if err := command.Execute(); err != nil {
	//	os.Exit(1)
	//}
}
