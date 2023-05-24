# ProtectionPolicy

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int64** | A unique id given to the replication peer configuration | [optional] [default to null]
**Guid** | **string** | A unique guid given to the  replication peer configuration | [optional] [default to null]
**Name** | **string** | The name of the replication peer configuration | [optional] [default to null]
**Url** | **string** | Direct link to the replication policy | [optional] [default to null]
**TargetName** | **string** | The target peer name | [optional] [default to null]
**TargetObjectId** | **int** | The id of the target peer | [optional] [default to null]
**Prefix** | **string** | The prefix to be given to the replicated data | [optional] [default to null]
**CloneType** | **string** | The type the replication | [optional] [default to null]
**Frames** | [**[]ProtectionPolicySchedule**](ProtectionPolicySchedule.md) | List of snapshots schedules | [optional] [default to null]
**Indestructible** | **bool** | Is the snapshot indestructable | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

