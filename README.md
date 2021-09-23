# xwc-controller-agent
controller-agent为controller的job具体执行的任务，即到边缘节点去安装集群。

运行在云管区，需要解决如何连接到边缘节点，进行集群安装？通过消息队列nats实现。

## nats

默认provider为nats，使用agent为nats client，连接的是托管集群的vmserver服务，提供接口：

(1) 获取token: GET /token

(2) 获取所有集群节点的nodeid：GET /pks/api/v1/workers

(3) 执行命令：POST /pks/api/v1/execute

(4) 获取命令执行结果：/pks/api/v1/execute/<taskID>

原理分析: vmserver通过rest api接收指定，然后发送到nats指定的队列，在边缘k8s集群上的每个节点上，
部署vmagent，用于接收队列中的指令，执行完成后发送出去，然后vmserver订阅到消息，通过api提供。
 

## 针对集群安装， controller-agent的工作流程

（1）ssh认证

（2）