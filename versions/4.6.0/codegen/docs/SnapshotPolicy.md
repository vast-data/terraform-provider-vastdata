# SnapshotPolicy

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Guid** | **string** |  | [optional] [default to null]
**Name** | **string** |  | [default to null]
**Path** | **string** | Snapshot source path | [default to null]
**SnapshotExpiration** | **string** | m: minutes, D:days, W:weeks, M:months | [optional] [default to null]
**MaxCreatedSnapshots** | **int32** | The maximum number of snapshots that will be retained locally | [optional] [default to null]
**LastOperationState** | **string** | The state of the last policy operation | [optional] [default to null]
**Schedule** | **string** | The schedule to take the snapshot | [optional] [default to null]
**HumanizeSchedule** | **string** | Humanize readable schedule | [optional] [default to null]
**Prefix** | **string** | The prefix of the snapshot that will be create | [optional] [default to null]
**Cluster** | **string** | Parent Cluster | [optional] [default to null]
**ClusterId** | **int32** | Parent Cluster ID | [optional] [default to null]
**Handle** | **string** | Parent handle | [optional] [default to null]
**Enabled** | **bool** | Enable the policy | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


