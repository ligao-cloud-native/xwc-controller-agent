package k8scluster

import (
	"github.com/ligao-cloud-native/xwc-controller-agent/cluster"
	"github.com/ligao-cloud-native/xwc-controller-agent/cmd/utils"
	"github.com/ligao-cloud-native/xwc-controller-agent/pkg/types"
	"github.com/spf13/cobra"
)


type Reset struct {
	operator
	Action   types.OperatorType
}

func NewResetOption() *Reset {
	reset := Reset{Action: types.ResetOperatorType}

	reset.Agent = utils.GetAgent()
	reset.Nodes =	utils.GetNodes()

	reset.Callback = NewCallback(reset.Agent, reset.Nodes)

	return &reset
}

func NewResetCmd() *cobra.Command{
	reset := NewResetOption()

	cmd := &cobra.Command{
		Use:                   string(types.ResetOperatorType),
		Short:                 "revert any changes to this cluster",
		Long: "revert any changes to this cluster",
		Run: func(cmd *cobra.Command, args []string) {
			reset.Start()
		},
	}

	return cmd
}

func (r *Reset) Start() {
	actuator := cluster.NewActuator(r.Agent, r.Nodes, utils.Env)
	if actuator.K8sResetNodes() {
		r.Callback.Command = cluster.CallBackCmd(types.ResetOperatorType, true, utils.Env.PkgServerUrl)
	} else {
		r.Callback.Command = cluster.CallBackCmd(types.ResetOperatorType, false, utils.Env.PkgServerUrl)
	}

	r.Callback.Execute()

}