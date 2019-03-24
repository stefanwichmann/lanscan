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

import "time"
import "net"
import "fmt"

const defaultTimeout = 50 * time.Millisecond

// ProbeHosts will read hosts to probe from the given channel and check the given port and protocol for each of them.
// Responding hosts will be written back into the second channel.
func ProbeHosts(hosts <-chan string, port int, protocol string, respondingHosts chan<- string, done chan<- bool) {
	adjustedTimeout := defaultTimeout
	for host := range hosts {
		start := time.Now()
		con, err := net.DialTimeout(protocol, fmt.Sprintf("%s:%d", host, port), adjustedTimeout)
		duration := time.Now().Sub(start)
		if err == nil {
			// Host did respond
			con.Close()
			respondingHosts <- host
		}
		// Adjust timeout to current network speed
		if duration < adjustedTimeout {
			difference := adjustedTimeout - duration
			adjustedTimeout = adjustedTimeout - (difference / 2)
		}
	}
	done <- true
}
