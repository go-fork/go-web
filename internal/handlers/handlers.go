// Package handlers provides HTTP request handlers for the Go web application.
// It includes handlers for serving static pages, package listings, and category views.
package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/go-fork/go-web/internal/services"
)

// Handler manages HTTP request handlers with dependencies
type Handler struct {
	packageService *services.PackageService
	templates      *template.Template
}

// NewHandler creates a new handler instance with dependencies
func NewHandler(packageService *services.PackageService) *Handler {
	// Load templates
	templates := template.Must(template.ParseGlob("templates/*.gohtml"))
	return &Handler{
		packageService: packageService,
		templates:      templates,
	}
}

// Home renders the home page with package listing
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	packages := h.packageService.GetAllPackages()

	data := struct {
		Title    string
		Packages []services.Package
	}{
		Title:    "Go Fork - Package List",
		Packages: packages,
	}

	if err := h.templates.ExecuteTemplate(w, "packages.gohtml", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RouteRequest is the main router function that dispatches requests
func (h *Handler) RouteRequest(w http.ResponseWriter, r *http.Request) {
	// For root path, show all packages
	if r.URL.Path == "/" {
		h.AllPackages(w, r)
		return
	}

	// Check if this might be a package path
	path := strings.TrimPrefix(r.URL.Path, "/")

	// Check for go-get parameter
	goGet := r.URL.Query().Get("go-get")

	// If go-get=1 is present, handle go-get request
	if goGet == "1" {
		// Find a matching package
		pkg := h.packageService.GetPackageByPartialPath(path)
		if pkg != nil {
			data := struct {
				Title   string
				Package services.Package
			}{
				Title:   pkg.Path + " - Go Fork",
				Package: *pkg,
			}

			// Render go-get template
			if err := h.templates.ExecuteTemplate(w, "go-get.gohtml", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	// For everything else, return 404
	h.NotFound(w, r)
}

// AllPackages renders the all packages page
func (h *Handler) AllPackages(w http.ResponseWriter, r *http.Request) {
	packages := h.packageService.GetAllPackages()

	data := struct {
		Title    string
		Packages []services.Package
	}{
		Title:    "All Packages - Go Fork",
		Packages: packages,
	}

	if err := h.templates.ExecuteTemplate(w, "packages.gohtml", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// PackageDetail handles package details/go-get views
func (h *Handler) PackageDetail(w http.ResponseWriter, r *http.Request) {
	// Get path directly without /package/ prefix
	packagePath := strings.TrimPrefix(r.URL.Path, "/")
	if packagePath == "" {
		h.NotFound(w, r)
		return
	}

	// Check for go-get query parameter
	goGet := r.URL.Query().Get("go-get")

	// Find package
	pkg := h.packageService.GetPackageByPartialPath(packagePath)
	if pkg == nil {
		// If package not found, return 404
		h.NotFound(w, r)
		return
	}

	data := struct {
		Title   string
		Package services.Package
	}{
		Title:   pkg.Path + " - Go Fork",
		Package: *pkg,
	}

	// Use go-get template for go modules
	if goGet == "1" {
		if err := h.templates.ExecuteTemplate(w, "go-get.gohtml", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Otherwise, redirect to pkg.go.dev
	http.Redirect(w, r, "https://pkg.go.dev/"+pkg.Path, http.StatusFound)
}

// CategoryPackages handles category-based package listing
func (h *Handler) CategoryPackages(w http.ResponseWriter, r *http.Request) {
	// This is a placeholder - for future implementation
	// Currently redirects to all packages
	http.Redirect(w, r, "/packages", http.StatusFound)
}

// NotFound renders the 404 page
func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
	}{
		Title: "Page Not Found - Go Fork",
	}

	w.WriteHeader(http.StatusNotFound)
	if err := h.templates.ExecuteTemplate(w, "404.gohtml", data); err != nil {
		w.Write([]byte("Page Not Found"))
	}
}
