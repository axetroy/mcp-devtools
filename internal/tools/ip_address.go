package tools

import (
	"context"
	"fmt"
	"net"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ipAddressOutput struct {
	Addresses []string `json:"addresses" jsonschema:"List of IP addresses"`
	Primary   string   `json:"primary" jsonschema:"Primary IP address (first non-loopback IPv4)"`
}

// GetIPAddress returns the current computer's IP addresses
func GetIPAddress(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, *ipAddressOutput, error) {
	addresses := []string{}
	primary := ""

	// Get all network interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	for _, iface := range ifaces {
		// Skip loopback and down interfaces
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ipStr := ip.String()
			addresses = append(addresses, ipStr)

			// Set primary as the first non-loopback IPv4 address
			if primary == "" && ip.To4() != nil {
				primary = ipStr
			}
		}
	}

	if len(addresses) == 0 {
		return nil, nil, fmt.Errorf("no IP addresses found")
	}

	if primary == "" && len(addresses) > 0 {
		primary = addresses[0]
	}

	return nil, &ipAddressOutput{
		Addresses: addresses,
		Primary:   primary,
	}, nil
}
