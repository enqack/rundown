//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	binary    = "rundown"
	mainPath  = "./cmd/rundown"
	buildPath = "."
)

var (
	goexe = "go"
)

// Build builds the rundown binary for the current platform
func Build() error {
	mg.Deps(Clean)
	fmt.Println("Building rundown...")

	env := map[string]string{
		"CGO_ENABLED": "0",
	}

	return sh.RunWith(env, goexe, "build", "-o", binary, mainPath)
}

// Install installs the rundown binary to $GOPATH/bin
func Install() error {
	fmt.Println("Installing rundown...")

	env := map[string]string{
		"CGO_ENABLED": "0",
	}

	return sh.RunWith(env, goexe, "install", mainPath)
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning...")

	// Remove binary
	if err := sh.Rm(binary); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Remove dist directory
	if err := sh.Rm("dist"); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// Test runs all tests
func Test() error {
	fmt.Println("Running tests...")
	return sh.RunV(goexe, "test", "-v", "./...")
}

// Fmt formats code with gofmt
func Fmt() error {
	fmt.Println("Formatting code...")
	return sh.RunV(goexe, "fmt", "./...")
}

// Lint runs golangci-lint
func Lint() error {
	fmt.Println("Running linter...")
	return sh.RunV("golangci-lint", "run", "./...")
}

// Vet runs go vet
func Vet() error {
	fmt.Println("Running vet...")
	return sh.RunV(goexe, "vet", "./...")
}

// Dev builds and runs rundown with base16 theme
func Dev() error {
	mg.Deps(Build)
	fmt.Println("Running rundown...")
	return sh.RunV("./"+binary, "-t", "base16")
}

// Dist builds binaries for multiple platforms
func Dist() error {
	mg.Deps(Clean)
	fmt.Println("Building distribution binaries...")

	platforms := []struct {
		os   string
		arch string
	}{
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"windows", "amd64"},
	}

	// Create dist directory
	if err := os.MkdirAll("dist", 0755); err != nil {
		return err
	}

	for _, p := range platforms {
		binaryName := fmt.Sprintf("rundown-%s-%s", p.os, p.arch)
		if p.os == "windows" {
			binaryName += ".exe"
		}

		outputPath := filepath.Join("dist", binaryName)

		fmt.Printf("Building %s...\n", binaryName)

		env := map[string]string{
			"CGO_ENABLED": "0",
			"GOOS":        p.os,
			"GOARCH":      p.arch,
		}

		if err := sh.RunWith(env, goexe, "build", "-o", outputPath, mainPath); err != nil {
			return err
		}
	}

	fmt.Println("Distribution binaries built successfully in dist/")
	return nil
}

// Version displays the current Go version
func Version() error {
	fmt.Printf("Go version: %s\n", runtime.Version())
	return sh.RunV(goexe, "version")
}

// Deps ensures all dependencies are downloaded and tidy
func Deps() error {
	fmt.Println("Downloading dependencies...")
	if err := sh.RunV(goexe, "mod", "download"); err != nil {
		return err
	}

	fmt.Println("Tidying go.mod...")
	return sh.RunV(goexe, "mod", "tidy")
}

// All runs format, lint, test, and build
func All() error {
	mg.Deps(Fmt, Lint, Vet, Test, Build)
	fmt.Println("✅ All tasks completed successfully!")
	return nil
}
