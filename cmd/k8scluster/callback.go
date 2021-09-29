package k8scluster

import (
	"bytes"
	"github.com/ligao-cloud-native/xwc-controller-agent/cmd/utils"
	"github.com/ligao-cloud-native/xwc-controller-agent/pkg/types"
	"github.com/ligao-cloud-native/xwc-controller-agent/provider"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"strings"
)

// Callback用于修改pwc-controller服务修改集群状态
type Callback struct {
	Step int
	Second int
	// CallbackUrl为pwc-controller服务暴漏的rest server
	// http://pwc-controller.pks-system:7000
	CallbackUrl string
	Command string
	Agent provider.Agent
	Nodes types.Nodes
}

func NewCallback(agent provider.Agent, nodes types.Nodes) *Callback {
	return &Callback{
		Step: 5,
		Second: 5,
		CallbackUrl: utils.Env.CallbackUrl,
		Agent: agent,
		Nodes: nodes,
	}
}

func (c *Callback) Execute() {
	if len(c.Nodes.Masters) > 0 {
		initMasterNode := c.Nodes.Masters[0]
		result, success := c.Agent.CommandExecute(initMasterNode, c.Command)
		if success {
			var res []byte

			if strings.Contains(result, "errorInfo") {
				index := strings.Index(result, "errorInfo")
				res = []byte(result[:index])
			} else {
				res = []byte(result)
			}
			c.sendMessage(res)
		}

	}

}

// sendMessage调用pwc-controller服务暴漏的rest server，修改集群状态
func (c *Callback) sendMessage(data []byte) {
	body := bytes.NewReader(data)
	url := c.CallbackUrl + "/task-completion-callback"
	req ,err := http.NewRequest("POST", url, body)
	if err != nil {
		klog.Error(err)
		return
	}
	req.Header.Set("Content-Type","application/json;charset=UTF-8")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		klog.Error(err)
		return
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		klog.Error(err)
		return
	}

	klog.Info(string(resBody))



}