package cluster

import (
	"github.com/ligao-cloud-native/xwc-controller-agent/cluster/execute"
	"os"
	"strconv"
)

var ClusterName = os.Getenv("PWC_NAME")

var (
	logInstallCmd = "2>&1 |tee /opt/install.log"
	logInstallAppendCmd = "2>&1 |tee /opt/install.log"
	logResetCmd = "2>&1 | tee /opt/reset.log"
	logJoinCmd = "2>&1 | tee /opt/join.log"

	env = "export PKG_SERVER=%s && export CLUSTER_NAME=" + ClusterName
	envKubeConfig = "export KUBECONFIG=/etc/kubernetes/admin.conf && export CLUSTERNAME=" + ClusterName

	sshAuthCmd = "curl -k %s/sshd.tar.gz | tar zxv -C /tmp/ && chmod +x /tmp/sshd/exec.sh && sh /tmp/sshd/exec.sh"

	etcdCleanCmd = "curl -k %s/etcd-installer/clean-etcd.sh | sh -s -- -y " + logResetCmd
	etcdDownlodPkgCmd = env + " && curl -k %s/install-etcd | sh " + logInstallCmd
	etcdInstallCmd = env + " && curl -k %s/install-etcd | %s sh -s init " + logInstallAppendCmd

	k8sInstallCmd = env + "&& curl -k %s/install | %s sh -s init %s " + logInstallCmd
	k8sPrepareCmd = env + "&& curl -k %s/install | sh -s prepare %s " + logInstallAppendCmd
	k8sResetCmd = "curl -k %s/pks-installer/k8s/reset-node.sh | sh -s -- -y " + logResetCmd
	k8sMasterLableCmd = envKubeConfig + " && /usr/local/bin/kubectl lable node %s %s"

	k8sJoinTokenCmd = `/usr/local/bin/kubeadm token create --print-join-command |awk '{print $3" "$5" "$7}'`
	k8sJoinMasterCmd = "export CLUSTERNAME=" + ClusterName + " && cd /tmp/pks-installer/k8s/ && ./install.sh master %s " + logJoinCmd
	k8sJoinWorkerCmd = "export CLUSTERNAME=" + ClusterName + " && cd /tmp/pks-installer/k8s/ && ./install.sh worker %s %s" + logJoinCmd
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

