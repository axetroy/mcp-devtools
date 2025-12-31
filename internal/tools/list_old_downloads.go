package tools

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type listOldDownloadsOutput struct {
	System string    `json:"system" jsonschema:"Operating system of the server"`
	Files  []oldFile `json:"files" jsonschema:"List of file paths to check for old downloads"`
}

type oldFile struct {
	Name           string    `json:"name" jsonschema:"Name of the old file"`
	LastModifyTime time.Time `json:"last_modify" jsonschema:"Last modify time of the file"`
	Size           int64     `json:"size" jsonschema:"Size of the file in bytes"`
}

// ListOldDownloads lists files in the Download directory that haven't been accessed in a long time.
func ListOldDownloads(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, *listOldDownloadsOutput, error) {
	downloadDir, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	downloadDir = downloadDir + string(os.PathSeparator) + "Downloads"

	files, err := os.ReadDir(downloadDir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read Downloads directory: %w", err)
	}

	var oldFiles []oldFile

	cutoff := time.Now().AddDate(0, -3, 0) // 3 months ago

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			oldFiles = append(oldFiles, oldFile{
				Name:           info.Name(),
				LastModifyTime: info.ModTime(),
				Size:           info.Size(),
			})
		}
	}

	return nil, &listOldDownloadsOutput{
		System: runtime.GOOS,
		Files:  oldFiles,
	}, nil
}
