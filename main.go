package main

import (
	"fmt"
	"log"

	"github.com/nitishm/ahoy/pkg/istio"
)

func main() {
	podName := "istio-ingressgateway-1-0-0-40-dbg-9d9c95d8-jjnw4"
	podNamespace := "fed-test-host"
	cd, err := istio.New(podName, podNamespace)
	if err != nil {
		log.Fatal(err)
	}

	listeners, err := cd.Listeners()
	if err != nil {
		log.Fatal(err)
	}

	for _, listener := range listeners {
		fmt.Println(listener.Name)
	}

}
