// MIT License
//
// Copyright (c) 2017 Stefan Wichmann
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
package lanscan

import (
	"errors"
	"fmt"
	"net"
	"time"
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
