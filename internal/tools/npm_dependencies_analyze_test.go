package tools

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestNpmDependenciesAnalyze(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	// Test with a well-known package
	input := npmPackageInput{
		PackageName: "express",
	}

	_, output, err := NpmDependenciesAnalyze(ctx, req, input)
	if err != nil {
		t.Fatalf("Failed to analyze express package: %v", err)
	}

	if output == nil {
		t.Fatal("Output should not be nil")
	}

	if output.Name != "express" {
		t.Errorf("Expected package name 'express', got '%s'", output.Name)
	}

	if output.Version == "" {
		t.Error("Version should not be empty")
	}

	if output.LatestVersion == "" {
		t.Error("Latest version should not be empty")
	}

	if output.Description == "" {
		t.Error("Description should not be empty")
	}

	if output.Dependencies == nil {
		t.Error("Dependencies map should not be nil")
	}

	t.Logf("Package: %s", output.Name)
	t.Logf("Version: %s", output.Version)
	t.Logf("Latest: %s", output.LatestVersion)
	t.Logf("Dependencies: %d", output.DependencyCount)
	t.Logf("Description: %s", output.Description)
}

func TestNpmDependenciesAnalyzeWithVersion(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	// Test with a specific version
	input := npmPackageInput{
		PackageName: "lodash",
		Version:     "4.17.21",
	}

	_, output, err := NpmDependenciesAnalyze(ctx, req, input)
	if err != nil {
		t.Fatalf("Failed to analyze lodash package: %v", err)
	}

	if output == nil {
		t.Fatal("Output should not be nil")
	}

	if output.Name != "lodash" {
		t.Errorf("Expected package name 'lodash', got '%s'", output.Name)
	}

	if output.Version != "4.17.21" {
		t.Errorf("Expected version '4.17.21', got '%s'", output.Version)
	}

	t.Logf("Package: %s@%s", output.Name, output.Version)
	t.Logf("Latest: %s", output.LatestVersion)
}

func TestNpmDependenciesAnalyzeNotFound(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	// Test with a non-existent package
	input := npmPackageInput{
		PackageName: "this-package-definitely-does-not-exist-12345",
	}

	_, _, err := NpmDependenciesAnalyze(ctx, req, input)
	if err == nil {
		t.Error("Expected error for non-existent package, got nil")
	}
}

func TestNpmDependenciesAnalyzeEmptyPackageName(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	// Test with empty package name
	input := npmPackageInput{
		PackageName: "",
	}

	_, _, err := NpmDependenciesAnalyze(ctx, req, input)
	if err == nil {
		t.Error("Expected error for empty package name, got nil")
	}
}

func TestNpmDependenciesAnalyzeScopedPackage(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	// Test with a scoped package
	input := npmPackageInput{
		PackageName: "@types/node",
	}

	_, output, err := NpmDependenciesAnalyze(ctx, req, input)
	if err != nil {
		t.Fatalf("Failed to analyze @types/node package: %v", err)
	}

	if output == nil {
		t.Fatal("Output should not be nil")
	}

	if output.Name != "@types/node" {
		t.Errorf("Expected package name '@types/node', got '%s'", output.Name)
	}

	t.Logf("Scoped Package: %s@%s", output.Name, output.Version)
}
