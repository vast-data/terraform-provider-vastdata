/*
 * VastData API
 *
 * A API document representing VastData API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type Role struct {
	// A unique id given to the role
	Id int64 `json:"id,omitempty"`
	// A uniqe name of the role
	Name string `json:"name,omitempty"`
	// List of allowed permissions
	PermissionsList []string `json:"permissions_list,omitempty"`
	// List of allowed permissions returned from the VMS
	Permissions []string `json:"permissions,omitempty"`
	// List of tenants to which this role is associated with
	Tenants []int64 `json:"tenants,omitempty"`
	// Is the role is an admin role
	IsAdmin bool `json:"is_admin,omitempty"`
	// Is the role is a default role
	IsDefault bool `json:"is_default,omitempty"`
	// LDAP group(s) associated with the role. Members of the specified groups on a connected LDAP/Active Directory provider can access VMS and are granted whichever permissions are included in the role. A group can be associated with multiple roles.
	LdapGroups []string `json:"ldap_groups,omitempty"`
}
