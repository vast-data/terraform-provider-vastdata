# ProtectionPolicy

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int64** | The unique ID of the protection policy. | [optional] [default to null]
**Guid** | **string** | The unique GUID of the protection policy. | [optional] [default to null]
**Name** | **string** | The name of the protection policy. | [optional] [default to null]
**Url** | **string** | Direct URL of the protection policy. | [optional] [default to null]
**TargetName** | **string** | The name of the destination peer. | [optional] [default to null]
**TargetObjectId** | **int** | The ID of the destination peer. | [optional] [default to null]
**Prefix** | **string** | The prefix to be given to the replicated data. | [optional] [default to null]
**CloneType** | **string** | The type of replication. | [optional] [default to null]
**Frames** | [**[]ProtectionPolicySchedule**](ProtectionPolicySchedule.md) | A list of snapshot schedules. | [optional] [default to null]
**Indestructible** | **bool** | If true, the snapshot is  indestructable. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

