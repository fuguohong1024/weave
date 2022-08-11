//go:build iface && !mcast
// +build iface,!mcast

package main

import (
	weavenet "github.com/fuguohong1024/weave/net"
)

func checkNetwork() error {
	_, err := weavenet.EnsureInterface(weavenet.VethName)
	return err
}
