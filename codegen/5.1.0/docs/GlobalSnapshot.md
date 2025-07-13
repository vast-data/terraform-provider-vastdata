# GlobalSnapshot

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int64** | The unique ID of the global snapshot. | [optional] [default to null]
**Guid** | **string** | The unique GUID of the global snapshot. | [optional] [default to null]
**Name** | **string** | The name of the snapshot. | [optional] [default to null]
**LoaneeTenantId** | **int** | The tenant ID on the destination peer. | [optional] [default to null]
**LoaneeRootPath** | **string** | The path where to store the snapshot on the destination peer. | [optional] [default to null]
**RemoteTargetId** | **int** | The remote replication peering ID. | [optional] [default to null]
**RemoteTargetGuid** | **string** | The remote replication peering GUID. | [optional] [default to null]
**RemoteTargetPath** | **string** | The path on the remote cluster. | [optional] [default to null]
**Enabled** | **bool** | Sets the snapshot to be enabled or disabled. | [optional] [default to null]
**OwnerRootSnapshot** | [***GlobalSnapshotOwnerRootSnapshot**](GlobalSnapshotOwnerRootSnapshot.md) |  | [optional] [default to null]
**OwnerTenant** | [***GlobalSnapshotOwnerTenant**](GlobalSnapshotOwnerTenant.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

