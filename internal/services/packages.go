package services

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

// Package represents a Go package with its path, repository, description, and display information
type Package struct {
	Path        string `json:"path"`
	Repository  string `json:"repository"`
	Description string `json:"description"`
	Display     bool   `json:"display"`
}

// PackageData represents the structure of the packages.json file
type PackageData struct {
	Packages []Package `json:"packages"`
}

// PackageService provides methods to interact with package data
type PackageService struct {
	packages     []Package
	packagesByID map[string]*Package
}

// NewPackageService creates a new PackageService instance and loads packages from JSON file
func NewPackageService(path string) *PackageService {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read package data: %v", err)
	}

	var packageData PackageData
	if err := json.Unmarshal(data, &packageData); err != nil {
		log.Fatalf("Failed to parse package data: %v", err)
	}

	service := &PackageService{
		packages:     packageData.Packages,
		packagesByID: make(map[string]*Package),
	}

	// Index packages by path for quick lookups
	for i := range service.packages {
		service.packagesByID[service.packages[i].Path] = &service.packages[i]
	}

	return service
}

// GetAllPackages returns all displayable packages
func (s *PackageService) GetAllPackages() []Package {
	var displayablePackages []Package
	for _, pkg := range s.packages {
		// Only include packages where display is true
		// Note: Since all packages in our JSON now have the display field explicitly set,
		// we don't need to handle legacy packages anymore
		if pkg.Display {
			displayablePackages = append(displayablePackages, pkg)
		}
	}
	return displayablePackages
}

// Note: Legacy package handling has been removed as it's no longer needed

// GetPackage retrieves a package by its full path
func (s *PackageService) GetPackage(path string) *Package {
	pkg, exists := s.packagesByID[path]
	if !exists {
		return nil
	}
	return pkg
}

// GetPackageByPartialPath looks up a package by a partial path
// This is useful for go-get handler where users may request with partial paths
func (s *PackageService) GetPackageByPartialPath(partialPath string) *Package {
	// Handle the case where the URL path is just a part of the package path
	// For example, URL path "di" should match "go.fork.vn/di"

	// First check for an exact match
	if pkg, exists := s.packagesByID[partialPath]; exists {
		// Only return visible packages
		if pkg.Display {
			return pkg
		}
		return nil
	}

	// Check for exact match with domain prefix
	domainPrefixedPath := "go.fork.vn/" + partialPath
	if pkg, exists := s.packagesByID[domainPrefixedPath]; exists {
		// Only return visible packages
		if pkg.Display {
			return pkg
		}
		return nil
	}

	// Then look for packages that end with the partial path
	var matches []*Package
	for _, pkg := range s.packages {
		// Skip packages that should not be displayed
		if !pkg.Display {
			continue
		}

		// Check if package path ends with the partial path
		// e.g., "di" should match "go.fork.vn/di"
		pathParts := strings.Split(pkg.Path, "/")
		if len(pathParts) > 0 && pathParts[len(pathParts)-1] == partialPath {
			p := pkg // Create a copy to avoid reference issues
			matches = append(matches, &p)
		} else if strings.HasPrefix(pkg.Path, partialPath) || strings.HasSuffix(pkg.Path, "/"+partialPath) {
			// Also match if the partial path is a prefix or a suffix with a separator
			p := pkg // Create a copy to avoid reference issues
			matches = append(matches, &p)
		}
	}

	// No matches
	if len(matches) == 0 {
		return nil
	}

	// Return the shortest match (most specific)
	shortestMatch := matches[0]
	shortestLen := len(shortestMatch.Path)

	for _, match := range matches {
		if len(match.Path) < shortestLen {
			shortestMatch = match
			shortestLen = len(match.Path)
		}
	}

	return shortestMatch
}
