# SupportBundleCreateParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Prefix** | **string** | Specify an identifying label to include in the bundle file name. | [optional] [default to null]
**Level** | **string** | Bundle level e.g. small, medium, large | [optional] [default to null]
**Path** | **string** | Bundle path | [optional] [default to null]
**Aggregated** | **bool** | Saves an aggregated bundle file on the management CNode in addition to the separate bundle files that are saved per node. | [optional] [default to null]
**Preset** | **string** | A predefined preset bundle | [optional] [default to null]
**Obfuscated** | **bool** | Converts all bundled objects to text and obfuscates them. Any data that cannot be converted to text is not included in the bundle. The following types of information are replaced with a non-reversible hash: file and directory names, IP addresses, host names, user names, passwords, MAC addresses. | [optional] [default to null]
**StartTime** | **string** | Start time of logs in UTC+3 | [optional] [default to null]
**EndTime** | **string** | End time of logs in UTC+3 | [optional] [default to null]
**Text** | **bool** | Convert all bundled objects to a textual format. Any data that cannot be converted to text is not included in the bundle. | [optional] [default to null]
**HubbleArgs** | **string** | Arguments for the hubble command. Note: times should be in UTC (Use with caution) | [optional] [default to null]
**AstronArgs** | **string** | Arguments for the astron command. Note: times should be in UTC (Use with caution) | [optional] [default to null]
**CnodeIds** | **string** | Collect from specific CNodes. Specify as a comma separated array of CNode IDs. If not specified, logs are collected from all CNodes. | [optional] [default to null]
**DnodeIds** | **string** | Collect from specific DNodes. Specify as a comma separated array of DNode IDs. If not specified, logs are collected from all DNodes. | [optional] [default to null]
**VippoolIds** | **string** | Collect support bundle from CNodes in these vip-pools IDs | [optional] [default to null]
**CnodesOnly** | **bool** | Collect logs from CNodes only | [optional] [default to null]
**DnodesOnly** | **bool** | Collect logs from DNodes only | [optional] [default to null]
**MaxSize** | **float32** | Maximum data limit to apply to the collection of binary trace files, in GB, per node. | [optional] [default to null]
**SendNow** | **bool** | Upload Support Bundle immediately after creation | [optional] [default to null]
**UploadViaVms** | **bool** | If true, upload non-aggregated Support Bundle via VMS (requires proxy). Otherwise, upload from each node. | [optional] [default to null]
**BucketSubdir** | **string** | Sub-Directory in support bucket | [optional] [default to null]
**AccessKey** | **string** | S3 Bucket access key | [optional] [default to null]
**SecretKey** | **string** | S3 Bucket secret key | [optional] [default to null]
**BucketName** | **string** | S3 Bucket for upload | [optional] [default to null]
**DeleteAfterSend** | **bool** | Delete bundle immediately after successfully uploading | [optional] [default to null]
**MdToolHandles** | **string** | A comma separated list of handles to send to md_tool_cli | [optional] [default to null]
**MdToolAddresses** | **string** | A comma separated list of block addresses to send to md_tool_cli | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


