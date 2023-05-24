# Snapshot

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Guid** | **string** |  | [optional] [default to null]
**Name** | **string** |  | [default to null]
**Path** | **string** | Snapshot path | [optional] [default to null]
**ExpirationTime** | **string** | Snapshot expiration time UTC | [optional] [default to null]
**State** | **string** | Snapshot stats | [optional] [default to null]
**Policy** | **string** | Associated snapshot policy | [optional] [default to null]
**PolicyId** | **int32** | Associated snapshot policy ID | [optional] [default to null]
**Cluster** | **string** | Parent Cluster | [optional] [default to null]
**ClusterId** | **int32** | Parent Cluster ID | [optional] [default to null]
**Handle** | **string** | Parent handle | [optional] [default to null]
**Created** | **string** | Snapshot created time | [optional] [default to null]
**Locked** | **bool** | Lock the snapshot from being deleted by cleanup | [optional] [default to null]
**CloneId** | **int32** |  | [optional] [default to null]
**AggrPhysEstimation** | **int64** | The usable capacity reclaimable by deleting the snapshot and all older snapshots on the protected path | [optional] [default to null]
**UniquePhysEstimation** | **int64** | The usable capacity reclaimable by deleting the snapshot without deleting other snapshots on the path | [optional] [default to null]
**ProtectionPolicyId** | **int32** |  | [optional] [default to null]
**ProtectionPolicy** | **string** | Protection Policy Name | [optional] [default to null]
**Type_** | **string** |  | [optional] [default to null]
**Indestructible** | **bool** | Prevent the snapshot from being deleted | [optional] [default to null]
**TenantId** | **int32** | Tenant ID | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


