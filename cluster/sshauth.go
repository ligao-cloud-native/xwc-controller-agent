package cluster

import (
	"fmt"
	"github.com/ligao-cloud-native/xwc-controller-agent/cluster/executer"
	"github.com/ligao-cloud-native/xwc-controller-agent/pkg/types"
	"github.com/ligao-cloud-native/xwc-controller-agent/provider"
	"k8s.io/klog/v2"
	"time"
)

func SSHAuth(agent provider.Agent, nodes types.Nodes, env types.Env) {
	klog.Infof("Step1: SSHAuthorized")

	concurrent := &executer.Concurrent{
		Agent: agent,
		Timeout: 480 * time.Second,
		ConcurrentNum: len(nodes.Masters) + len(nodes.Workers) + len(nodes.Etcd),

	}
	concurrent.Nodes = allClusterNodes(nodes)

	//定义不同类型节点上执行的命令
	cmd := fmt.Sprintf(sshAuthCmd, env.PkgServerUrl)
	concurrent.Command = types.Command{
		InitNodeEtcdCmd: cmd,
		EtcdCmd: cmd,
		InitNodeMasterCmd: cmd,
		MasterCmd: cmd,
		WorkerCmd: cmd,
	}

	concurrent.Execute()
}

func allClusterNodes(nodes types.Nodes) []executer.ClusterNodes {
	var clusterNodes []executer.ClusterNodes

	for i, etcd := range nodes.Etcd {
		initNode := false
		if i == 0 {
			initNode = true
		}

		clusterNodes = append(clusterNodes, executer.ClusterNodes{
			Node:etcd, Role:executer.NodeRoleEtcd, IsInitNode:initNode})
	}
	for i, master := range nodes.Masters {
		initNode := false
		if i == 0 {
			initNode = true
		}

		clusterNodes = append(clusterNodes, executer.ClusterNodes{
			Node:master, Role:executer.NodeRoleMaster, IsInitNode:initNode})
	}
	for _, worker := range nodes.Workers {
		clusterNodes = append(clusterNodes, executer.ClusterNodes{
			Node:worker, Role:executer.NodeRoleWorker, IsInitNode:false})
	}

	return clusterNodes

}
