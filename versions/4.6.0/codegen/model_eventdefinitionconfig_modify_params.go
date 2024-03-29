/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type EventdefinitionconfigModifyParams struct {
	// SMTP server host name for alert emails.
	SmtpHost string `json:"smtp_host,omitempty"`
	// The port used by the SMTP server to send outgoing emails. 
	SmtpPort string `json:"smtp_port,omitempty"`
	// User for SMTP authentication
	SmtpUser string `json:"smtp_user,omitempty"`
	// Password for SMTP authentication
	SmtpPassword string `json:"smtp_password,omitempty"`
	// Set to true to send email over a TLS connection.
	SmtpUseTls bool `json:"smtp_use_tls,omitempty"`
	// Optional and global email subject for all alarm notification emails. Leave blank to send alarm info in the subject.
	EmailSubject string `json:"email_subject,omitempty"`
	// Global for all alarm notification emails, the sender email that appears in the emails.
	EmailSender string `json:"email_sender,omitempty"`
	// Default email recipients. These recipients receive notifications of all alarms except those triggered by events that have a different list of email recipients specified in the event definition or for which actions are disabled.
	EmailRecipients []string `json:"email_recipients,omitempty"`
	// The URL of the API endpoint of an external application, including parameters.
	WebhookUrl string `json:"webhook_url,omitempty"`
	// The payload, if required, for the endpoint. You can use the $event variable to include the event message.
	WebhookData string `json:"webhook_data,omitempty"`
	// The HTTP method to invoke with the webhook trigger.
	WebhookMethod string `json:"webhook_method,omitempty"`
	// The syslog server's IP address, for sending event logs to a syslog server.
	SyslogHost string `json:"syslog_host,omitempty"`
	// The port number used by the syslog server to listen on for syslog requests.
	SyslogPort string `json:"syslog_port,omitempty"`
	// The protocol used for communicating with the remote syslog server.
	SyslogProtocol string `json:"syslog_protocol,omitempty"`
	// VMS audit
	SyslogVmsAudit bool `json:"syslog_vms_audit,omitempty"`
	// CNode and DNode shell commands
	SyslogShellAudit bool `json:"syslog_shell_audit,omitempty"`
	// CNode and DNode IPMI commands
	SyslogIpmiAudit bool `json:"syslog_ipmi_audit,omitempty"`
	// Audit logs retention in days
	AuditLogsRetention int32 `json:"audit_logs_retention,omitempty"`
	// A default suffix to add to a username in case there is no specific email address
	QuotaEmailSuffix string `json:"quota_email_suffix,omitempty"`
	QuotaEmailProvider string `json:"quota_email_provider,omitempty"`
	QuotaEmailInterval string `json:"quota_email_interval,omitempty"`
	// Maximum quota alert emails VMS will send per hour
	QuotaEmailHourlyLimit int32 `json:"quota_email_hourly_limit,omitempty"`
	// Set to true to disable default actions for events.
	DisableActions bool `json:"disable_actions,omitempty"`
	Enabled bool `json:"enabled,omitempty"`
}
