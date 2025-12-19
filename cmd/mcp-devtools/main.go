package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

// max returns the maximum of three float64 values
func max(a, b, c float64) float64 {
	if a > b {
		if a > c {
			return a
		}
		return c
	}
	if b > c {
		return b
	}
	return c
}

// ColorInput represents the input for color conversion tool
type ColorInput struct {
	Color string `json:"color" jsonschema:"CSS color value (e.g., '#ff5733', 'rgb(255, 87, 51)', 'hsl(9, 100%, 60%)', 'red')"`
}

// ColorOutput represents the output of color conversion
type ColorOutput struct {
	Hex       string  `json:"hex" jsonschema:"Hexadecimal color representation"`
	RGB       string  `json:"rgb" jsonschema:"RGB color representation"`
	HSL       string  `json:"hsl" jsonschema:"HSL color representation"`
	HSV       string  `json:"hsv" jsonschema:"HSV color representation"`
	CMYK      string  `json:"cmyk" jsonschema:"CMYK color representation"`
	LAB       string  `json:"lab" jsonschema:"LAB color representation"`
	XYZ       string  `json:"xyz" jsonschema:"XYZ color representation"`
	LinearRGB string  `json:"linear_rgb" jsonschema:"Linear RGB color representation"`
	Luminance float64 `json:"luminance" jsonschema:"Relative luminance (0-1)"`
	IsLight   bool    `json:"is_light" jsonschema:"Whether the color is light (luminance > 0.5)"`
	IsDark    bool    `json:"is_dark" jsonschema:"Whether the color is dark (luminance <= 0.5)"`
	Original  string  `json:"original" jsonschema:"Original input color value"`
}

// IPAddressOutput represents the output of IP address tool
type IPAddressOutput struct {
	Addresses []string `json:"addresses" jsonschema:"List of IP addresses"`
	Primary   string   `json:"primary" jsonschema:"Primary IP address (first non-loopback IPv4)"`
}

// ColorConversionTool converts CSS color values to various color formats
func ColorConversionTool(ctx context.Context, req *mcp.CallToolRequest, input ColorInput) (*mcp.CallToolResult, *ColorOutput, error) {
	// Parse the color
	color, err := colorful.Hex(input.Color)
	if err != nil {
		// Try parsing as named color or other CSS format
		color, err = parseColor(input.Color)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse color '%s': %w", input.Color, err)
		}
	}

	// Get various color representations
	r, g, b := color.RGB255()
	h, s, l := color.Hsl()
	hv, sv, v := color.Hsv()

	// Calculate CMYK manually
	rf, gf, bf := float64(r)/255.0, float64(g)/255.0, float64(b)/255.0
	k := 1.0 - max(rf, gf, bf)
	var c, m, y float64
	if k < 1.0 {
		c = (1.0 - rf - k) / (1.0 - k)
		m = (1.0 - gf - k) / (1.0 - k)
		y = (1.0 - bf - k) / (1.0 - k)
	}

	lab_l, lab_a, lab_b := color.Lab()
	x, yv, z := color.Xyz()
	lr, lg, lb := color.LinearRgb()
	luminance := (0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)) / 255.0

	output := &ColorOutput{
		Hex:       color.Hex(),
		RGB:       fmt.Sprintf("rgb(%d, %d, %d)", r, g, b),
		HSL:       fmt.Sprintf("hsl(%.1f, %.1f%%, %.1f%%)", h, s*100, l*100),
		HSV:       fmt.Sprintf("hsv(%.1f, %.1f%%, %.1f%%)", hv, sv*100, v*100),
		CMYK:      fmt.Sprintf("cmyk(%.1f%%, %.1f%%, %.1f%%, %.1f%%)", c*100, m*100, y*100, k*100),
		LAB:       fmt.Sprintf("lab(%.2f, %.2f, %.2f)", lab_l, lab_a, lab_b),
		XYZ:       fmt.Sprintf("xyz(%.3f, %.3f, %.3f)", x, yv, z),
		LinearRGB: fmt.Sprintf("linear-rgb(%.3f, %.3f, %.3f)", lr, lg, lb),
		Luminance: luminance,
		IsLight:   luminance > 0.5,
		IsDark:    luminance <= 0.5,
		Original:  input.Color,
	}

	return nil, output, nil
}

// parseColor attempts to parse various CSS color formats
func parseColor(colorStr string) (colorful.Color, error) {
	colorStr = strings.TrimSpace(colorStr)

	// Try hex format
	if strings.HasPrefix(colorStr, "#") {
		return colorful.Hex(colorStr)
	}

	// Try RGB format
	if strings.HasPrefix(colorStr, "rgb") {
		var r, g, b uint8
		_, err := fmt.Sscanf(colorStr, "rgb(%d,%d,%d)", &r, &g, &b)
		if err != nil {
			_, err = fmt.Sscanf(colorStr, "rgb(%d, %d, %d)", &r, &g, &b)
		}
		if err == nil {
			return colorful.Color{R: float64(r) / 255.0, G: float64(g) / 255.0, B: float64(b) / 255.0}, nil
		}
	}

	// Try HSL format
	if strings.HasPrefix(colorStr, "hsl") {
		var h, s, l float64
		_, err := fmt.Sscanf(colorStr, "hsl(%f,%f%%,%f%%)", &h, &s, &l)
		if err != nil {
			_, err = fmt.Sscanf(colorStr, "hsl(%f, %f%%, %f%%)", &h, &s, &l)
		}
		if err == nil {
			return colorful.Hsl(h, s/100.0, l/100.0), nil
		}
	}

	// Try named colors
	namedColors := map[string]string{
		"red": "#ff0000", "green": "#008000", "blue": "#0000ff",
		"white": "#ffffff", "black": "#000000", "yellow": "#ffff00",
		"cyan": "#00ffff", "magenta": "#ff00ff", "gray": "#808080",
		"orange": "#ffa500", "purple": "#800080", "pink": "#ffc0cb",
		"brown": "#a52a2a", "lime": "#00ff00", "navy": "#000080",
	}

	if hex, ok := namedColors[strings.ToLower(colorStr)]; ok {
		return colorful.Hex(hex)
	}

	return colorful.Color{}, fmt.Errorf("unable to parse color")
}

// GetIPAddressTool returns the current computer's IP addresses
func GetIPAddressTool(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, *IPAddressOutput, error) {
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

	return nil, &IPAddressOutput{
		Addresses: addresses,
		Primary:   primary,
	}, nil
}

func main() {
	// Create MCP server using the official SDK
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "mcp-devtools",
			Version: version,
		},
		&mcp.ServerOptions{
			Instructions: "A collection of useful developer tools including color conversion and network information.",
		},
	)

	// Register color conversion tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "color_convert",
		Description: "Convert CSS color values to various color formats (Hex, RGB, HSL, HSV, CMYK, LAB, XYZ, Linear RGB). Supports hex (#ff5733), rgb(255, 87, 51), hsl(9, 100%, 60%), and named colors (red, blue, etc.)",
	}, ColorConversionTool)

	// Register IP address tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_ip_address",
		Description: "Get the current computer's IP addresses, including all network interfaces and the primary IP address",
	}, GetIPAddressTool)

	// Run the server over stdin/stdout
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
