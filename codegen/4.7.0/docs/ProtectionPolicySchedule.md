# ProtectionPolicySchedule

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Every** | **string** | How often to make a snapshot, format is &lt;integer&gt;&lt;time period&gt; , while time period can be D - Days ,W - Weeks ,s - Seconds ,m - Minutes, H - Hours, M - Months, Y - Years , Ex 1D &#x3D; 1 Day | [optional] [default to null]
**StartAt** | **string** | The start data of the replication | [optional] [default to null]
**KeepLocal** | **string** | For how long to keep a local copy of the replication, format is &lt;integer&gt;&lt;time period&gt; , while time period can be D - Days ,W - Weeks ,s - Seconds ,m - Minutes, H - Hours, M - Months, Y - Years , Ex 1D &#x3D; 1 Day | [optional] [default to null]
**KeepRemote** | **string** | For how long to keep the copy on the remote side, format is &lt;integer&gt;&lt;time period&gt; , while time period can be D - Days ,W - Weeks ,s - Seconds ,m - Minutes, H - Hours, M - Months, Y - Years , Ex 1D &#x3D; 1 Day | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

