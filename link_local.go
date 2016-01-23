package main

import (
	"log"
	"net"
	"syscall"

	"github.com/vishvananda/netlink"
)

var metadataIP = net.IPv4(169, 254, 169, 254)

func setupLinkLocalNetworking() {
	if _, err := netlink.RouteGet(metadataIP); err != syscall.ENETUNREACH {
		return
	}

	links, err := netlink.LinkList()
	if err != nil {
		log.Fatalf("couldn't list network links: %s", err)
	}

	var device *netlink.Device
	for _, link := range links {
		if dev, ok := link.(*netlink.Device); ok && dev.Flags&net.FlagLoopback == 0 {
			device = dev
			break
		}
	}

	if device == nil {
		log.Fatalf("couldn't find a network interface")
	}

	if err := netlink.AddrAdd(device, &netlink.Addr{IPNet: &net.IPNet{
		IP:   net.IPv4(169, 254, 0, 1),
		Mask: net.IPv4Mask(255, 255, 0, 0),
	}}); err != nil {
		log.Fatalf("couldn't set up link-local networking: %s", err)
	}
}
