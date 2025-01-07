# Volume

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int64** | A unique id given to the manager | [optional] [default to null]
**Guid** | **string** | A uniqe GUID assigned to the manager | [optional] [default to null]
**Name** | **string** | A uniqe name given to the volume | [optional] [default to null]
**Size** | **int64** | The volume size of the volume in bytes | [optional] [default to null]
**ViewId** | **int64** | The View ID to relate this volume with , must be a View with protocol defined as BLOCK | [optional] [default to null]
**BlockHostIds** | **[]int64** | List of blockhosts associated with this volume | [optional] [default to null]
**VolumeTags** | **[]string** | List of tags at the key:value format | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

