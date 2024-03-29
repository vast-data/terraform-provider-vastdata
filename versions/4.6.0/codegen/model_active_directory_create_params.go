/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ActiveDirectoryCreateParams struct {
	// ID of the LDAP configuration for binding to the LDAP domain of the AD server.
	LdapId string `json:"ldap_id,omitempty"`
	// The fully qualified domain name (FQDN) of the Active Directory domain to join.
	DomainName string `json:"domain_name,omitempty"`
	// The name for the machine object representing the VAST Cluster to be created within the OU
	MachineAccountName string `json:"machine_account_name,omitempty"`
	// A non default organizational unit (OU) in the AD domain in which to create the machine object. If left empty, the machine object will be created in the default Computers OU.
	OrganizationalUnit string `json:"organizational_unit,omitempty"`
	// Specify multiple DCs using 'urls' parameter in LDAP configuration.
	PreferredDcList []ErrorUnknown `json:"preferred_dc_list,omitempty"`
	// Indicates if AD is allowed for SMB. There may only be 1 such AD.
	SmbAllowed bool `json:"smb_allowed,omitempty"`
}
