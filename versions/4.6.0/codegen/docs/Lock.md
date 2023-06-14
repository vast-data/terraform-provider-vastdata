# Lock

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**LockType** | **string** | NLM4, ... | [optional] [default to null]
**Caller** | **string** | An identifier of the client that acquired the lock. This could be an IP or host name of the client. | [optional] [default to null]
**Owner** | **string** | An identifier internal to the client kernel for the specific process that owns the lock. | [optional] [default to null]
**IsExclusive** | **bool** | If true, the lock is an exclusive (write) lock. If false, the lock is a shared (read) lock. | [optional] [default to null]
**CreateTimeNano** | **int32** | The time the lock was acquired. | [optional] [default to null]
**Offset** | **int32** | The number of bytes from the beginning of the file&#39;s byte range from which the lock begins. | [optional] [default to null]
**Length** | **int32** | The number of bytes of the file locked by the lock. A length of 0 means the lock reaches until the end of the file.  | [optional] [default to null]
**Svid** | **int32** | A kernel identifier of the owning process on the client machine. | [optional] [default to null]
**Path** | **string** | The path that the locks are taken on | [optional] [default to null]
**State** | **string** | Lock state | [optional] [default to null]
**LockPath** | **string** | The path that the locks are taken on | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


