package execute

import "github.com/ligao-cloud-native/kubemc/pkg/apis/xwc/v1"

type NodeRole string

const (
	NodeRoleMaster NodeRole = "master"
	NodeRoleWorker NodeRole = "worker"
	NodeRoleEtcd NodeRole = "etcd"
)

type ClusterNodes struct {
	v1.Node
	Role NodeRole
	IsInitNode bool
}
