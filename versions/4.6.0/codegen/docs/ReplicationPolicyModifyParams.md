# ReplicationPolicyModifyParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | [optional] [default to null]
**ScheduleFrequency** | [**time.Time**](time.Time.md) | schedule frequency, in datetime format | [optional] [default to null]
**ScheduleStartTime** | [**time.Time**](time.Time.md) | Schedule the first restore point after the initial sync | [optional] [default to null]
**ReplicationTarget** | **string** | replication target id | [optional] [default to null]
**BandwidthLimitationRules** | **string** | bandwith limitation rules | [optional] [default to null]
**Priority** | **string** | low / normal / high | [optional] [default to null]
**AwsPreferredStorage** | **string** | Amazon S3 / Amazon Glacier | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


