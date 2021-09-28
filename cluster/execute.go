package cluster

import (
	"fmt"
	"github.com/ligao-cloud-native/kubemc/pkg/apis/xwc/v1"
	"github.com/ligao-cloud-native/xwc-controller-agent/cluster/execute"
	"github.com/ligao-cloud-native/xwc-controller-agent/pkg/types"
	"github.com/ligao-cloud-native/xwc-controller-agent/provider"
	"k8s.io/klog/v2"
	"time"
)

const (
	DownloadPkgTimeout = 90 * time.Second
)

type Actuator struct {
	Agent provider.Agent
	Nodes []execute.ClusterNodes
	Env types.Env
	Concurrent *execute.Concurrent
}


func NewActuator(agent provider.Agent, nodes types.Nodes, env types.Env) *Actuator {
	act := &Actuator{Agent: agent, Env: env}
	act.Nodes = allClusterNodes(nodes)
	act.Concurrent = &execute.Concurrent{
		Agent: agent,
		Nodes: act.Nodes,
		ConcurrentNum: len(act.Nodes),
		Timeout: DownloadPkgTimeout,
	}

	return act
}


func (a *Actuator) SSHAuth() bool {
	klog.Infof("Step1: SSHAuthorized")

	//定义不同类型节点上执行的命令
	cmd := fmt.Sprintf(sshAuthCmd, a.Env.PkgServerUrl)
	a.Concurrent.Command = types.Command{
		InitNodeEtcdCmd: cmd,
		EtcdCmd: cmd,
		InitNodeMasterCmd: cmd,
		MasterCmd: cmd,
		WorkerCmd: cmd,
	}

	return a.Concurrent.Execute()
}

func (a *Actuator) EtcdDownloadPkg() bool {
	klog.Infof("Step2: EtcdDownloadPackage")

	cleanCmd := fmt.Sprintf(etcdCleanCmd, a.Env.PkgServerUrl)
	pkgCmd := fmt.Sprintf(etcdDownlodPkgCmd, a.Env.PkgServerUrl, a.Env.PkgServerUrl)
	cmd := fmt.Sprintf("%s && %s", cleanCmd, pkgCmd)
	a.Concurrent.Command = types.Command{
		InitNodeEtcdCmd: cmd,
		EtcdCmd: cmd,
	}

	return a.Concurrent.Execute()
}


func (a Actuator) EtcdInstall() bool {
	klog.Infof("Step3: EtcdInstall")

	cmdParam := EtcdCmdParam(a.Nodes)
	cmd := fmt.Sprintf(etcdInstallCmd, a.Env.PkgServerUrl, a.Env.PkgServerUrl, cmdParam)
	a.Concurrent.Command = types.Command{
		InitNodeEtcdCmd: cmd,
	}

	return a.Concurrent.Execute()
}


func (a *Actuator) K8sDownloadPkgAndInstallInitMaster() bool {
	klog.Infof("Step3: K8sDownloadPkgAndInstall")

	etcdCmdParam := EtcdCmdParam(a.Nodes)
	k8sCmdParam := CniK8sverRuntimeCmdParam(a.Nodes)

	//第一个master节点上执行的命令
	initNodeDownloadPkgCmd := fmt.Sprintf(k8sInstallCmd, a.Env.PkgServerUrl, a.Env.PkgServerUrl,
		etcdCmdParam, InitNodeDownloadPkgCmd(a.Nodes) + k8sCmdParam )
	initNodeLabelCmd := fmt.Sprintf(k8sMasterLableCmd, "nodeName", "labelKey=labelValue")
	initNodeCmd := fmt.Sprintf(k8sResetCmd, a.Env.PkgServerUrl) + " && " +
		initNodeDownloadPkgCmd + initNodeLabelCmd

	//其他节点上执行的命令
	otherNodeDownloadPkgCmd := fmt.Sprintf(k8sPrepareCmd, a.Env.PkgServerUrl, a.Env.PkgServerUrl,
		k8sCmdParam)
	otherNodeCmd := fmt.Sprintf(k8sResetCmd, a.Env.PkgServerUrl) + " && " + otherNodeDownloadPkgCmd


	a.Concurrent.Command = types.Command{
		InitNodeEtcdCmd: initNodeCmd,
		MasterCmd: otherNodeCmd,
		WorkerCmd: otherNodeCmd,
	}


	return a.Concurrent.Execute()

}


func (a *Actuator) K8sJoinOtherMastersAndWorkers() bool {
	// 获取join token
	initMasterNode := v1.Node{}
	for _, node := range a.Nodes {
		if node.Role == execute.NodeRoleMaster && node.IsInitNode {
			initMasterNode = node.Node
			break
		}
	}
	result, success := a.Agent.CommandExecute(initMasterNode, k8sJoinTokenCmd)
	if success {
		//TODO:处理result
	}

	// join
	a.Concurrent.Command = types.Command{
		MasterCmd: fmt.Sprintf(k8sJoinMasterCmd, result),
		WorkerCmd: fmt.Sprintf(k8sJoinWorkerCmd, "labels", "params"),
	}

	return a.Concurrent.Execute()

}


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
