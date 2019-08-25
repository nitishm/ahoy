package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"istio.io/istio/istioctl/pkg/kubernetes"
	"istio.io/istio/istioctl/pkg/writer/envoy/configdump"
)

func main() {
	podName := "istio-ingressgateway-1-0-0-40-dbg-9d9c95d8-jjnw4"
	podNamespace := "fed-test-host"
	configWriter, err := setupConfigdumpEnvoyConfigWriter(podName, podNamespace, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	configWriter.PrintListenerSummary(configdump.ListenerFilter{})
}

func setupConfigdumpEnvoyConfigWriter(podName, podNamespace string, out io.Writer) (*configdump.ConfigWriter, error) {
	kubeClient, err := kubernetes.NewClient("~/.kube/config", "")
	if err != nil {
		return nil, err
	}

	path := "config_dump"
	debug, err := kubeClient.EnvoyDo(podName, podNamespace, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute command on envoy: %v", err)
	}

	cw := &configdump.ConfigWriter{Stdout: out}
	err = cw.Prime(debug)
	if err != nil {
		return nil, err
	}

	return cw, nil
}
