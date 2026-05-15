package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

// HandleSecurityDocs returns HTML documentation for configuring Cosign and Trivy integrations
func (h *Handlers) HandleSecurityDocs(w http.ResponseWriter, r *http.Request) {
	var html strings.Builder

	html.WriteString(`
<div class="space-y-6">
	<div class="bg-indigo-50 dark:bg-indigo-900/30 border border-indigo-200 dark:border-indigo-800 rounded-lg p-5">
		<h3 class="text-lg font-semibold text-indigo-900 dark:text-indigo-100 mb-2">Goal</h3>
		<p class="text-sm text-indigo-800 dark:text-indigo-200">
			Wire external Cosign and Trivy tooling into The Engine without rebuilding the backend. Configure the paths, keys, and policies your platform team already manages.
		</p>
	</div>

	<div class="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg">
		<div class="border-b border-gray-200 dark:border-gray-700 px-5 py-3 flex items-center justify-between">
			<h4 class="font-semibold text-gray-900 dark:text-white">Configuration Blueprint</h4>
			<span class="text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400">config.yaml</span>
		</div>
		<pre class="overflow-auto text-sm leading-relaxed bg-gray-900 text-gray-100 dark:bg-black/60 px-5 py-4 rounded-b-lg">
security:
  cosign:
    enabled: true
    publicKeyPath: "/etc/engine/cosign.pub"
    enableRekor: true
    rekorURL: "https://rekor.sigstore.dev"
  trivy:
    enabled: true
    binaryPath: "/usr/local/bin/trivy"
    cacheDir: "/var/lib/the-engine/trivy-cache"
    severities: ["CRITICAL", "HIGH"]
    timeout: "5m"
    additionalArgs: []
		</pre>
	</div>

	<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
		<div class="border border-gray-200 dark:border-gray-700 rounded-lg p-5 bg-white dark:bg-gray-900">
			<h4 class="font-semibold text-gray-900 dark:text-white mb-3">Cosign Install Checklist</h4>
			<ol class="list-decimal list-inside space-y-2 text-sm text-gray-700 dark:text-gray-300">
				<li>Install Cosign: <code class="bg-gray-100 dark:bg-gray-800 px-1 py-0.5 rounded">curl -sSL https://raw.githubusercontent.com/sigstore/cosign/main/install.sh | sudo sh -s -- -b /usr/local/bin</code></li>
				<li>Fetch or create signing public key and mount it where The Engine can read it.</li>
				<li>Allow outbound HTTPS access to <code>rekor.sigstore.dev</code> if Rekor verification is enabled.</li>
				<li>Grant the service account read access to container registries that store signatures.</li>
			</ol>
		</div>
		<div class="border border-gray-200 dark:border-gray-700 rounded-lg p-5 bg-white dark:bg-gray-900">
			<h4 class="font-semibold text-gray-900 dark:text-white mb-3">Trivy Install Checklist</h4>
			<ol class="list-decimal list-inside space-y-2 text-sm text-gray-700 dark:text-gray-300">
				<li>Install Trivy: <code class="bg-gray-100 dark:bg-gray-800 px-1 py-0.5 rounded">brew install aquasecurity/trivy/trivy</code> or download the binary for Linux.</li>
				<li>Create a writable cache directory and grant the Engine process access.</li>
				<li>Schedule periodic <code>trivy --download-db-only</code> to keep the CVE database fresh.</li>
				<li>Ensure the Engine host can reach target container registries for image pulls.</li>
			</ol>
		</div>
	</div>

	<div class="border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-900 p-5">
		<h4 class="font-semibold text-gray-900 dark:text-white mb-3">Environment Variables</h4>
		<table class="w-full text-sm text-left text-gray-700 dark:text-gray-300">
			<thead class="text-xs uppercase bg-gray-50 dark:bg-gray-800 text-gray-500 dark:text-gray-400">
				<tr>
					<th class="px-3 py-2">Variable</th>
					<th class="px-3 py-2">Purpose</th>
				</tr>
			</thead>
			<tbody>
				<tr class="border-t border-gray-200 dark:border-gray-700">
					<td class="px-3 py-2 font-mono text-xs">ENGINE_SECURITY_COSIGN_KEY</td>
					<td class="px-3 py-2">Override the public key path using an environment variable.</td>
				</tr>
				<tr class="border-t border-gray-200 dark:border-gray-700">
					<td class="px-3 py-2 font-mono text-xs">ENGINE_SECURITY_TRIVY_BIN</td>
					<td class="px-3 py-2">Override the Trivy binary location at runtime.</td>
				</tr>
				<tr class="border-t border-gray-200 dark:border-gray-700">
					<td class="px-3 py-2 font-mono text-xs">ENGINE_SECURITY_TRIVY_CACHE</td>
					<td class="px-3 py-2">Define a cache directory if not using the YAML configuration.</td>
				</tr>
			</tbody>
		</table>
	</div>

	<div class="border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-900 p-5">
		<h4 class="font-semibold text-gray-900 dark:text-white mb-3">Operator Workflow</h4>
		<ol class="list-decimal list-inside space-y-2 text-sm text-gray-700 dark:text-gray-300">
			<li>Provision Cosign public keys and optional Fulcio/Rekor access.</li>
			<li>Install Trivy on the same host as The Engine UI backend.</li>
			<li>Populate the configuration above via the UI or configuration file.</li>
			<li>Use the "Test Security Integration" action (coming soon) to validate the setup.</li>
		</ol>
		<p class="text-xs text-gray-500 dark:text-gray-400 mt-3">
			The current release surfaces configuration guidance only. Runtime enforcement will invoke the configured binaries in a future update.
		</p>
	</div>
</div>
`)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, html.String())
}
