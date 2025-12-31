package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Luminance coefficients for relative luminance calculation (ITU-R BT.709)
const (
	RedLuminance   = 0.2126
	GreenLuminance = 0.7152
	BlueLuminance  = 0.0722
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

// rgbToCMYK converts RGB values (0-1 range) to CMYK values (0-1 range)
func rgbToCMYK(r, g, b float64) (c, m, y, k float64) {
	k = 1.0 - max(r, g, b)
	if k < 1.0 {
		c = (1.0 - r - k) / (1.0 - k)
		m = (1.0 - g - k) / (1.0 - k)
		y = (1.0 - b - k) / (1.0 - k)
	}
	return c, m, y, k
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

// colorInput represents the input for color conversion tool
type colorInput struct {
	Color string `json:"color" jsonschema:"CSS color value (e.g., '#ff5733', 'rgb(255, 87, 51)', 'hsl(9, 100%, 60%)', 'red')"`
}

// colorOutput represents the output of color conversion
type colorOutput struct {
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

// ColorConversion converts CSS color values to various color formats
func ColorConversion(ctx context.Context, req *mcp.CallToolRequest, input colorInput) (*mcp.CallToolResult, *colorOutput, error) {
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

	// Calculate CMYK
	rf, gf, bf := float64(r)/255.0, float64(g)/255.0, float64(b)/255.0
	c, m, y, k := rgbToCMYK(rf, gf, bf)

	lab_l, lab_a, lab_b := color.Lab()
	x, yv, z := color.Xyz()
	lr, lg, lb := color.LinearRgb()
	luminance := (RedLuminance*float64(r) + GreenLuminance*float64(g) + BlueLuminance*float64(b)) / 255.0

	output := &colorOutput{
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
