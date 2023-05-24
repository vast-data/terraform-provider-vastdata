# ReplicationPolicy

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Name** | **string** |  | [optional] [default to null]
**Guid** | **string** | unique identifier | [optional] [default to null]
**ScheduleFrequency** | [**time.Time**](time.Time.md) | schedule frequency, in datetime format | [optional] [default to null]
**ScheduleStartTime** | [**time.Time**](time.Time.md) | Schedule the first restore point after the initial sync | [optional] [default to null]
**ReplicationTarget** | **string** | replication target name | [optional] [default to null]
**ReplicationTargetName** | **string** | replication target name | [optional] [default to null]
**BandwidthLimitationRules** | **string** | bandwith limitation rules | [optional] [default to null]
**Priority** | **string** | low / normal / high | [optional] [default to null]
**VipPool** | **string** | ip pool | [optional] [default to null]
**AwsPreferredStorage** | **string** | Amazon S3 / Amazon Glacier | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


