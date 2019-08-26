package istio

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"istio.io/istio/istioctl/pkg/kubernetes"
	"istio.io/istio/istioctl/pkg/util/configdump"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/config/host"
)

type ProxyConfig struct {
	configdump *configdump.Wrapper
}

func New(podName, podNamespace string) (*ProxyConfig, error) {
	kubeClient, err := kubernetes.NewClient("~/.kube/config", "")
	if err != nil {
		return nil, err
	}

	path := "config_dump"
	b, err := kubeClient.EnvoyDo(podName, podNamespace, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute command on envoy: %v", err)
	}

	cd := &configdump.Wrapper{}
	err = json.Unmarshal(b, &cd)
	if err != nil {
		return nil, err
	}

	return &ProxyConfig{cd}, nil
}

func (c *ProxyConfig) Listeners() ([]*xdsapi.Listener, error) {
	listenerDump, err := c.configdump.GetListenerConfigDump()
	if err != nil {
		return nil, err
	}
	listeners := make([]*xdsapi.Listener, 0)
	for _, listener := range listenerDump.DynamicActiveListeners {
		if listener.Listener != nil {
			listeners = append(listeners, listener.Listener)
		}
	}

	for _, listener := range listenerDump.StaticListeners {
		if listener.Listener != nil {
			listeners = append(listeners, listener.Listener)
		}
	}
	if len(listeners) == 0 {
		return nil, fmt.Errorf("no listeners found")
	}
	return listeners, nil
}

func (c *ProxyConfig) Clusters() ([]*xdsapi.Cluster, error) {
	safelyParseSubsetKey := func(key string) (model.TrafficDirection, string, host.Name, int) {
		if len(strings.Split(key, "|")) > 3 {
			return model.ParseSubsetKey(key)
		}
		name := host.Name(key)
		return "", "", name, 0
	}

	clusterDump, err := c.configdump.GetClusterConfigDump()
	if err != nil {
		return nil, err
	}
	clusters := make([]*xdsapi.Cluster, 0)
	for _, cluster := range clusterDump.DynamicActiveClusters {
		if cluster.Cluster != nil {
			clusters = append(clusters, cluster.Cluster)
		}
	}
	for _, cluster := range clusterDump.StaticClusters {
		if cluster.Cluster != nil {
			clusters = append(clusters, cluster.Cluster)
		}
	}
	if len(clusters) == 0 {
		return nil, fmt.Errorf("no clusters found")
	}
	sort.Slice(clusters, func(i, j int) bool {
		iDirection, iSubset, iName, iPort := safelyParseSubsetKey(clusters[i].Name)
		jDirection, jSubset, jName, jPort := safelyParseSubsetKey(clusters[j].Name)
		if iName == jName {
			if iSubset == jSubset {
				if iPort == jPort {
					return iDirection < jDirection
				}
				return iPort < jPort
			}
			return iSubset < jSubset
		}
		return iName < jName
	})
	return clusters, nil
}

func (c *ProxyConfig) Routes() ([]*xdsapi.RouteConfiguration, error) {
	routeDump, err := c.configdump.GetRouteConfigDump()
	if err != nil {
		return nil, err
	}
	routes := make([]*xdsapi.RouteConfiguration, 0)
	for _, route := range routeDump.DynamicRouteConfigs {
		if route.RouteConfig != nil {
			routes = append(routes, route.RouteConfig)
		}
	}
	for _, route := range routeDump.StaticRouteConfigs {
		if route.RouteConfig != nil {
			routes = append(routes, route.RouteConfig)
		}
	}
	if len(routes) == 0 {
		return nil, fmt.Errorf("no routes found")
	}
	sort.Slice(routes, func(i, j int) bool {
		iName, err := strconv.Atoi(routes[i].Name)
		if err != nil {
			return false
		}
		jName, err := strconv.Atoi(routes[j].Name)
		if err != nil {
			return false
		}
		return iName < jName
	})
	return routes, nil
}
