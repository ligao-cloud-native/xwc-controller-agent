package executer

import (
	"github.com/ligao-cloud-native/kubemc/pkg/apis/xwc/v1"
	"github.com/ligao-cloud-native/xwc-controller-agent/pkg/types"
	"github.com/ligao-cloud-native/xwc-controller-agent/provider"
	"k8s.io/klog/v2"

	"sync"
	"time"
)

type Concurrent struct {
	Agent provider.Agent
	ConcurrentNum int
	Timeout time.Duration
	Nodes []ClusterNodes
	Command types.Command
}


func (c *Concurrent) Execute() bool {
	statistics := make(chan bool)

	execResultCh := make(chan types.CmdExecResult, c.ConcurrentNum)
	finishedCh := make(chan bool)

	go c.execute(execResultCh, finishedCh)

	go func() {
		execSuccess := true
		select {
		case <- time.After(c.Timeout):
			execSuccess = false
			klog.Errorf("timeout")
		case <- finishedCh:
			for result := range execResultCh {
				if !result.Success {
					execSuccess = false
					klog.Errorf("exec command error, host:%s, cmd:%s",result.Host, result.CmdList)
				}
			}
		}
		statistics <- execSuccess
	}()

	return <-statistics
}


func (c *Concurrent) execute(execResultCh chan types.CmdExecResult, finishedCh chan bool) {
	var wg sync.WaitGroup

	for _, node := range c.Nodes {
		wg.Add(1)

		result := types.CmdExecResult{Host: node.IP}
		go func(host ClusterNodes) {
			defer wg.Done()

			switch host.Role {
			case NodeRoleEtcd:
				//第一个节点作为etcd集群的master节点
				if node.IsInitNode {
					result.CmdList = c.Command.InitNodeEtcdCmd
				} else {
					result.CmdList = c.Command.EtcdCmd
				}
				result.Result, result.Success = c.procExec(host.Node, result.CmdList)
			case NodeRoleMaster:
				//第一个节点作为k8s集群的第一个master节点
				if node.IsInitNode {
					result.CmdList = c.Command.InitNodeMasterCmd
				} else {
					result.CmdList = c.Command.MasterCmd
				}
				result.Result, result.Success =c.procExec(host.Node, result.CmdList)
			case NodeRoleWorker:
				result.CmdList = c.Command.WorkerCmd
				result.Result, result.Success = c.procExec(host.Node, result.CmdList)
			}

			execResultCh <- result

		}(node)
	}

	wg.Wait()
	finishedCh <- true
}


func (c *Concurrent) procExec(node v1.Node, cmd string) (string, bool){
	agent := c.Agent
	return agent.CommandExecute(node, cmd)
}
