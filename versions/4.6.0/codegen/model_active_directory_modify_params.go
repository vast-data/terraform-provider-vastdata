/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ActiveDirectoryModifyParams struct {
	// An Active Directory admin user with permission to join the Active Directory server.
	AdminUsername string `json:"admin_username,omitempty"`
	// The password for the specified Active Directory admin user.
	AdminPasswd string `json:"admin_passwd,omitempty"`
	// Set to true to join Active Directory. Set to false to leave Active Directory.
	Enabled bool `json:"enabled,omitempty"`
	// Indicates if AD is allowed for SMB. There may only be 1 such AD.
	SmbAllowed bool `json:"smb_allowed,omitempty"`
}
