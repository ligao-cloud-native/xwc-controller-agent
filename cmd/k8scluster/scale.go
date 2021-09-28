package k8scluster

import (
	"github.com/spf13/cobra"
)

func NewScaleCmd() *cobra.Command{
	cmd := &cobra.Command{
		Use:                   "scale",
		Short:                 "add k8s cluster node",
		Long: "add k8s cluster node",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return cmd
}