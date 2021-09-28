package nats

import (
	"github.com/ligao-cloud-native/xwc-controller-agent/pkg/types"
	"time"
)

type Client struct {
	config *config
}

type config struct {
	url string
	token string
}

func NewConfig(url, token string) *config {
	return &config{url: url, token:token}
}

func (cfg *config) BuildClientOrDie() *Client {
	client := &Client{config: cfg}
	if cfg.token == "" {
		client.geneToken()
	}

	return client
}

func (c *Client) geneToken() {
	// TODO: get token, GET c.config.url/token
	c.config.token = ""
}

func (c *Client) GetWorkers() {
	url := c.config.url + "/pks/api/v1/worker"
	//Todo : do get rest api
}

func (c *Client) ExecCmdAndGetInfo(cmd types.CmdExecInfo, retry int, retrySpan time.Duration) (CmdExecInfo,error){
	return CmdExecInfo{}, nil
}

func (c *Client) ExecCmd() {
	url := c.config.url + "/pks/api/v1/execute"
	//Todo : do post rest api
}

func (c *Client) GetExecCmdInfo() {
	url := c.config.url + "/pks/api/v1/execute/<taskID>"
	//Todo : do get rest api
}

