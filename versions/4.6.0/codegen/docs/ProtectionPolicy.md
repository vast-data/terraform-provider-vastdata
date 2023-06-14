# ProtectionPolicy

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Name** | **string** |  | [optional] [default to null]
**Title** | **string** |  | [optional] [default to null]
**Guid** | **string** | unique identifier | [optional] [default to null]
**TargetObjectId** | **int32** | target object id | [optional] [default to null]
**TargetName** | **string** | Target Name | [optional] [default to null]
**Frames** | [***interface{}**](interface{}.md) | schedules | [optional] [default to null]
**Prefix** | **string** | The prefix of the snapshot that will be created | [optional] [default to null]
**CloneType** | **string** | Specify the type of data protection. CLOUD_REPLICATION is S3 backup. LOCAL means local snapshots without replication. (LOCAL | NATIVE_REPLICATION | CLOUD_REPLICATION) | [optional] [default to null]
**State** | **string** | State of Protection Policy | [optional] [default to null]
**Internal** | **bool** |  | [optional] [default to null]
**Indestructible** | **bool** | Indestructible the Protection policy from being deleted | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


