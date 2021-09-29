package types

import (
	"github.com/ligao-cloud-native/kubemc/pkg/apis/xwc/v1"
)

type OperatorType string

const (
	InstallOperatorType = "install"
	RemoveOperatorType = "remove"
	ScaleOperatorType = "scale"
	ResetOperatorType = "reset"
)


type Nodes struct {
	Masters []v1.Node `json:"masters,omitempty"`
	Workers []v1.Node `json:"workers,omitempty"`
	Etcd []v1.Node `json:"workers,omitempty"`
}

type Env struct {
	PkgServerUrl string `json:"pkgserver_url,omitempty"`

	CallbackUrl string `json:"callback_url,omitempty"`

	// NatsServerUrl 为vmserver服务暴漏的rest server
	// http://vmserver.kmc-nats:8000
	NatsServerUrl string `json:"natserver_url,omitempty"`
	NatsServerToken string `json:"natserver_token,omitempty"`
}