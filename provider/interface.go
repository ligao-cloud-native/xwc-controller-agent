package provider

import v1 "github.com/ligao-cloud-native/kubemc/pkg/apis/xwc/v1"

type Agent interface {
	CommandExecute(node v1.Node, cmd string) (string, bool)
}