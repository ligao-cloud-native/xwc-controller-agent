package utils

import (
	"github.com/ligao-cloud-native/xwc-controller-agent/provider"
	"os"
	"time"
)

func GetAgent() provider.Agent{
	var agent provider.Agent

	switch os.Getenv("PROVIDER") {
	case NatsAgent:
		opt := provider.ExecCmdOpt{Retry: 3, RetrySpan: 2 * time.Second}
		agent = provider.NewNatsAgent(Env.NatsServerUrl, Env.NatsServerToken, opt)
	default:
		opt := provider.ExecCmdOpt{Retry: 3, RetrySpan: 2 * time.Second}
		agent = provider.NewNatsAgent(Env.NatsServerUrl, Env.NatsServerToken, opt)
	}

	return agent
}
