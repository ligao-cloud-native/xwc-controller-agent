package types

type Command struct {
	InitNodeEtcdCmd string
	EtcdCmd string
	InitNodeMasterCmd string
	MasterCmd string
	WorkerCmd string
}


type CmdExecInfo struct {
	Node Node
	Cmd string
}

type Node struct {
	Ip string
	Uuid string
}

type CmdExecResult struct {
	Host string
	CmdList string
	Success bool
	Result string
}
