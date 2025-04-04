/*
 * VastData API
 *
 * A API document representing VastData API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type QosUser struct {
	// User FQDN
	Fqdn string `json:"fqdn,omitempty"`
	IsSid bool `json:"is_sid,omitempty"`
	// The user SID
	SidStr string `json:"sid_str,omitempty"`
	UidOrGid int64 `json:"uid_or_gid,omitempty"`
	// How to display the user
	Label string `json:"label,omitempty"`
	// The user name
	Value string `json:"value,omitempty"`
	// The user login name
	LoginName string `json:"login_name,omitempty"`
	// The user name
	Name string `json:"name,omitempty"`
	// The user type of idetify
	IdentifierType string `json:"identifier_type,omitempty"`
	// The value to use fo the identifier_type
	IdentifierValue string `json:"identifier_value,omitempty"`
}
