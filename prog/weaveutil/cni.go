package main

import (
	"os"

	cni "github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/version"

	weaveapi "github.com/fuguohong1024/weave/api"
	"github.com/fuguohong1024/weave/common"
	ipamplugin "github.com/fuguohong1024/weave/plugin/ipam"
	netplugin "github.com/fuguohong1024/weave/plugin/net"
)

func cniNet(args []string) error {
	weave := weaveapi.NewClient(os.Getenv("WEAVE_HTTP_ADDR"), common.Log)
	n := netplugin.NewCNIPlugin(weave)
	cni.PluginMain(n.CmdAdd, n.CmdDel, version.PluginSupports("0.1.0", "0.2.0", "0.3.0"))
	return nil
}

func cniIPAM(args []string) error {
	weave := weaveapi.NewClient(os.Getenv("WEAVE_HTTP_ADDR"), common.Log)
	i := ipamplugin.NewIpam(weave)
	cni.PluginMain(i.CmdAdd, i.CmdDel, version.PluginSupports("0.1.0", "0.2.0", "0.3.0"))
	return nil
}
