/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Lock struct {
	// NLM4, ...
	LockType string `json:"lock_type,omitempty"`
	// An identifier of the client that acquired the lock. This could be an IP or host name of the client.
	Caller string `json:"caller,omitempty"`
	// An identifier internal to the client kernel for the specific process that owns the lock.
	Owner string `json:"owner,omitempty"`
	// If true, the lock is an exclusive (write) lock. If false, the lock is a shared (read) lock.
	IsExclusive bool `json:"is_exclusive,omitempty"`
	// The time the lock was acquired.
	CreateTimeNano int32 `json:"create_time_nano,omitempty"`
	// The number of bytes from the beginning of the file's byte range from which the lock begins.
	Offset int32 `json:"offset,omitempty"`
	// The number of bytes of the file locked by the lock. A length of 0 means the lock reaches until the end of the file. 
	Length int32 `json:"length,omitempty"`
	// A kernel identifier of the owning process on the client machine.
	Svid int32 `json:"svid,omitempty"`
	// The path that the locks are taken on
	Path string `json:"path,omitempty"`
	// Lock state
	State string `json:"state,omitempty"`
	// The path that the locks are taken on
	LockPath string `json:"lock_path,omitempty"`
}