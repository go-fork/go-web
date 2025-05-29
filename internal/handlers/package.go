package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"go.for.vn/go-web/internal/models"
	"go.for.vn/go-web/internal/services"
)

type PackageHandler struct {
	packageService *services.PackageService
	templates      *template.Template
}

func NewPackageHandler(packageService *services.PackageService) *PackageHandler {
	// Load templates
	templates, err := template.ParseGlob("public/templates/*.gohtml")
	if err != nil {
		panic(fmt.Sprintf("Failed to load templates: %v", err))
	}

	return &PackageHandler{
		packageService: packageService,
		templates:      templates,
	}
}

// HandleRoot serves the main page with package list or handles module requests
func (h *PackageHandler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		h.handleHomePage(w, r)
		return
	}

	// Handle Go module requests
	h.handleModuleRequest(w, r)
}

type HomePageData struct {
	Title       string
	Packages    []models.Package
	Categories  map[string]models.Category
	Selected    string
	SearchQuery string
	CurrentYear int
	GoGetMeta   *GoGetMeta
}

func (h *PackageHandler) handleHomePage(w http.ResponseWriter, r *http.Request) {
	// Check if filtering by category or searching
	category := r.URL.Query().Get("category")
	searchQuery := r.URL.Query().Get("search")

	var packages []models.Package
	if searchQuery != "" {
		// Search functionality
		allPackages := h.packageService.GetAllPackages()
		searchQuery = strings.ToLower(searchQuery)
		for _, pkg := range allPackages {
			if strings.Contains(strings.ToLower(pkg.Name), searchQuery) ||
				strings.Contains(strings.ToLower(pkg.Description), searchQuery) ||
				strings.Contains(strings.ToLower(pkg.Path), searchQuery) {
				packages = append(packages, pkg)
				continue
			}
			// Search in tags
			for _, tag := range pkg.Tags {
				if strings.Contains(strings.ToLower(tag), searchQuery) {
					packages = append(packages, pkg)
					break
				}
			}
		}
	} else if category != "" {
		packages = h.packageService.GetPackagesByCategory(category)
	} else {
		packages = h.packageService.GetAllPackages()
	}

	categories := h.packageService.GetAllCategories()

	// Determine title
	title := "Trang chủ"
	if category != "" {
		if cat, exists := categories[category]; exists {
			title = cat.Name
		}
	} else if searchQuery != "" {
		title = "Kết quả tìm kiếm: " + searchQuery
	}

	// Get current year for footer copyright
	currentYear := time.Now().Year()

	data := HomePageData{
		Title:       title,
		Packages:    packages,
		Categories:  categories,
		Selected:    category,
		SearchQuery: searchQuery,
		CurrentYear: currentYear,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := h.templates.ExecuteTemplate(w, "home.gohtml", data)
	if err != nil {
		// Chỉ sử dụng log.Printf thay vì http.Error để tránh ghi đè header
		log.Printf("Error executing template: %v", err)
	}
}

type GoGetMeta struct {
	Path       string
	Repository string
}

type PackagePageData struct {
	Title       string
	Package     models.Package
	Categories  map[string]models.Category
	GoGetMeta   *GoGetMeta
	CurrentYear int
	SearchQuery string
	Selected    string
}

func (h *PackageHandler) handleModuleRequest(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	// Find matching package
	var pkg *models.Package
	packages := h.packageService.GetAllPackages()
	for _, p := range packages {
		modulePath := strings.TrimPrefix(p.Path, "go.fork.vn/")
		if path == modulePath || strings.HasPrefix(path, modulePath+"/") || path == p.Name {
			pkg = &p
			break
		}
	}

	if pkg == nil {
		http.NotFound(w, r)
		return
	}

	// Handle different Go module requests
	if strings.Contains(r.URL.RawQuery, "go-get=1") {
		h.handleGoGet(w, r, pkg)
		return
	}

	// For web browser requests, show package info
	h.handlePackageInfo(w, r, pkg)
}

func (h *PackageHandler) handleGoGet(w http.ResponseWriter, r *http.Request, pkg *models.Package) {
	currentYear := time.Now().Year()

	data := PackagePageData{
		Title:   "Go Get - " + pkg.Path,
		Package: *pkg,
		GoGetMeta: &GoGetMeta{
			Path:       pkg.Path,
			Repository: pkg.Repository,
		},
		CurrentYear: currentYear,
		Selected:    pkg.Category, // Thêm category của package làm Selected
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := h.templates.ExecuteTemplate(w, "go-get.gohtml", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *PackageHandler) handlePackageInfo(w http.ResponseWriter, r *http.Request, pkg *models.Package) {
	categories := h.packageService.GetAllCategories()

	currentYear := time.Now().Year()

	data := PackagePageData{
		Title:       pkg.Name + " - " + pkg.Description,
		Package:     *pkg,
		Categories:  categories,
		CurrentYear: currentYear,
		Selected:    pkg.Category, // Thêm category của package làm Selected
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := h.templates.ExecuteTemplate(w, "package.gohtml", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}
