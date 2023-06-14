# EventDefinition

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**ObjectType** | **string** |  | [optional] [default to null]
**EventType** | **string** |  | [optional] [default to null]
**Severity** | **string** | The severity of the alarm | [optional] [default to null]
**TriggerOn** | **string** | Trigger to turn on the alarm | [optional] [default to null]
**TriggerOff** | **string** | Trigger to turn off the alarm | [optional] [default to null]
**UserModified** | **bool** | Did a user modify this event definition | [optional] [default to null]
**TimeFrame** | **string** | Time frame for rate alarms | [optional] [default to null]
**EmailRecipients** | **string** | List of emails you want to notify in case this alarm occurs (separated by comma) | [optional] [default to null]
**WebhookUrl** | **string** | The URL that the webhook will go to | [optional] [default to null]
**WebhookData** | **string** | Use $event as event message parameter | [optional] [default to null]
**WebhookMethod** | **string** | The method that the webhook will use | [optional] [default to null]
**Property** | **string** |  | [optional] [default to null]
**WebhookParams** | **string** |  | [optional] [default to null]
**AlarmDefinitions** | **string** |  | [optional] [default to null]
**ActionDefinitions** | **string** |  | [optional] [default to null]
**EventMessage** | **string** |  | [optional] [default to null]
**DisableActions** | **bool** |  | [optional] [default to null]
**Enable** | **bool** |  | [optional] [default to null]
**Internal** | **bool** |  | [optional] [default to null]
**AlarmOnly** | **bool** | When this is enabled, only alarms will lead to email and webhook actions | [optional] [default to null]
**Metadata** | [***interface{}**](interface{}.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


