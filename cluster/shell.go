package cluster

var (
	sshAuthCmd = "curl -k %s/sshd.tar.gz | tar zxv -C /tmp/ && chmod +x /tmp/sshd/exec.sh && sh /tmp/sshd/exec.sh"
)
