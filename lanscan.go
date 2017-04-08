package lanscan

import (
	"errors"
	"fmt"
	"time"
	"net"
)

const maxAddressesPerSubnet = 1000

func ScanLinkLocal(network string, port int, threads int, timeout time.Duration) ([]string, error) {
	if !validateNetwork(network) {
		return []string{}, errors.New(fmt.Sprintf("Invalid network %s (Valid options: %v)", network, validNetworks))
	}

	if port < 0 || port > 65535 {
		return []string{}, errors.New(fmt.Sprintf("Invalid port %d (Valid options: 0 - 65535)", port))
	}

	hosts := make(chan string, 1000)
	results := make(chan string, 1000)
	done := make(chan bool, threads)

	for worker := 0; worker < threads; worker++ {
		go ProbeHosts(hosts, port, network, results, done)
	}


	for _, current := range LinkLocalAddresses(network) {
		allIPs := CalculateSubnetIPs(current, maxAddressesPerSubnet)
		ip, _, err := net.ParseCIDR(current)
		if err != nil {
			continue
		}
		startIndex := 0
		for index, value := range allIPs {
			if value == ip.String() {
				startIndex = index
				break
			}
		}

		for i := startIndex + 1; i < len(allIPs); i++ {
			hosts <- allIPs[i]
			if (startIndex - i) >= 0 {
				hosts <- allIPs[startIndex-i]
			}
		}
	}
	close(hosts)

	var responses = []string{}
	for {
		select {
		case found := <-results:
			responses = append(responses, found)
		case <-done:
			threads--
			if threads == 0 {
				return responses, nil
			}
		case <-time.After(timeout):
			return responses, nil
		}
	}
}

func validateNetwork(network string) bool {
	for _, net := range validNetworks {
		if network == net {
			return true
		}
	}
	return false
}
