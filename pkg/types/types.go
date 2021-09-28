package types

import (
	"github.com/ligao-cloud-native/kubemc/pkg/apis/xwc/v1"
)

type Nodes struct {
	Masters []v1.Node `json:"masters,omitempty"`
	Workers []v1.Node `json:"workers,omitempty"`
	Etcd []v1.Node `json:"workers,omitempty"`
}

type Env struct {
	PkgServerUrl string `json:"pkgserver_url,omitempty"`

	CallbackUrl string `json:"callback_url,omitempty"`

	NatsServerUrl string `json:"natserver_url,omitempty"`
	NatsServerToken string `json:"natserver_token,omitempty"`
}