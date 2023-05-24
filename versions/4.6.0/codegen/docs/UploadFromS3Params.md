# UploadFromS3Params

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**S3Url** | **string** | S3 URL to upgrade package. If not provided, will be taken from db | [optional] [default to null]
**SkipPrepare** | **bool** | Skips preparing the cluster for upgrade, including: pre-upgrade validations, copying the bundle to other hosts, and pulling the image on all CNodes. | [optional] [default to null]
**SkipHwCheck** | **bool** | Skips validation of hardware component health. Use with caution since component redundancy is important in NDU. Do not use with OS upgrade. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


