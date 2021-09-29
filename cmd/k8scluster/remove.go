package k8scluster

import (
	"github.com/ligao-cloud-native/xwc-controller-agent/cluster"
	"github.com/ligao-cloud-native/xwc-controller-agent/cmd/utils"
	"github.com/ligao-cloud-native/xwc-controller-agent/pkg/types"
	"github.com/spf13/cobra"
	"strings"
)

type Remove struct {
	operator
	Action types.OperatorType
	removeNodes string
}

func NewRemoveOption() *Remove {
	remove := Remove{Action: types.RemoveOperatorType}

	remove.Agent = utils.GetAgent()
	remove.Nodes =	utils.GetNodes()

	remove.Callback = NewCallback(remove.Agent, remove.Nodes)

	return &remove
}

func NewRemoveCmd() *cobra.Command{
	remove := NewRemoveOption()

	cmd := &cobra.Command{
		Use:                   string(types.RemoveOperatorType),
		Short:                 "remove k8s cluster",
		Long: "remove k8s cluster node, as: k8s remove",
		Run: func(cmd *cobra.Command, args []string) {
			remove.Start()
		},
	}

	cmd.PersistentFlags().StringVar(&remove.removeNodes, "nodes", "", "remove cluster node")

	return cmd
}

func (r *Remove) Start() {
	if r.removeNodes != "" {
		var RemovedNodes types.Nodes

		// 可能某个节点既是master又是worker
		for _, removedNode := range strings.Split(r.removeNodes, ",") {
			for _, master := range r.Nodes.Masters {
				if master.IP == removedNode {
					RemovedNodes.Masters = append(RemovedNodes.Masters, master)
				}
			}
			for _, worker := range r.Nodes.Workers {
				if worker.IP == removedNode {
					RemovedNodes.Workers = append(RemovedNodes.Workers, worker)
				}
			}
		}

		actuator := cluster.NewActuator(r.Agent, RemovedNodes, utils.Env)
		if actuator.K8sRemoveNodes() {
			r.Callback.Command = cluster.CallBackCmd(types.RemoveOperatorType, true, utils.Env.PkgServerUrl)
		} else {
			r.Callback.Command = cluster.CallBackCmd(types.RemoveOperatorType, false, utils.Env.PkgServerUrl)
		}

		r.Callback.Execute()

	}
}