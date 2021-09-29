package cluster

import (
	"bytes"
	"fmt"
	"github.com/ligao-cloud-native/xwc-controller-agent/cluster/execute"
	"github.com/ligao-cloud-native/xwc-controller-agent/pkg/types"
	"os"
	"strconv"
	"strings"
)

var ClusterName = os.Getenv("PWC_NAME")

var (
	cmdStatusOk = " 1"
	cmdStatusFailed = " 2"


	logInstallCmd = "2>&1 |tee /opt/install.log"
	logInstallAppendCmd = "2>&1 |tee /opt/install.log"
	logResetCmd = "2>&1 | tee /opt/reset.log"
	logJoinCmd = "2>&1 | tee /opt/join.log"
	logRemoveCmd = "2>&1 | tee /opt/remove.log"
	logScaleCmd = "2>&1 | tee /opt/scale.log"

	env = "export PKG_SERVER=%s && export CLUSTER_NAME=" + ClusterName
	envClusterName = "export CLUSTERNAME=" + ClusterName
	envKubeConfig = "export KUBECONFIG=/etc/kubernetes/admin.conf && " + envClusterName

	sshAuthCmd = "curl -k %s/sshd.tar.gz | tar zxv -C /tmp/ && chmod +x /tmp/sshd/exec.sh && sh /tmp/sshd/exec.sh"

	etcdCleanCmd = "curl -k %s/etcd-installer/clean-etcd.sh | sh -s -- -y " + logResetCmd
	etcdDownlodPkgCmd = env + " && curl -k %s/install-etcd | sh " + logInstallCmd
	etcdInstallCmd = env + " && curl -k %s/install-etcd | %s sh -s init " + logInstallAppendCmd

	k8sInstallCmd = env + "&& curl -k %s/install | %s sh -s init %s " + logInstallCmd
	k8sScaleNodeCmd = env + "&& curl -k %s/install | sh -s worker --init %s " + logScaleCmd
	k8sPrepareCmd = env + "&& curl -k %s/install | sh -s prepare %s " + logInstallAppendCmd
	k8sResetCmd = "curl -k %s/pks-installer/k8s/reset-node.sh | sh -s -- -y " + logResetCmd
	k8sMasterLableCmd = envKubeConfig + " && /usr/local/bin/kubectl lable node %s %s"

	k8sJoinTokenCmd = `/usr/local/bin/kubeadm token create --print-join-command |awk '{print $3" "$5" "$7}'`
	k8sJoinMasterCmd = envClusterName + " && cd /tmp/pks-installer/k8s/ && ./install.sh master %s " + logJoinCmd
	k8sJoinWorkerCmd = envClusterName + " && cd /tmp/pks-installer/k8s/ && ./install.sh worker %s %s" + logJoinCmd

	k8sMasterDeleteCmd = " && /usr/local/bin/kubectl delete node %s"
	k8sMasterDrainCmd = " && /usr/local/bin/kubectl drain %s --delete-local-data --force --ignore-daemonsets"
	k8sMasterCleanCmd = envKubeConfig + k8sMasterDrainCmd + k8sMasterDeleteCmd
	k8sNodeCleanCmd = "curl -k %s/pks-installer/k8s/reset-node.sh | sh -s -- -y " + logRemoveCmd

	callbackCmd = "curl -s -k %s/pks-installer/k8s/callback.sh |sh -s "
	callbackInstallCmd = callbackCmd + "install " + ClusterName + " %s"
	callbackResetCmd = callbackCmd + "reset " + ClusterName
	callbackRemoveCmd = callbackCmd + "remove " + ClusterName
	callbackScaleCmd = callbackCmd + "scale " + ClusterName


)

func EtcdCmdParam(node []execute.ClusterNodes) string {
	cmdParam := ""
	for _, node := range node {
		if node.Role == execute.NodeRoleEtcd {
			i := 0
			cmdParam = cmdParam + " HOST" + strconv.Itoa(i) + "=" + node.IP
			i++
		}
	}

	return cmdParam
}

func CniK8sverRuntimeCmdParam(node []execute.ClusterNodes) string {
	// network
	cmd := " --cni" + "networkType" + ":" + "networkVersion"
	cmd = cmd + " --cni-extra-args \"" + "cniExtraArgsKey cniExtraArgsValue" + "\""

	//k8s versiom
	cmd = cmd + " --k8s-version " + "k8sVersion"

	// runtime
	cmd = cmd + " -rt " + "runtimeType" + "-rv " + "runtimeVersion"
	cmd = cmd + " --runtime-extra-args \"" + "runtimeExtraArgsKey runtimeExtraArgsValue" + "\""

	return cmd
}


func InitNodeDownloadPkgCmd(nodes []execute.ClusterNodes) string {
	cmd := ""

	// network calico_ipip=always
	cmd = cmd + " --calico-ipip"

	// network kube proxy mode
	cmd = cmd + "--kube-proxy-mode" + "kubeProxyMode"

	// loadbance
	cmd = cmd + "-l " + "loadbanceIP" + " -p" + "loadbancePort"

	// master ip
	for _, node := range nodes {
		if node.Role == execute.NodeRoleMaster {
			cmd = cmd + " -m " + node.IP
		}
	}

	// cert sans
	cmd = cmd + " -san " + "CertSan1"  + " -san " + "CertSan2"

	// addon
	cmd = cmd + "-ad " + "addonType1" + "/" + "addon1" + "-ad " + "addonType2" + "/" + "addon2"

	return cmd

}

func CallBackCmd(opt types.OperatorType, optParam interface{}, pkgServer string) string {
	switch opt {
	case types.InstallOperatorType:
		if nodeObj, ok := optParam.(types.Nodes); ok {
			var nodeIp bytes.Buffer
			for _, master := range nodeObj.Masters {
				nodeIp.WriteString(master.IP + "#")
			}
			for _, worker := range nodeObj.Workers {
				nodeIp.WriteString(worker.IP + "#")
			}

			return fmt.Sprintf(callbackInstallCmd, pkgServer, strings.TrimRight(nodeIp.String(), "#"))
		}
		return fmt.Sprintf(callbackInstallCmd, pkgServer)

	case types.ScaleOperatorType:
		if nodeObj, ok := optParam.(types.Nodes); ok {
			var nodeIp bytes.Buffer
			for _, worker := range nodeObj.Workers {
				nodeIp.WriteString(worker.IP + "#")
			}

			return fmt.Sprintf(callbackScaleCmd, pkgServer, strings.TrimRight(nodeIp.String(), "#"))
		}
		return fmt.Sprintf(callbackInstallCmd, pkgServer)
	case types.ResetOperatorType:
		cmdParam := cmdStatusFailed
		if boolOjb, ok := optParam.(bool); ok && boolOjb {
			cmdParam = cmdStatusOk
		}

		return fmt.Sprintf(callbackResetCmd + cmdParam, pkgServer)

	case types.RemoveOperatorType:
		cmdParam := cmdStatusFailed
		if boolOjb, ok := optParam.(bool); ok && boolOjb {
			cmdParam = cmdStatusOk
		}

		return fmt.Sprintf(callbackRemoveCmd + cmdParam, pkgServer)
	}



	return ""
}

