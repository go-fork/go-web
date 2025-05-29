package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.for.vn/go-web/internal/models"
	"go.for.vn/go-web/internal/services"
)

type APIHandler struct {
	packageService *services.PackageService
}

func NewAPIHandler(packageService *services.PackageService) *APIHandler {
	return &APIHandler{
		packageService: packageService,
	}
}

func (h *APIHandler) HandleAPI(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api")

	switch {
	case path == "/packages" && r.Method == http.MethodGet:
		h.handleGetPackages(w, r)
	case path == "/packages" && r.Method == http.MethodPost:
		h.handleCreatePackage(w, r)
	case path == "/search" && r.Method == http.MethodGet:
		h.handleSearchPackages(w, r)
	case path == "/stats" && r.Method == http.MethodGet:
		h.handleGetStats(w, r)
	case path == "/categories" && r.Method == http.MethodGet:
		h.handleGetCategories(w, r)
	case strings.HasPrefix(path, "/packages/") && r.Method == http.MethodGet:
		name := strings.TrimPrefix(path, "/packages/")
		h.handleGetPackage(w, r, name)
	case strings.HasPrefix(path, "/packages/") && r.Method == http.MethodPut:
		name := strings.TrimPrefix(path, "/packages/")
		h.handleUpdatePackage(w, r, name)
	case strings.HasPrefix(path, "/packages/") && r.Method == http.MethodDelete:
		name := strings.TrimPrefix(path, "/packages/")
		h.handleDeletePackage(w, r, name)
	case strings.HasPrefix(path, "/categories/") && r.Method == http.MethodGet:
		category := strings.TrimPrefix(path, "/categories/")
		h.handleGetPackagesByCategory(w, r, category)
	default:
		http.NotFound(w, r)
	}
}

func (h *APIHandler) handleGetPackages(w http.ResponseWriter, r *http.Request) {
	packages := h.packageService.GetAllPackages()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"packages": packages,
		"count":    len(packages),
	})
}

func (h *APIHandler) handleGetPackage(w http.ResponseWriter, r *http.Request, name string) {
	pkg := h.packageService.FindPackageByName(name)
	if pkg == nil {
		http.Error(w, "Package not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pkg)
}

func (h *APIHandler) handleCreatePackage(w http.ResponseWriter, r *http.Request) {
	var pkg models.Package
	if err := json.NewDecoder(r.Body).Decode(&pkg); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.packageService.AddPackage(pkg); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pkg)
}

func (h *APIHandler) handleUpdatePackage(w http.ResponseWriter, r *http.Request, name string) {
	var pkg models.Package
	if err := json.NewDecoder(r.Body).Decode(&pkg); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.packageService.UpdatePackage(name, pkg); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pkg)
}

func (h *APIHandler) handleDeletePackage(w http.ResponseWriter, r *http.Request, name string) {
	if err := h.packageService.DeletePackage(name); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *APIHandler) handleSearchPackages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	packages := h.packageService.GetAllPackages()
	var results []models.Package

	query = strings.ToLower(query)
	for _, pkg := range packages {
		// Search in name, description, tags
		if strings.Contains(strings.ToLower(pkg.Name), query) ||
			strings.Contains(strings.ToLower(pkg.Description), query) ||
			strings.Contains(strings.ToLower(pkg.Path), query) {
			results = append(results, pkg)
			continue
		}

		// Search in tags
		for _, tag := range pkg.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				results = append(results, pkg)
				break
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"query":    query,
		"packages": results,
		"count":    len(results),
	})
}

func (h *APIHandler) handleGetStats(w http.ResponseWriter, r *http.Request) {
	packages := h.packageService.GetAllPackages()
	categories := h.packageService.GetAllCategories()

	stats := map[string]interface{}{
		"total_packages":       len(packages),
		"total_categories":     len(categories),
		"packages_by_category": make(map[string]int),
		"latest_packages":      make([]models.Package, 0),
	}

	// Count packages by category
	categoryStats := make(map[string]int)
	for _, pkg := range packages {
		categoryStats[pkg.Category]++
	}
	stats["packages_by_category"] = categoryStats

	// Get latest 5 packages (assuming they're ordered)
	latestCount := 5
	if len(packages) < latestCount {
		latestCount = len(packages)
	}
	if latestCount > 0 {
		stats["latest_packages"] = packages[len(packages)-latestCount:]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *APIHandler) handleGetCategories(w http.ResponseWriter, r *http.Request) {
	categories := h.packageService.GetAllCategories()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"categories": categories,
		"count":      len(categories),
	})
}

func (h *APIHandler) handleGetPackagesByCategory(w http.ResponseWriter, r *http.Request, category string) {
	packages := h.packageService.GetPackagesByCategory(category)
	categories := h.packageService.GetAllCategories()

	categoryInfo, exists := categories[category]
	if !exists {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"category": categoryInfo,
		"packages": packages,
		"count":    len(packages),
	})
}
