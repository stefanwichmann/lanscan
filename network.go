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

import "net"
import "strings"

var linkLocalCIDR = []string{
	"10.0.0.0/8",     // RFC ???
	"172.16.0.0/12",  // RFC ???
	"192.168.0.0/16", // RFC ???
	"169.254.0.0/16", // RFC ???
	"fc00::/7",       // RFC ???
	"fe80::/64"}      // RFC ???

var validNetworks = []string{"tcp", "tcp4", "tcp6", "udp", "udp4", "udp6", "ip", "ip4", "ip6", "unix", "unixgram", "unixpacket"}

func IsLinkLocalAddress(ip net.IP) bool {
	for _, curnet := range linkLocalCIDR {
		_, ipnet, err := net.ParseCIDR(curnet)
		if err != nil {
			return false
		}

		if !ip.IsLoopback() && ipnet.Contains(ip) {
			return true
		}
	}

	return false
}

func LinkLocalAddresses(network string) []string {
	var addrs = []string{}
	interfaces, err := net.Interfaces()
	if err != nil {
		return addrs
	}

	for _, i := range interfaces {
		addresses, _ := i.Addrs()
		for _, address := range addresses {
			ip, _, err := net.ParseCIDR(address.String())
			if err != nil {
				continue
			}

			// should scan only v4 Hosts and current IP is v6
			if strings.Contains(network, "4") && ip.To4() == nil {
				continue
			}

			// should scan only v6 Hosts and current IP is v4
			if strings.Contains(network, "6") && ip.To16() == nil {
				continue
			}

			if IsLinkLocalAddress(ip) {
				addrs = append(addrs, address.String())
			}
		}
	}

	return addrs
}

func CalculateSubnetIPs(cidr string, maxAddresses int) []string {
	var ips = []string{}

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return ips
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
		if len(ips) > maxAddresses {
			return ips[1:len(ips)] // remove network address
		}
	}

	// remove network address and broadcast address
	if len(ips) > 1 {
		return ips[1 : len(ips)-1]
	}
	return ips
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
