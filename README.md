# **Lanscan** - Blazing fast, local network scanning in Go
[![Build Status](https://travis-ci.org/stefanwichmann/lanscan.svg?branch=master)](https://travis-ci.org/stefanwichmann/lanscan)
[![Go Report Card](https://goreportcard.com/badge/github.com/stefanwichmann/lanscan)](https://goreportcard.com/report/github.com/stefanwichmann/lanscan)
[![license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/stefanwichmann/lanscan/blob/master/LICENSE)

# Introduction
Lanscan is a small, blazing fast and easy to use golang library to scan for hosts in your local network. It's job is to identify hosts nearby that are listening on a specified port. The goal is to make scans like these as fast and simple as possible. Just provide a port, the number of parallel threads and a timeout. Lanscan will take care of the rest...

# Features
- [x] Automatic restriction to your link local network
- [x] Automatic prioritisation of nearby IPs
- [x] Automatic scan on all network interfaces
- [x] Automatic adaption to your network latency
- [x] Full support for parallel scanning
- [x] Ability to stop after a timeout

# Getting started
```go
package main

import "github.com/stefanwichmann/lanscan"
import "time"
import "log"

func main() {
  // Scan for hosts listening on tcp port 80.
  // Use 20 threads and timeout after 5 seconds.
  hosts, err := lanscan.ScanLinkLocal("tcp4", 80, 20, 5*time.Second)
  if err != nil {
    log.Fatal(err)
  }
  for _, host := range hosts {
    log.Printf("Host %v responded.", host)
  }
}

```

# Status
Lanscan is still work in progress and far from done! Right now it's working stable in an IPv4 environment when scanning for TCP services. Once a TCP handshake is successfully completed Lanscan considers the host reachable. Open tasks right now are:

- [ ] IPv6 discovery
- [ ] UDP discovery
- [ ] Provide a proper command line tool for scanning
