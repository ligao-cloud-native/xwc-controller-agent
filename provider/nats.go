package provider

import "github.com/ligao-cloud-native/xwc-controller-agent/pkg/agent/nats"

type NatsAgent struct {
	client *nats.Client
}

