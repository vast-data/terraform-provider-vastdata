# ProtectionPolicySchedule

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Every** | **string** | How often to create a snapshot. The format is &lt;integer&gt;&lt;time period&gt;. The time period can be D - Days, W - Weeks, s - Seconds, m - Minutes, H - Hours, M - Months, Y - Years. For example: 1D &#x3D; 1 Day | [optional] [default to null]
**StartAt** | **string** | Replication start date and time. | [optional] [default to null]
**KeepLocal** | **string** | For how long to keep the local copy. The format is &lt;integer&gt;&lt;time period&gt;. The time period can be D - Days, W - Weeks, s - Seconds, m - Minutes, H - Hours, M - Months, Y - Years. For example: 1D &#x3D; 1 Day | [optional] [default to null]
**KeepRemote** | **string** | For how long to keep the copy on the remote peer. The format is &lt;integer&gt;&lt;time period&gt;. The time period can be D - Days, W - Weeks, s - Seconds, m - Minutes, H - Hours, M - Months, Y - Years. For example: 1D &#x3D; 1 Day | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

