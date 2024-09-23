# ProtocolsAudit

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CreateDeleteFilesDirsObjects** | **bool** | Audit operations that create or delete files, directories, or objects. | [optional] [default to false]
**LogDeletedFilesDirs** | **bool** | Log deleted files and directories. | [optional] [default to false]
**LogFullPath** | **bool** | Log full Element Store path to the requested resource. Enabled by default. May affect performance. When disabled, the view path is recorded. | [optional] [default to true]
**LogUsername** | **bool** | Log username of requesting user. Disabled by default | [optional] [default to false]
**LogHostname** | **bool** | Log the accessing Hostname | [optional] [default to null]
**ModifyDataMd** | **bool** | Audit operations that modify data (including operations that change the file size) and metadata | [optional] [default to false]
**ReadData** | **bool** | Audit operations that read data and metadata | [optional] [default to false]
**ModifyData** | **bool** |  | [optional] [default to false]
**ReadDataMd** | **bool** |  | [optional] [default to false]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

