package services

import (
	"encoding/json"
	"fmt"
	"os"

	"go.for.vn/go-web/internal/models"
)

type PackageService struct {
	configPath string
	packages   []models.Package
	categories map[string]models.Category
}

func NewPackageService(configPath string) *PackageService {
	return &PackageService{
		configPath: configPath,
	}
}

func (s *PackageService) LoadPackages() error {
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config models.PackageConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	s.packages = config.Packages
	s.categories = config.Categories
	return nil
}

func (s *PackageService) GetAllPackages() []models.Package {
	return s.packages
}

func (s *PackageService) GetAllCategories() map[string]models.Category {
	return s.categories
}

func (s *PackageService) GetPackagesByCategory(category string) []models.Package {
	var packages []models.Package
	for _, pkg := range s.packages {
		if pkg.Category == category {
			packages = append(packages, pkg)
		}
	}
	return packages
}

func (s *PackageService) FindPackage(path string) *models.Package {
	for _, pkg := range s.packages {
		if pkg.Path == path || pkg.Name == path {
			return &pkg
		}
	}
	return nil
}

func (s *PackageService) FindPackageByName(name string) *models.Package {
	for _, pkg := range s.packages {
		if pkg.Name == name {
			return &pkg
		}
	}
	return nil
}

func (s *PackageService) AddPackage(pkg models.Package) error {
	// Check if package already exists
	if s.FindPackage(pkg.Path) != nil {
		return fmt.Errorf("package with path %s already exists", pkg.Path)
	}

	s.packages = append(s.packages, pkg)
	return s.saveConfig()
}

func (s *PackageService) UpdatePackage(name string, pkg models.Package) error {
	for i, existingPkg := range s.packages {
		if existingPkg.Name == name {
			s.packages[i] = pkg
			return s.saveConfig()
		}
	}
	return fmt.Errorf("package %s not found", name)
}

func (s *PackageService) DeletePackage(name string) error {
	for i, pkg := range s.packages {
		if pkg.Name == name {
			s.packages = append(s.packages[:i], s.packages[i+1:]...)
			return s.saveConfig()
		}
	}
	return fmt.Errorf("package %s not found", name)
}

func (s *PackageService) saveConfig() error {
	config := models.PackageConfig{
		Packages:   s.packages,
		Categories: s.categories,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(s.configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
