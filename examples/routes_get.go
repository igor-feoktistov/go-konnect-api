package main

import (
	"fmt"
	"time"
	"strings"
	"go-konnect-api/konnect"
)

func main() {
	token := "kpat_<reducted>"
	c := konnect.NewClient(
		"https://us.api.konghq.com",
		&konnect.ClientOptions {
		    Token: token,
		    SSLVerify: false,
		    Debug: false,
    		    Timeout: 60 * time.Second,
		},
	)
	controlPlanes, err := c.ControlPlanesGet([]string{})
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	} else {
		for _, controlPlane := range controlPlanes {
			fmt.Printf("%s\n", controlPlane.Name)
			routes, err := c.RoutesGet(controlPlane.Id, []string{})
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			} else {
				for _, route := range routes {
					if len(route.Hosts) > 0 {
						var routeName, routeNamespace string
						for _, tag := range route.Tags {
							if strings.HasPrefix(tag, "k8s-name:") {
								routeName = tag[9:]
							}
							if strings.HasPrefix(tag, "k8s-namespace:") {
								routeNamespace = tag[14:]
							}
						}
						fmt.Printf("\tname=%s/%s, host=%s\n", routeNamespace, routeName, route.Hosts[0])
					}
				}
			}
		}
	}
}
