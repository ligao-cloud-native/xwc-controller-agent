package k8scluster

import (
	"github.com/ligao-cloud-native/xwc-controller-agent/cluster"
	"github.com/ligao-cloud-native/xwc-controller-agent/cmd/utils"
	"github.com/ligao-cloud-native/xwc-controller-agent/pkg/types"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)


type Install struct {
	operator
	Action   types.OperatorType
}

func NewInstallOption() *Install {
	install := Install{Action: types.InstallOperatorType}

	install.Agent = utils.GetAgent()
	install.Nodes =	utils.GetNodes()

	install.Callback = NewCallback(install.Agent, install.Nodes)

	return &install
}

func NewInstallCmd() *cobra.Command{
	install := NewInstallOption()

	cmd := &cobra.Command{
		Use:                   string(types.InstallOperatorType),
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
	if !actuator.SSHAuth() {
		return
	}

	// 是否使用外部etcd集群，如果使用则需要安装。
	if len(i.Nodes.Etcd) >= 3 {
		if !(actuator.EtcdDownloadPkg() && actuator.EtcdInstall()) {
			return
		}
	}

	// 安装k8s集群
	if actuator.K8sJoinOtherMastersAndWorkers() && actuator.K8sJoinOtherMastersAndWorkers() {
		//TODO: master node schedule

		// callback execute
		i.Callback.Command = cluster.CallBackCmd(types.InstallOperatorType, i.Nodes, utils.Env.PkgServerUrl)
		i.Callback.Execute()


	}

}
