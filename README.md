[![Build Status](https://github.com/axetroy/mcp-server-devtools/workflows/ci/badge.svg)](https://github.com/axetroy/mcp-server-devtools/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/mcp-server-devtools)](https://goreportcard.com/report/github.com/axetroy/mcp-server-devtools)
![Latest Version](https://img.shields.io/github/v/release/axetroy/mcp-server-devtools.svg)
![License](https://img.shields.io/github/license/axetroy/mcp-server-devtools.svg)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/mcp-server-devtools.svg)

## MCP DevTools

> A Model Context Protocol (MCP) server that provides useful developer tools for local development.

MCP DevTools is an MCP server implementation built with the [official MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) that exposes useful development tools through the Model Context Protocol. It allows AI assistants and other MCP clients to interact with your local development environment through a standardized interface.

### Features

This MCP server provides the following tools:

#### Color Conversion
- **color_convert** - Convert CSS color values to various color formats
  - Input: CSS color value (e.g., `#ff5733`, `rgb(255, 87, 51)`, `hsl(9, 100%, 60%)`, or named colors like `red`)
  - Output: Hex, RGB, HSL, HSV, CMYK, LAB, XYZ, Linear RGB representations
  - Additional info: Luminance, whether the color is light or dark

#### Network Information
- **get_ip_address** - Get the current computer's IP addresses
  - Returns all active network interface IP addresses
  - Identifies the primary IP address (first non-loopback IPv4)

#### NPM Package Analysis
- **npm_dependencies_analyze** - Get npm package information and analyze its dependencies
  - Input: Package name (e.g., `express`, `react`, `@types/node`) and optional version
  - Output: Package metadata, dependencies, dev dependencies, peer dependencies, version information
  - Fetches data from the official npm registry
  - Analyzes latest version by default, or specify a version

### Usage

The MCP DevTools server communicates via the Model Context Protocol over stdin/stdout. It follows the MCP specification and is built with the official Go SDK.

To use with an MCP client:

```bash
mcp-server-devtools
```

The server will start and wait for MCP requests on stdin, sending responses to stdout.

#### Example: Color Conversion

Request:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "color_convert",
    "arguments": {
      "color": "#ff5733"
    }
  }
}
```

Response:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "hex": "#ff5733",
    "rgb": "rgb(255, 87, 51)",
    "hsl": "hsl(9.0, 100.0%, 60.0%)",
    "hsv": "hsv(9.0, 80.0%, 100.0%)",
    "cmyk": "cmyk(0.0%, 65.9%, 80.0%, 0.0%)",
    "lab": "lab(61.57, 56.45, 51.48)",
    "xyz": "xyz(0.469, 0.305, 0.074)",
    "linear_rgb": "linear-rgb(1.000, 0.106, 0.030)",
    "luminance": 0.428,
    "is_light": false,
    "is_dark": true,
    "original": "#ff5733"
  }
}
```

#### Example: Get IP Address

Request:
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "get_ip_address",
    "arguments": {}
  }
}
```

Response:
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "addresses": ["192.168.1.100", "fe80::1"],
    "primary": "192.168.1.100"
  }
}
```

#### Example: NPM Package Analysis

Request:
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "npm_dependencies_analyze",
    "arguments": {
      "package_name": "express"
    }
  }
}
```

Response:
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "result": {
    "name": "express",
    "version": "5.2.1",
    "description": "Fast, unopinionated, minimalist web framework",
    "license": "MIT",
    "homepage": "http://expressjs.com/",
    "repository": "https://github.com/expressjs/express",
    "dependencies": {
      "accepts": "~1.3.8",
      "body-parser": "2.1.0",
      "content-disposition": "0.5.5",
      "cookie": "0.7.2",
      "cookie-signature": "1.2.3",
      "debug": "2.6.9",
      "escape-html": "~1.0.3",
      "etag": "~1.8.1",
      "finalhandler": "1.3.1",
      "methods": "~1.1.2",
      "mime-types": "~2.1.18",
      "on-finished": "2.4.1",
      "parseurl": "~1.3.3",
      "path-to-regexp": "0.1.12",
      "proxy-addr": "~2.0.7",
      "qs": "6.13.2",
      "range-parser": "~1.2.1",
      "safe-buffer": "5.2.1",
      "send": "1.1.1",
      "serve-static": "2.1.2",
      "setprototypeof": "1.2.0",
      "statuses": "2.0.1",
      "type-is": "~1.6.18",
      "utils-merge": "1.0.1",
      "vary": "~1.1.2"
    },
    "dev_dependencies": {},
    "peer_dependencies": {},
    "dependency_count": 28,
    "author": "TJ Holowaychuk <tj@vision-media.ca>",
    "keywords": ["express", "framework", "sinatra", "web", "http", "rest", "restful", "router", "app", "api"],
    "latest_version": "5.2.1",
    "publish_time": "2024-12-25T14:49:15.000Z"
  }
}
```

You can also analyze a specific version:
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "npm_dependencies_analyze",
    "arguments": {
      "package_name": "express",
      "version": "4.18.0"
    }
  }
}
```

### Install

1. Shell (Mac/Linux)

   ```bash
   curl -fsSL https://github.com/release-lab/install/raw/v1/install.sh | bash -s -- -r=axetroy/mcp-server-devtools -e=mcp-server-devtools
   ```

2. PowerShell (Windows):

   ```powershell
   $r="axetroy/mcp-server-devtools";$e="mcp-server-devtools";iwr https://github.com/release-lab/install/raw/v1/install.ps1 -useb | iex
   ```

3. [Github release page](https://github.com/axetroy/mcp-server-devtools/releases) (All platforms)

   download the executable file and put the executable file to `$PATH`

4. Build and install from source using [Golang](https://golang.org) (All platforms)

   ```bash
   go install github.com/axetroy/mcp-server-devtools/cmd/mcp-server-devtools@latest
   ```

### Development

Build the project:

```bash
make build
```

Run tests:

```bash
make test
```

Format code:

```bash
make format
```

### License

The [MIT License](LICENSE)
