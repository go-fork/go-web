package models

// Package represents a Go package configuration
type Package struct {
	Name          string            `json:"name"`
	Path          string            `json:"path"`
	Repository    string            `json:"repository"`
	Description   string            `json:"description"`
	Category      string            `json:"category"`
	Version       string            `json:"version,omitempty"`
	Author        string            `json:"author,omitempty"`
	License       string            `json:"license,omitempty"`
	Tags          []string          `json:"tags,omitempty"`
	Dependencies  []string          `json:"dependencies,omitempty"`
	Requires      map[string]string `json:"requires,omitempty"`
	Homepage      string            `json:"homepage,omitempty"`
	Documentation string            `json:"documentation,omitempty"`
}

// PackageConfig represents the entire package configuration
type PackageConfig struct {
	Packages   []Package           `json:"packages"`
	Categories map[string]Category `json:"categories,omitempty"`
}

// Category represents a package category
type Category struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon,omitempty"`
	Color       string `json:"color,omitempty"`
}
