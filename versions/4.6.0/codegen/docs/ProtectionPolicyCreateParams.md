# ProtectionPolicyCreateParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** | Do not specify this parameter. | [optional] [default to null]
**Name** | **string** |  | [optional] [default to null]
**Guid** | **string** | Do not specify this parameter. | [optional] [default to null]
**TargetObjectId** | **int32** | ID of the remote peer. Specify ID of a ReplicationTarget (aka S3 replication peer) if clone_type is CLOUD_REPLICATION. Specify the ID of a NativeReplicationRemoteTarget if clone_type is NATIVE_REPLICATION. | [optional] [default to null]
**Frames** | [***interface{}**](interface{}.md) | Defines the schedule for snapshot creation and the local and remote retention policies. | [optional] [default to null]
**Prefix** | **string** | The prefix for names of snapshots created by the policy | [optional] [default to null]
**CloneType** | **string** | Specify the type of data protection. CLOUD_REPLICATION is S3 backup. LOCAL means local snapshots without replication. | [optional] [default to null]
**Indestructible** | **bool** | Indestructible the Protection policy from being deleted | [optional] [default to null]
**BigCatalog** | **bool** | Indicates if Protection Policy will be used for big catalogue. There may only be 1 such policy. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


