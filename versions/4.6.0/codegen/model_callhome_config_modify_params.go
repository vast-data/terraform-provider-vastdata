/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type CallhomeConfigModifyParams struct {
	// Set to true to enable periodic sending of bundles to the Support server
	BundleEnabled bool `json:"bundle_enabled,omitempty"`
	// The frequency for sending bundles to the Support server
	BundleInterval int32 `json:"bundle_interval,omitempty"`
	// Set to true to enable system state data to be logged to the Support server
	LogEnabled bool `json:"log_enabled,omitempty"`
	// The frequency for sending system state data to the Support server.
	LogInterval int32 `json:"log_interval,omitempty"`
	// Company name
	Customer string `json:"customer,omitempty"`
	// Site name
	Site string `json:"site,omitempty"`
	// Site location
	Location string `json:"location,omitempty"`
	// Proxy IP/hostname
	ProxyHost string `json:"proxy_host,omitempty"`
	// Proxy Port
	ProxyPort string `json:"proxy_port,omitempty"`
	// Proxy username
	ProxyUsername string `json:"proxy_username,omitempty"`
	// Proxy password
	ProxyPassword string `json:"proxy_password,omitempty"`
	// Set to true to enable test mode
	TestMode bool `json:"test_mode,omitempty"`
	// Set to true to enable SSL verification. Set to false to disable. VAST Cluster recognizes SSL certificates from a large range of widely recognized certificate authorities (CAs). VAST Cluster may not recognize an SSL certificate signed by your own in-house CA.
	VerifySsl bool `json:"verify_ssl,omitempty"`
	// 
	ProxyScheme string `json:"proxy_scheme,omitempty"`
	// Set to true to enable the VAST Support channel.
	SupportChannel bool `json:"support_channel,omitempty"`
	// Set to true to enable reporting to VAST Cloud Services
	CloudEnabled bool `json:"cloud_enabled,omitempty"`
	// Cloud Services API key
	CloudApiKey string `json:"cloud_api_key,omitempty"`
	//  Cloud Services API domain name
	CloudApiDomain string `json:"cloud_api_domain,omitempty"`
	// Cloud Services subdomain, unique per customer, common to all reporting clusters
	CloudSubdomain string `json:"cloud_subdomain,omitempty"`
	// The ID issued to the customer
	CustomerId string `json:"customer_id,omitempty"`
	// The maximum number of parts of a file to upload simultaneously.
	MaxUploadConcurrency int32 `json:"max_upload_concurrency,omitempty"`
	// If true, call home data is obfuscated.
	Obfuscated bool `json:"obfuscated,omitempty"`
	// If true, send aggregated callhome logs, otherwise upload logs from each node
	Aggregated bool `json:"aggregated,omitempty"`
	// If true, upload non-aggregated Callhome Bundle via VMS (requires proxy). Otherwise, upload from each node.
	UploadViaVms bool `json:"upload_via_vms,omitempty"`
	// Compression method for callhome bundles (by default zstd)
	CompressMethod string `json:"compress_method,omitempty"`
}
