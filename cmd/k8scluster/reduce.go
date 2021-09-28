package k8scluster

import (
	"github.com/spf13/cobra"
)

func NewReduceCmd() *cobra.Command{
	cmd := &cobra.Command{
		Use:                   "reduce",
		Short:                 "remove k8s cluster node",
		Long: "remove k8s cluster node, as: k8s reduce 192.168.1.1,192.168.1.2",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return cmd
}