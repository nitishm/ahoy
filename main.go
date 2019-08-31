package main

import (
	"flag"
	"log"

	"github.com/nitishm/ahoy/pkg/istio"
	"istio.io/istio/istioctl/pkg/writer/envoy/configdump"
)

var (
	name      string
	namespace string
)

func init() {
	flag.StringVar(&name, "pod", "", "podname")
	flag.StringVar(&namespace, "ns", "default", "namespace")
}
func main() {
	flag.Parse()

	if name == "" {
		log.Fatal("Pod name cannot be an empty string")
	}

	cd, err := istio.New(name, namespace)
	if err != nil {
		log.Fatal(err)
	}

	listeners, err := cd.Listeners(configdump.ListenerFilter{})
	if err != nil {
		log.Fatal(err)
	}

	for _, listener := range listeners {
		routes, err := cd.FetchRoutes(listener)
		if err != nil {
			log.Fatal(err)
		}

		for _, route := range routes {
			_, err = cd.FetchClusters(route)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}
