# EventdefinitionModifyParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Severity** | **string** | The severity of an alarm triggered by this event. INFO means no alarm is triggered. | [optional] [default to null]
**TriggerOn** | [**[]ErrorUnknown**](.md) | For &#39;Object Modified&#39; alarms, the value of the monitored property at which to trigger an alarm. | [optional] [default to null]
**TriggerOff** | [**[]ErrorUnknown**](.md) | For &#39;Object Modified&#39; alarms, the value of the monitored property at which to trigger off an alarm.  | [optional] [default to null]
**TimeFrame** | **string** | For rate alarms, the The time frame over which to monitor the property. | [optional] [default to null]
**EmailRecipients** | **[]string** | Comma separated list of email recipients for alarms | [optional] [default to null]
**WebhookUrl** | **string** | The URL of the API endpoint of an external application to trigger on alarms. | [optional] [default to null]
**WebhookMethod** | **string** | The HTTP method to invoke with the webhook trigger. | [optional] [default to null]
**WebhookData** | **string** | The payload, if required, to send with a POST command. You can use the $event variable to include the event message. | [optional] [default to null]
**WebhookParams** | **string** | The URL parameters to send with a GET command. | [optional] [default to null]
**DisableActions** | **bool** | Set to true to disable alert actions for the event definition. | [optional] [default to null]
**Enabled** | **bool** | Set to true to enable events, alarms and actions. | [optional] [default to null]
**Internal** | **bool** |  | [optional] [default to null]
**AlarmOnly** | **bool** | Set to true for only alarms to trigger configured actions such as email and webhook. Set to false for all events of the definition to trigger the configured actions. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


