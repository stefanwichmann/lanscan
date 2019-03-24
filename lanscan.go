// MIT License
//
// Copyright (c) 2019 Stefan Wichmann
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package lanscan contains a blazing fast port scanner for local networks
package lanscan

import (
	"fmt"
	"net"
	"time"
)

const maxAddressesPerSubnet = 1000

// ScanLinkLocal scans all link local networks on all interfaces found on the current computer for hosts
// responding on the given port. It will use the given amout of threads and will return after the given timeout
// or after finishing the scan.
func ScanLinkLocal(network string, port int, threads int, timeout time.Duration) ([]string, error) {
	// Validate parameters
	if !validateNetwork(network) {
		return []string{}, fmt.Errorf("Invalid network %s (Valid options: %v)", network, validNetworks)
	}
	if port < 0 || port > 65535 {
		return []string{}, fmt.Errorf("Invalid port %d (Valid options: 0 - 65535)", port)
	}

	hosts := make(chan string, 100)
	results := make(chan string, 10)
	done := make(chan bool, threads)

	// Start workers
	for worker := 0; worker < threads; worker++ {
		go ProbeHosts(hosts, port, network, results, done)
	}

	// Generate host list to check
	for _, current := range LinkLocalAddresses(network) {
		allIPs := CalculateSubnetIPs(current, maxAddressesPerSubnet)

		startIndex := findIndex(current, allIPs)
		for i := startIndex + 1; i < len(allIPs); i++ {
			hosts <- allIPs[i] // add all following hosts to channel
			if (startIndex - i) >= 0 {
				hosts <- allIPs[startIndex-i] // add all previous hosts to channel
			}
		}
	}
	close(hosts)

	// collect responses
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

func findIndex(candidate string, hosts []string) int {
	ip, _, err := net.ParseCIDR(candidate)
	if err != nil {
		return 0
	}
	for index, value := range hosts {
		if value == ip.String() {
			return index
		}
	}
	return 0
}
