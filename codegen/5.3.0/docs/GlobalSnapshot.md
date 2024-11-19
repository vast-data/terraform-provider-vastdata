# GlobalSnapshot

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int64** | A unique id given to the global snapshot | [optional] [default to null]
**Guid** | **string** | A unique guid given to the global snapshot | [optional] [default to null]
**Name** | **string** | The name of the snapshot | [optional] [default to null]
**LoaneeTenantId** | **int** | The tenant ID of the target | [optional] [default to null]
**LoaneeRootPath** | **string** | The path where to store the snapshot on a Target | [optional] [default to null]
**RemoteTargetId** | **int** | The remote replication peering id | [optional] [default to null]
**RemoteTargetGuid** | **string** | The remote replication peering guid | [optional] [default to null]
**RemoteTargetPath** | **string** | The path on the remote cluster | [optional] [default to null]
**Enabled** | **bool** | Is the snapshot enabled | [optional] [default to null]
**OwnerRootSnapshot** | [***GlobalSnapshotOwnerRootSnapshot**](GlobalSnapshotOwnerRootSnapshot.md) |  | [optional] [default to null]
**OwnerTenant** | [***GlobalSnapshotOwnerTenant**](GlobalSnapshotOwnerTenant.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

