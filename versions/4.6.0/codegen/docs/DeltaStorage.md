# DeltaStorage

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CurrentSequence** | **string** | Runtime sequence id for delta. Reset if HA happened. | [default to null]
**CurrentGeneration** | **int32** |  | [default to null]
**Records** | [**[]DeltaRecord**](DeltaRecord.md) |  | [default to null]
**Status** | **string** | ok - all ok, reset - journal was reset | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


