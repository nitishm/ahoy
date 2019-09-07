package main

import (
	"flag"
	"fmt"
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
		routeConfigurations, err := cd.FetchRouteConfigurations(listener)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\n\n== Listener ==\n%v\n", listener.GetName())
		for _, routeConfiguration := range routeConfigurations {
			virtualHosts, err := cd.FetchVirtualHosts(routeConfiguration)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("== RouteConfiguration ==\n%v\n", routeConfiguration.GetName())
			for _, virtualHost := range virtualHosts {
				routes, err := cd.FetchRoutes(virtualHost)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("== VirtualHost ==\n%v\n", routeConfiguration.GetName())
				for _, route := range routes {
					cluster, err := cd.FetchCluster(route)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Printf("== Route ==\n%v\n== Cluster== \n%v\n", route.GetMatch().GetPrefix(), cluster.GetName())
				}
			}
		}
	}

}
