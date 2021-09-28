package cluster

import (
	"fmt"
	"github.com/ligao-cloud-native/xwc-controller-agent/cluster/execute"
	"github.com/ligao-cloud-native/xwc-controller-agent/pkg/types"
	"github.com/ligao-cloud-native/xwc-controller-agent/provider"
	"k8s.io/klog/v2"
	"time"
)

type Actuator struct {
	Agent provider.Agent
	Nodes []execute.ClusterNodes
	Env types.Env
}


func NewActuator(agent provider.Agent, nodes types.Nodes, env types.Env) *Actuator {
	act := &Actuator{Agent: agent, Env: env}
	act.Nodes = allClusterNodes(nodes)

	return act
}


func (a Actuator) SSHAuth() bool {
	klog.Infof("Step1: SSHAuthorized")

	concurrent := &execute.Concurrent{
		Agent: a.Agent,
		Nodes: a.Nodes,
		ConcurrentNum: len(a.Nodes),
		Timeout: 480 * time.Second,
	}

	//定义不同类型节点上执行的命令
	cmd := fmt.Sprintf(sshAuthCmd, a.Env.PkgServerUrl)
	concurrent.Command = types.Command{
		InitNodeEtcdCmd: cmd,
		EtcdCmd: cmd,
		InitNodeMasterCmd: cmd,
		MasterCmd: cmd,
		WorkerCmd: cmd,
	}

	return concurrent.Execute()
}

func (a Actuator) DownloadEtcdPackage() bool {}



func allClusterNodes(nodes types.Nodes) []execute.ClusterNodes {
	var clusterNodes []execute.ClusterNodes

	for i, etcd := range nodes.Etcd {
		initNode := false
		if i == 0 {
			initNode = true
		}

		clusterNodes = append(clusterNodes, execute.ClusterNodes{
			Node:etcd, Role: execute.NodeRoleEtcd, IsInitNode:initNode})
	}
	for i, master := range nodes.Masters {
		initNode := false
		if i == 0 {
			initNode = true
		}

		clusterNodes = append(clusterNodes, execute.ClusterNodes{
			Node:master, Role: execute.NodeRoleMaster, IsInitNode:initNode})
	}
	for _, worker := range nodes.Workers {
		clusterNodes = append(clusterNodes, execute.ClusterNodes{
			Node:worker, Role: execute.NodeRoleWorker, IsInitNode:false})
	}

	return clusterNodes

}
