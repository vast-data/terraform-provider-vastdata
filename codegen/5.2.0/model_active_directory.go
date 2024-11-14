/*
 * VastData API
 *
 * A API document representing VastData API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type ActiveDirectory struct {
	Id int32 `json:"id,omitempty"`
	Guid string `json:"guid,omitempty"`
	MachineAccountName string `json:"machine_account_name,omitempty"`
	OrganizationalUnit string `json:"organizational_unit,omitempty"`
	LdapId int64 `json:"ldap_id,omitempty"`
}