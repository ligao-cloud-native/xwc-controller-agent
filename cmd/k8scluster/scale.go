package k8scluster

import (
	"github.com/ligao-cloud-native/xwc-controller-agent/cluster"
	"github.com/ligao-cloud-native/xwc-controller-agent/cmd/utils"
	"github.com/ligao-cloud-native/xwc-controller-agent/pkg/types"
	"github.com/spf13/cobra"
)

type Scale struct {
	operator
	Action   types.OperatorType
}

func NewScaleOption() *Scale {
	reset := Scale{Action: types.ScaleOperatorType}

	reset.Agent = utils.GetAgent()
	reset.Nodes = utils.GetNodes()

	reset.Callback = NewCallback(reset.Agent, reset.Nodes)

	return &reset

}

func NewScaleCmd() *cobra.Command{
	scale := NewScaleOption()

	cmd := &cobra.Command{
		Use:                   string(types.ScaleOperatorType),
		Short:                 "add k8s cluster node",
		Long: "add k8s cluster node",
		Run: func(cmd *cobra.Command, args []string) {
			scale.Start()
		},
	}

	return cmd
}

func (s *Scale) Start() {
	actuator := cluster.NewActuator(s.Agent, s.Nodes, utils.Env)

	//设置无密登录
	if actuator.SSHAuth() {
		actuator.K8sScaleNodes()
	}

	s.Callback.Command = cluster.CallBackCmd(types.ScaleOperatorType, s.Nodes, utils.Env.PkgServerUrl)
	s.Callback.Execute()




}