package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// SecurityConfig represents security integration configuration
type SecurityConfig struct {
	CosignEnabled     bool   `json:"cosign_enabled"`
	CosignKeyPath     string `json:"cosign_key_path"`
	CosignEnableRekor bool   `json:"cosign_enable_rekor"`
	CosignRekorURL    string `json:"cosign_rekor_url"`
	
	TrivyEnabled      bool     `json:"trivy_enabled"`
	TrivyBinaryPath   string   `json:"trivy_binary_path"`
	TrivyCacheDir     string   `json:"trivy_cache_dir"`
	TrivySeverities   []string `json:"trivy_severities"`
	TrivyTimeout      string   `json:"trivy_timeout"`
}

// HandleSecurityConfigForm returns HTML form for security configuration
func (h *Handlers) HandleSecurityConfigForm(w http.ResponseWriter, r *http.Request) {
	// Load current config (defaults)
	config := h.loadSecurityConfig()
	
	var html strings.Builder
	
	html.WriteString(`
<div class="space-y-6">
	<div class="bg-blue-50 dark:bg-blue-900/30 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
		<p class="text-sm text-blue-800 dark:text-blue-200">
			Configure your Cosign and Trivy integrations below. These settings will be used when the security verification features are enabled.
		</p>
	</div>
	
	<form hx-post="/api/security/config" hx-target="#security-config-result" hx-swap="innerHTML" class="space-y-6">
		<!-- Cosign Configuration -->
		<div class="border border-gray-200 dark:border-gray-700 rounded-lg p-5 bg-white dark:bg-gray-900">
			<h4 class="font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
				<svg class="w-5 h-5 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path>
				</svg>
				Cosign Signature Verification
			</h4>
			
			<div class="space-y-4">
				<div class="flex items-center gap-3">
					<input type="checkbox" id="cosign-enabled" name="cosign_enabled" `)
	if config.CosignEnabled {
		html.WriteString(`checked`)
	}
	html.WriteString(` class="w-4 h-4 text-purple-600 rounded border-gray-300 focus:ring-purple-500">
					<label for="cosign-enabled" class="text-sm font-medium text-gray-700 dark:text-gray-300">Enable Cosign Verification</label>
				</div>
				
				<div>
					<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Public Key Path</label>
					<input type="text" name="cosign_key_path" value="` + config.CosignKeyPath + `"
						placeholder="/etc/security/cosign_pubkey.pem"
						class="w-full bg-gray-50 dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-gray-900 dark:text-white focus:ring-purple-500 focus:border-purple-500">
					<p class="text-xs text-gray-500 dark:text-gray-400 mt-1">Path to your Cosign public key PEM file</p>
				</div>
				
				<div class="flex items-center gap-3">
					<input type="checkbox" id="cosign-rekor" name="cosign_enable_rekor" `)
	if config.CosignEnableRekor {
		html.WriteString(`checked`)
	}
	html.WriteString(` class="w-4 h-4 text-purple-600 rounded border-gray-300 focus:ring-purple-500">
					<label for="cosign-rekor" class="text-sm font-medium text-gray-700 dark:text-gray-300">Enable Rekor Transparency Log</label>
				</div>
				
				<div>
					<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Rekor URL</label>
					<input type="text" name="cosign_rekor_url" value="` + config.CosignRekorURL + `"
						placeholder="https://rekor.sigstore.dev"
						class="w-full bg-gray-50 dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-gray-900 dark:text-white focus:ring-purple-500 focus:border-purple-500">
				</div>
			</div>
		</div>
		
		<!-- Trivy Configuration -->
		<div class="border border-gray-200 dark:border-gray-700 rounded-lg p-5 bg-white dark:bg-gray-900">
			<h4 class="font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
				<svg class="w-5 h-5 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"></path>
				</svg>
				Trivy Vulnerability Scanning
			</h4>
			
			<div class="space-y-4">
				<div class="flex items-center gap-3">
					<input type="checkbox" id="trivy-enabled" name="trivy_enabled" `)
	if config.TrivyEnabled {
		html.WriteString(`checked`)
	}
	html.WriteString(` class="w-4 h-4 text-green-600 rounded border-gray-300 focus:ring-green-500">
					<label for="trivy-enabled" class="text-sm font-medium text-gray-700 dark:text-gray-300">Enable Trivy Scanning</label>
				</div>
				
				<div>
					<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Trivy Binary Path</label>
					<input type="text" name="trivy_binary_path" value="` + config.TrivyBinaryPath + `"
						placeholder="/usr/local/bin/trivy"
						class="w-full bg-gray-50 dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-gray-900 dark:text-white focus:ring-green-500 focus:border-green-500">
				</div>
				
				<div>
					<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Cache Directory</label>
					<input type="text" name="trivy_cache_dir" value="` + config.TrivyCacheDir + `"
						placeholder="/var/lib/trivy"
						class="w-full bg-gray-50 dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-gray-900 dark:text-white focus:ring-green-500 focus:border-green-500">
				</div>
				
				<div>
					<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Severity Levels</label>
					<div class="flex flex-wrap gap-2">
						<label class="flex items-center gap-2">
							<input type="checkbox" name="trivy_severities" value="CRITICAL" `)
	if contains(config.TrivySeverities, "CRITICAL") {
		html.WriteString(`checked`)
	}
	html.WriteString(` class="w-4 h-4 text-green-600 rounded border-gray-300 focus:ring-green-500">
							<span class="text-sm text-gray-700 dark:text-gray-300">CRITICAL</span>
						</label>
						<label class="flex items-center gap-2">
							<input type="checkbox" name="trivy_severities" value="HIGH" `)
	if contains(config.TrivySeverities, "HIGH") {
		html.WriteString(`checked`)
	}
	html.WriteString(` class="w-4 h-4 text-green-600 rounded border-gray-300 focus:ring-green-500">
							<span class="text-sm text-gray-700 dark:text-gray-300">HIGH</span>
						</label>
						<label class="flex items-center gap-2">
							<input type="checkbox" name="trivy_severities" value="MEDIUM" `)
	if contains(config.TrivySeverities, "MEDIUM") {
		html.WriteString(`checked`)
	}
	html.WriteString(` class="w-4 h-4 text-green-600 rounded border-gray-300 focus:ring-green-500">
							<span class="text-sm text-gray-700 dark:text-gray-300">MEDIUM</span>
						</label>
						<label class="flex items-center gap-2">
							<input type="checkbox" name="trivy_severities" value="LOW" `)
	if contains(config.TrivySeverities, "LOW") {
		html.WriteString(`checked`)
	}
	html.WriteString(` class="w-4 h-4 text-green-600 rounded border-gray-300 focus:ring-green-500">
							<span class="text-sm text-gray-700 dark:text-gray-300">LOW</span>
						</label>
					</div>
				</div>
				
				<div>
					<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Scan Timeout</label>
					<input type="text" name="trivy_timeout" value="` + config.TrivyTimeout + `"
						placeholder="5m"
						class="w-full bg-gray-50 dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-gray-900 dark:text-white focus:ring-green-500 focus:border-green-500">
					<p class="text-xs text-gray-500 dark:text-gray-400 mt-1">Duration format (e.g., 5m, 10m)</p>
				</div>
			</div>
		</div>
		
		<div class="flex gap-3">
			<button type="submit" class="bg-purple-600 hover:bg-purple-700 text-white font-medium px-6 py-2 rounded-lg transition-colors">
				Save Configuration
			</button>
			<button type="button" hx-get="/api/security/docs" hx-target="#security-config-result" hx-swap="innerHTML"
				class="bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 font-medium px-6 py-2 rounded-lg transition-colors">
				View Documentation
			</button>
		</div>
	</form>
	
	<div id="security-config-result" class="mt-4"></div>
</div>
`)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, html.String())
}

// HandleSecurityConfigSave handles saving security configuration
func (h *Handlers) HandleSecurityConfigSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Parse form data
	r.ParseForm()
	
	config := SecurityConfig{
		CosignEnabled:     r.FormValue("cosign_enabled") == "on",
		CosignKeyPath:     r.FormValue("cosign_key_path"),
		CosignEnableRekor: r.FormValue("cosign_enable_rekor") == "on",
		CosignRekorURL:    r.FormValue("cosign_rekor_url"),
		TrivyEnabled:      r.FormValue("trivy_enabled") == "on",
		TrivyBinaryPath:   r.FormValue("trivy_binary_path"),
		TrivyCacheDir:     r.FormValue("trivy_cache_dir"),
		TrivySeverities:   r.Form["trivy_severities"],
		TrivyTimeout:      r.FormValue("trivy_timeout"),
	}
	
	// Save config (in memory for now, could be persisted to file/database)
	h.cache.Set("security_config", config)
	
	// Return success message
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `
<div class="bg-green-50 dark:bg-green-900/30 border border-green-200 dark:border-green-800 rounded-lg p-4">
	<div class="flex items-center gap-2">
		<svg class="w-5 h-5 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
		</svg>
		<span class="font-medium text-green-800 dark:text-green-200">Configuration saved successfully</span>
	</div>
	<p class="text-sm text-green-700 dark:text-green-300 mt-2">
		Your security integration settings have been saved. The configuration will be used when security verification is enabled.
	</p>
</div>
`)
}

// loadSecurityConfig loads security configuration from cache or environment variables
func (h *Handlers) loadSecurityConfig() SecurityConfig {
	// Try to load from cache first
	if cached, found := h.cache.Get("security_config"); found {
		if config, ok := cached.(SecurityConfig); ok {
			return config
		}
	}
	
	// Load from environment variables with defaults
	config := SecurityConfig{
		CosignEnabled:     os.Getenv("ENGINE_SECURITY_COSIGN_ENABLED") == "true",
		CosignKeyPath:     getEnvWithDefault("ENGINE_SECURITY_COSIGN_KEY", "/etc/security/cosign_pubkey.pem"),
		CosignEnableRekor: os.Getenv("ENGINE_SECURITY_COSIGN_REKOR") == "true",
		CosignRekorURL:    getEnvWithDefault("ENGINE_SECURITY_COSIGN_REKOR_URL", "https://rekor.sigstore.dev"),
		TrivyEnabled:      os.Getenv("ENGINE_SECURITY_TRIVY_ENABLED") == "true",
		TrivyBinaryPath:   getEnvWithDefault("ENGINE_SECURITY_TRIVY_BIN", "/usr/local/bin/trivy"),
		TrivyCacheDir:     getEnvWithDefault("ENGINE_SECURITY_TRIVY_CACHE", "/var/lib/trivy"),
		TrivySeverities:   []string{"CRITICAL", "HIGH"},
		TrivyTimeout:      getEnvWithDefault("ENGINE_SECURITY_TRIVY_TIMEOUT", "5m"),
	}
	
	return config
}

// getEnvWithDefault returns environment variable value or default
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// contains checks if a string slice contains a value
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
