# SnapshotCreateParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Snapshot name | [default to null]
**Path** | **string** | The path to take a snapshot on | [default to null]
**ExpirationTime** | **string** | Snapshot expiration time | [optional] [default to null]
**ClusterId** | **int32** | Cluster ID | [optional] [default to null]
**Locked** | **bool** | Protect the snapshot from deletion. If locked, a snapshot cannot expire or be deleted without being unlocked first. | [optional] [default to null]
**Indestructible** | **bool** | Prevent the snapshot from being deleted | [optional] [default to null]
**TenantId** | **int32** | Tenant ID | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


