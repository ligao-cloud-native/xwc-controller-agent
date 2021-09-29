package k8scluster

import (
	"github.com/ligao-cloud-native/xwc-controller-agent/pkg/types"
	"github.com/ligao-cloud-native/xwc-controller-agent/provider"
)



type Operator interface {
	Start()
}

type operator struct {
	Nodes    types.Nodes
	Agent    provider.Agent
	Callback *Callback
}


