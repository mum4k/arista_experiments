/*
fetch_interfaces retrieves information about interfaces on an Arista box.

If the ifName interface doesn't have its primary IP set to ipAddress/mask, it will be configured.
*/
package main

import (
	"fmt"

	"github.com/aristanetworks/goeapi"
	"github.com/aristanetworks/goeapi/module"
	"github.com/kylelemons/godebug/pretty"
)

const (
	ifName    = "Management1"
	ipAddress = "10.1.2.90/24"
)

// connect connects to the arista box.
func connect() (*goeapi.Node, error) {
	goeapi.LoadConfig("eapi.conf")

	node, err := goeapi.ConnectTo("arista")
	if err != nil {
		return nil, err
	}
	return node, nil
}

// showInterfaces executes the 'show interfaces' command.
func showInterfaces(node *goeapi.Node) (*module.ShowInterface, error) {
	h, err := node.GetHandle("json")
	if err != nil {
		return nil, err
	}

	si := &module.ShowInterface{}
	if err := h.AddCommand(si); err != nil {
		return nil, err
	}

	if err := h.Call(); err != nil {
		return nil, err
	}
	return si, nil
}

// hasIP determines if the provided interface has the specified primary IP address.
func hasPrimaryIP(name string, node *goeapi.Node) (bool, error) {
	iface := module.IPInterface(node)
	cfg, err := iface.Get(name)
	if err != nil {
		return false, err
	}
	return cfg.Address() == ipAddress, nil
}

// setPrimaryIP sets the primary IP address on an interface.
func setPrimaryIP(name string, node *goeapi.Node, ip string) error {
	iface := module.IPInterface(node)
	if !iface.SetAddress(name, ip) {
		return fmt.Errorf("failed to configure IP %s on interface %s.", ip, name)
	}
	return nil
}

func main() {
	node, err := connect()
	if err != nil {
		panic(err)
	}

	si, err := showInterfaces(node)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ShowInterfaces:\n%s\n\n", pretty.Sprint(si))

	hasIP, err := hasPrimaryIP(ifName, node)
	if err != nil {
		panic(err)
	}
	if hasIP {
		fmt.Printf("Interface %s already has IP %s.\n", ifName, ipAddress)
		return
	}

	fmt.Printf("Adding IP %s to interface %s.\n", ipAddress, ifName)
	if err := setPrimaryIP(ifName, node, ipAddress); err != nil {
		panic(err)
	}
}
