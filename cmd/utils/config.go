package utils

import (
	"encoding/json"
	"fmt"
	"github.com/ligao-cloud-native/xwc-controller-agent/pkg/types"
	"io/ioutil"
	"os"
)

var (
	Env types.Env
)

func init() {
	Env = GetEnv()
}


func GetNodes() types.Nodes {
	nodes, err := GetConfig(nodeConfFile, types.Nodes{})
	nodeObject, ok := nodes.(types.Nodes)
	if !ok {
		panic(err)
	}
	return nodeObject
}

func GetEnv() types.Env {
	env, err := GetConfig(envConfFile, types.Env{})
	envObject, ok := env.(types.Env)
	if !ok {
		panic(err)
	}
	return envObject
}

func GetConfig(path string, result interface{}) (interface{}, error) {

	if ok, _ := PathExists(path); !ok {
		return nil, fmt.Errorf("No such file: %s ",path)
	}
	r, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Read file %s error: %v ", path, err)
	}

	if err := json.Unmarshal(r, &result); err != nil {
		return nil, fmt.Errorf("Unmarshal file %s error: %v ", path, err)
	}

	return result, nil

}


func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}