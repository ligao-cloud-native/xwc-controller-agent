package cmd

import (
	"github.com/ligao-cloud-native/xwc-controller-agent/cmd/k8scluster"
	"github.com/spf13/cobra"

)

func NewK8SClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "k8s",
		Short: "k8s client",
		Long: "k8s cluster manager, to create/scale/reduce/reset cluster.",
	}

	cmd.AddCommand(k8scluster.NewInstallCmd())
	cmd.AddCommand(k8scluster.NewRemoveCmd())
	cmd.AddCommand(k8scluster.NewScaleCmd())
	cmd.AddCommand(k8scluster.NewResetCmd())

	return cmd
}