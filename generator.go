package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"

	"github.com/coreos/go-systemd/unit"
)

func main() {
	outPath := os.Args[1]

	metadata := fetchMetadata()
	log.Printf("received metadata, configuring host %s", metadata.Hostname)

	units := make(map[string]*bytes.Buffer)
	for typ, ifaces := range metadata.Interfaces {
		for index, iface := range ifaces {
			file := fmt.Sprintf("digitalocean-%s-%d.network", typ, index)
			opts := iface.toUnitOptions()
			if len(opts) == 0 {
				continue
			}
			if typ == "public" {
				opts = append(opts, metadata.DNS.toUnitOptions()...)
			}
			buf := unit.Serialize(opts).(*bytes.Buffer)
			units[file] = buf
		}
	}

	for file, buf := range units {
		if err := ioutil.WriteFile(path.Join(outPath, file), buf.Bytes(), 0644); err != nil {
			log.Fatalf("couldn't write unit file: %s", err)
		}
	}
}

func (dns *DNS) toUnitOptions() (opts []*unit.UnitOption) {
	if dns == nil {
		return nil
	}
	for _, nameserver := range dns.Nameservers {
		if ip := net.ParseIP(nameserver); ip != nil {
			opts = append(opts, &unit.UnitOption{"Network", "DNS", ip.String()})
		}
	}
	return opts
}

func (iface *Interface) toUnitOptions() (opts []*unit.UnitOption) {
	if iface == nil {
		return nil
	}
	opts = append(opts, iface.IPv4.toUnitOptions()...)
	opts = append(opts, iface.IPv6.toUnitOptions()...)
	opts = append(opts, iface.AnchorIPv4.toUnitOptions()...)
	if len(opts) == 0 {
		return nil
	}
	return append(opts, &unit.UnitOption{"Match", "MACAddress", iface.MACAddress})
}

func (iface *InterfaceV4) toUnitOptions() []*unit.UnitOption {
	if iface == nil {
		return nil
	}
	network := net.IPNet{
		net.ParseIP(iface.IPAddress),
		net.IPMask(net.ParseIP(iface.Netmask)),
	}
	return []*unit.UnitOption{
		{"Network", "Address", network.String()},
		{"Network", "Gateway", iface.Gateway},
	}
}

func (iface *InterfaceV6) toUnitOptions() []*unit.UnitOption {
	if iface == nil {
		return nil
	}
	network := net.IPNet{
		net.ParseIP(iface.IPAddress),
		net.CIDRMask(iface.CIDR, 8*net.IPv6len),
	}
	return []*unit.UnitOption{
		{"Network", "Address", network.String()},
		{"Network", "Gateway", iface.Gateway},
	}
}
