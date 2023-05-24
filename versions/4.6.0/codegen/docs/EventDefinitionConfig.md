# EventDefinitionConfig

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**SmtpHost** | **string** | SMTP host for alert emails | [optional] [default to null]
**SmtpPort** | **int32** | Connection port on the SMTP host | [optional] [default to null]
**SmtpUser** | **string** | User for SMTP authentication | [optional] [default to null]
**SmtpPassword** | **string** | Password for SMTP authentication | [optional] [default to null]
**SmtpUseTls** | **bool** |  | [optional] [default to null]
**EmailSubject** | **string** |  | [optional] [default to null]
**EmailSender** | **string** |  | [optional] [default to null]
**EmailRecipients** | **string** |  | [optional] [default to null]
**WebhookUrl** | **string** |  | [optional] [default to null]
**WebhookData** | **string** |  | [optional] [default to null]
**WebhookMethod** | **string** |  | [optional] [default to null]
**DisabeActions** | **bool** |  | [optional] [default to null]
**SyslogHost** | **string** | Syslog host for events logging. Use commas for multiple hosts | [optional] [default to null]
**SyslogPort** | **int32** | Syslog port for events logging | [optional] [default to null]
**SyslogProtocol** | **string** | Syslog protocol for events logging. Default is UDP | [optional] [default to null]
**SyslogVmsAudit** | **bool** | Enable VMS audit | [optional] [default to null]
**SyslogShellAudit** | **bool** | Enable login/logout (GUI/CLI/VMS/SSH/IPMI),shell, clush, sudo and docker commands audit for CNode and DNode | [optional] [default to null]
**SyslogIpmiAudit** | **bool** | Enable CNode and DNode IPMI commands audit | [optional] [default to null]
**AuditLogsRetention** | **int32** | Audit logs retention in days | [optional] [default to null]
**QuotaEmailSuffix** | **string** |  | [optional] [default to null]
**QuotaEmailProvider** | **string** |  | [optional] [default to null]
**QuotaEmailInterval** | **string** | Minimum interval between emails to the same address. D HH:MM:SS | [optional] [default to null]
**QuotaEmailHourlyLimit** | **int32** | Maximum quota alert emails VMS will send per hour | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


