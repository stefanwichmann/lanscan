package lanscan

import "time"
import "net"
import "fmt"

const defaultTimeout = 50 * time.Millisecond

func ProbeHosts(hosts <-chan string, port int, protocol string, respondingHosts chan<- string, done chan<- bool) {
	adjustedTimeout := defaultTimeout
	for host := range hosts {
		adjustedTimeout := adjustedTimeout
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
