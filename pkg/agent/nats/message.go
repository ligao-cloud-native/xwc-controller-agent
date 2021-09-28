package nats

type CmdExecInfo struct {
	TaskInfo
	StdOut string
	StdErr string
	ExitCode int
}

type TaskInfo struct {
	TaskId string
	Status string
}
