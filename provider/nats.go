package provider

import (
	"fmt"
	v1 "github.com/ligao-cloud-native/kubemc/pkg/apis/xwc/v1"
	"github.com/ligao-cloud-native/xwc-controller-agent/pkg/agent/nats"
	"github.com/ligao-cloud-native/xwc-controller-agent/pkg/types"
	"time"
)

type NatsAgent struct {
	option ExecCmdOpt
	client *nats.Client

}

type ExecCmdOpt struct {
	Retry int
	RetrySpan time.Duration
}

func NewNatsAgent(url, token string, opt ExecCmdOpt) *NatsAgent {
	client := nats.NewConfig(url, token).BuildClientOrDie()

	return &NatsAgent{
		option: opt,
		client: client,
	}
}

func (n *NatsAgent) CommandExecute(node v1.Node, cmd string) (result string, success bool){
	execInfo := types.CmdExecInfo{
		Node: types.Node{Ip:node.IP, Uuid: node.NodeID},
		Cmd:cmd,
	}
	res, err := n.client.ExecCmdAndGetInfo(execInfo, n.option.Retry, n.option.RetrySpan)
	if err != nil {
		success = false
		result = fmt.Sprintf("%v", res)
		return result, success
	}

	if res.ExitCode != 0 {
		success = false
		result = fmt.Sprintf("%v", res)
		return result, success
	}

	return res.StdOut, true


}

