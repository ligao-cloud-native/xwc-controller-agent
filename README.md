# xwc-controller-agent
controller通过创建job执行任务，
controller-agent则执行job中的具体任务，通过下发命令，到边缘节点去创建/删除/扩缩容集群。

## 支持安装的集群类型
serverless集群：创建一个kubeless集群，能够实现fass和弹性扩缩容

istio集群：安装了istio的k8s集群，可用于流量管理

## nats
由于controller和controller-agent都运行在云管区，需要解决如何连接到边缘节点，进行集群操作？
通过消息队列nats实现。

默认provider为nats，使用agent为nats client，连接的是托管集群的vmserver服务，提供接口：

(1) 获取token: GET /token

(2) 获取所有集群节点的nodeid：GET /pks/api/v1/workers

(3) 执行命令：POST /pks/api/v1/execute

(4) 获取命令执行结果：/pks/api/v1/execute/<taskID>

原理分析: vmserver通过rest api接收指定，然后发送到nats指定的队列，在边缘k8s集群上的每个节点上，
部署vmagent，用于接收队列中的指令，执行完成后发送出去，然后vmserver订阅到消息，通过api提供。
 

## 针对集群安装， controller-agent的工作流程
主要执行kubemec-setup项目中的安装脚本

### ssh认证
根据节点的角色进行验证，如etcd节点，master节点，worker节点
主要通过下发指令：curl -k <fileserver>/sshd.tar.gz | tar zxv -C /tmp/ && chmod +x /tmp/sshd/exec.sh && sh /tmp/sshd/exec.sh

### download etcd package
下发的指令：
curl -k <fileserver>/etcd-installer/clean-etcd.sh | sh -s -- -y 2>&1 |tee /opt/reset.log
curl -k <fileserver>/install-etcd | sh 2>&1 |tee /opt/install.log

### install etcd cluster
curl -k <fileserver>/install-etcd | sh -s init 2>&1 |tee -a /opt/install.log 

### download k8s package and install k8s
curl -k <fileserver>/install | sh -s init  相当于执行sh install init 

### join k8s cluster
 