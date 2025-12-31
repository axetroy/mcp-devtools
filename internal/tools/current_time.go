package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type currentTimeOutput struct {
	Time string `json:"time" jsonschema:"Current server time in RFC1123 format"`
}

func GetCurrentTime(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, *currentTimeOutput, error) {
	currentTime := fmt.Sprintf("Current server time is: %s", time.Now().Format(time.RFC1123))
	return nil, &currentTimeOutput{Time: currentTime}, nil
}
