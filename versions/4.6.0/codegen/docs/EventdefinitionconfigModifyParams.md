# EventdefinitionconfigModifyParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SmtpHost** | **string** | SMTP server host name for alert emails. | [optional] [default to null]
**SmtpPort** | **string** | The port used by the SMTP server to send outgoing emails.  | [optional] [default to null]
**SmtpUser** | **string** | User for SMTP authentication | [optional] [default to null]
**SmtpPassword** | **string** | Password for SMTP authentication | [optional] [default to null]
**SmtpUseTls** | **bool** | Set to true to send email over a TLS connection. | [optional] [default to null]
**EmailSubject** | **string** | Optional and global email subject for all alarm notification emails. Leave blank to send alarm info in the subject. | [optional] [default to null]
**EmailSender** | **string** | Global for all alarm notification emails, the sender email that appears in the emails. | [optional] [default to null]
**EmailRecipients** | **[]string** | Default email recipients. These recipients receive notifications of all alarms except those triggered by events that have a different list of email recipients specified in the event definition or for which actions are disabled. | [optional] [default to null]
**WebhookUrl** | **string** | The URL of the API endpoint of an external application, including parameters. | [optional] [default to null]
**WebhookData** | **string** | The payload, if required, for the endpoint. You can use the $event variable to include the event message. | [optional] [default to null]
**WebhookMethod** | **string** | The HTTP method to invoke with the webhook trigger. | [optional] [default to null]
**SyslogHost** | **string** | The syslog server&#39;s IP address, for sending event logs to a syslog server. | [optional] [default to null]
**SyslogPort** | **string** | The port number used by the syslog server to listen on for syslog requests. | [optional] [default to null]
**SyslogProtocol** | **string** | The protocol used for communicating with the remote syslog server. | [optional] [default to null]
**SyslogVmsAudit** | **bool** | VMS audit | [optional] [default to null]
**SyslogShellAudit** | **bool** | CNode and DNode shell commands | [optional] [default to null]
**SyslogIpmiAudit** | **bool** | CNode and DNode IPMI commands | [optional] [default to null]
**AuditLogsRetention** | **int32** | Audit logs retention in days | [optional] [default to null]
**QuotaEmailSuffix** | **string** | A default suffix to add to a username in case there is no specific email address | [optional] [default to null]
**QuotaEmailProvider** | **string** |  | [optional] [default to null]
**QuotaEmailInterval** | **string** |  | [optional] [default to null]
**QuotaEmailHourlyLimit** | **int32** | Maximum quota alert emails VMS will send per hour | [optional] [default to null]
**DisableActions** | **bool** | Set to true to disable default actions for events. | [optional] [default to null]
**Enabled** | **bool** |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


