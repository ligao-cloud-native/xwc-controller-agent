package k8scluster

import (
	"github.com/spf13/cobra"
)

func NewResetCmd() *cobra.Command{
	cmd := &cobra.Command{
		Use:                   "reset",
		Short:                 "revert any changes to this cluster",
		Long: "revert any changes to this cluster",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return cmd
}