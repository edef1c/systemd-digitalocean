package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func fetchMetadata() (metadata *Metadata) {
	setupLinkLocalNetworking()

	resp, err := http.Get(fmt.Sprintf("http://%s/metadata/v1.json", metadataIP))
	if err != nil {
		log.Fatalf("couldn't retrieve metadata: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("couldn't retrieve metadata: %s %d %s", resp.Proto, resp.StatusCode, resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		log.Fatalf("couldn't parse metadata: %s", err)
	}
	return
}

type Metadata struct {
	Hostname string `json:"hostname,omitempty"`
	Region   string `json:"region,omitempty"`

	DNS        *DNS                   `json:"dns,omitempty"`
	Interfaces map[string][]Interface `json:"interfaces",omitempty`
}

type DNS struct {
	Nameservers []string `json:"nameservers,omitempty"`
}

type Interface struct {
	MACAddress string       `json:"mac,omitempty"`
	IPv4       *InterfaceV4 `json:"ipv4,omitempty"`
	IPv6       *InterfaceV6 `json:"ipv6,omitempty"`
	AnchorIPv4 *InterfaceV4 `json:"anchor_ipv4",omitempty`
}

type InterfaceV4 struct {
	IPAddress string `json:"ip_address,omitempty"`
	Netmask   string `json:"netmask,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
}

type InterfaceV6 struct {
	IPAddress string `json:"ip_address,omitempty"`
	CIDR      int    `json:"cidr,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
}
