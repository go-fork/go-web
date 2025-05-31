package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-fork/go-web/internal/handlers"
	"github.com/go-fork/go-web/internal/services"
	"golang.org/x/net/http2"
)

func main() {
	// Initialize package service
	packageService := services.NewPackageService("data/packages.json")

	// Initialize handlers
	handler := handlers.NewHandler(packageService)

	// Static files - only if public directory exists
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// Routes
	// Main handler that routes all requests
	http.HandleFunc("/", handler.RouteRequest)

	// Get host and port from environment variables or use defaults
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get certificate paths from environment variables
	certFile := os.Getenv("CERT_CRT")
	keyFile := os.Getenv("CERT_KEY")

	// Server configuration
	addr := host + ":" + port

	// Get HTTP/2 configuration from environment variables
	http2EnabledStr := os.Getenv("HTTP2_ENABLED")
	http2Enabled := http2EnabledStr == "" || http2EnabledStr == "true" // Default true if not specified
	alpnProtocols := os.Getenv("ALPN_PROTOCOLS")

	// Default ALPN protocols if not specified
	nextProtos := []string{"h2", "http/1.1"}
	if alpnProtocols != "" {
		nextProtos = strings.Split(alpnProtocols, ",")
	}

	// Check if TLS is configured
	if certFile != "" && keyFile != "" {
		// Configure TLS with modern cipher suites và ALPN support
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
			NextProtos: nextProtos, // Use configured ALPN protocols
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			},
			PreferServerCipherSuites: true,
		}

		// Tạo middleware handler để thêm các security headers
		secureHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			http.DefaultServeMux.ServeHTTP(w, r)
		})

		// Create an HTTP server with the TLS config và optimized settings
		server := &http.Server{
			Addr:      addr,
			TLSConfig: tlsConfig,
			Handler:   secureHandler,

			// Thiết lập timeout
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		}

		// Configure HTTP/2 with optimized settings
		http2Config := &http2.Server{
			// Tăng số lượng stream đồng thời tối đa cho mỗi kết nối
			MaxConcurrentStreams: 250,

			// Tăng kích thước bảng header để tối ưu nén header
			MaxDecoderHeaderTableSize: 4096,
			MaxEncoderHeaderTableSize: 4096,

			// Tối ưu hóa kích thước frame để cân bằng hiệu suất
			MaxReadFrameSize: 1 << 20, // 1MB

			// Thiết lập timeout để xử lý các kết nối không hoạt động
			IdleTimeout:     120 * time.Second,
			ReadIdleTimeout: 60 * time.Second,
			PingTimeout:     15 * time.Second,

			// Tối ưu hóa buffer upload để cải thiện hiệu suất
			MaxUploadBufferPerConnection: 1 << 21, // 2MB
			MaxUploadBufferPerStream:     1 << 20, // 1MB
		}

		// Chỉ cấu hình HTTP/2 nếu được bật trong cấu hình
		if http2Enabled {
			if err := http2.ConfigureServer(server, http2Config); err != nil {
				log.Fatalf("Failed to configure HTTP/2: %v", err)
			}
			log.Printf("Starting HTTPS server with HTTP/2 support on %s (ALPN: %s)", addr, strings.Join(nextProtos, ", "))
		} else {
			log.Printf("Starting HTTPS server on %s (HTTP/2 disabled)", addr)
		}

		log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
	} else {
		// Fallback to HTTP if no TLS certificates are configured
		log.Printf("Starting HTTP server on %s (no TLS certificates configured)", addr)
		log.Fatal(http.ListenAndServe(addr, nil))
	}
}
