// fetch_interfaces retrieves information about interfaces on an Arista box.
package main

import (
	"fmt"

	"github.com/aristanetworks/goeapi"
	"github.com/kylelemons/godebug/pretty"
)

type ShowInterfaces struct {
	Interfaces map[string]Interface
}

func (si *ShowInterfaces) GetCmd() string {
	return "show interfaces"
}

type Interface struct {
	Name string
}

func main() {
	goeapi.LoadConfig("eapi.conf")

	node, err := goeapi.ConnectTo("arista")
	if err != nil {
		panic(err)
	}

	si := &ShowInterfaces{}
	h, err := node.GetHandle("json")
	if err != nil {
		panic(err)
	}

	if err := h.AddCommand(si); err != nil {
		panic(err)
	}

	if err := h.Call(); err != nil {
		panic(err)
	}

	fmt.Printf("Output:\n%s\n", pretty.Sprint(si))
}
