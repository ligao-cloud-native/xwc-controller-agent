package k8scluster

import (
	"github.com/ligao-cloud-native/xwc-controller-agent/cluster"
	"github.com/ligao-cloud-native/xwc-controller-agent/cmd/utils"
	"github.com/ligao-cloud-native/xwc-controller-agent/pkg/types"
	"github.com/ligao-cloud-native/xwc-controller-agent/provider"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)


type Install struct {
	Name     string
	Action   string
	Nodes    types.Nodes
	Agent    provider.Agent
	Callback utils.Callback
}

func NewInstallOption() *Install {
	install := Install{Action: "install"}

	install.Agent = utils.GetAgent()
	install.Nodes =	utils.GetNodes()

	install.Callback = utils.NewCallback()

	return &install
}

func NewInstallCmd() *cobra.Command{
	install := NewInstallOption()

	cmd := &cobra.Command{
		Use:                   "install",
		Short:                 "install k8s cluster",
		Long: "install k8s cluster",
		Run: func(cmd *cobra.Command, args []string) {
			install.Start()
		},
	}

	return cmd
}

func (i *Install) Start() {
	if len(i.Nodes.Masters) <= 0 {
		klog.Fatal("k8s cluster master node is empty.")
	}
	actuator := cluster.NewActuator(i.Agent, i.Nodes, utils.Env)

	//设置无密登录
	if actuator.SSHAuth() {

	}

	// 是否使用外部etcd集群，如果使用则需要安装。
	if len(i.Nodes.Etcd) >= 3 {
		// TODO: install etcd cluster
	}

	// 安装k8s集群


}
