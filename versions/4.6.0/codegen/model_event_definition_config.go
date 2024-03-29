/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type EventDefinitionConfig struct {
	// 
	Id int32 `json:"id,omitempty"`
	// SMTP host for alert emails
	SmtpHost string `json:"smtp_host,omitempty"`
	// Connection port on the SMTP host
	SmtpPort int32 `json:"smtp_port,omitempty"`
	// User for SMTP authentication
	SmtpUser string `json:"smtp_user,omitempty"`
	// Password for SMTP authentication
	SmtpPassword string `json:"smtp_password,omitempty"`
	SmtpUseTls bool `json:"smtp_use_tls,omitempty"`
	EmailSubject string `json:"email_subject,omitempty"`
	EmailSender string `json:"email_sender,omitempty"`
	EmailRecipients string `json:"email_recipients,omitempty"`
	WebhookUrl string `json:"webhook_url,omitempty"`
	WebhookData string `json:"webhook_data,omitempty"`
	WebhookMethod string `json:"webhook_method,omitempty"`
	DisabeActions bool `json:"disabe_actions,omitempty"`
	// Syslog host for events logging. Use commas for multiple hosts
	SyslogHost string `json:"syslog_host,omitempty"`
	// Syslog port for events logging
	SyslogPort int32 `json:"syslog_port,omitempty"`
	// Syslog protocol for events logging. Default is UDP
	SyslogProtocol string `json:"syslog_protocol,omitempty"`
	// Enable VMS audit
	SyslogVmsAudit bool `json:"syslog_vms_audit,omitempty"`
	// Enable login/logout (GUI/CLI/VMS/SSH/IPMI),shell, clush, sudo and docker commands audit for CNode and DNode
	SyslogShellAudit bool `json:"syslog_shell_audit,omitempty"`
	// Enable CNode and DNode IPMI commands audit
	SyslogIpmiAudit bool `json:"syslog_ipmi_audit,omitempty"`
	// Audit logs retention in days
	AuditLogsRetention int32 `json:"audit_logs_retention,omitempty"`
	QuotaEmailSuffix string `json:"quota_email_suffix,omitempty"`
	QuotaEmailProvider string `json:"quota_email_provider,omitempty"`
	// Minimum interval between emails to the same address. D HH:MM:SS
	QuotaEmailInterval string `json:"quota_email_interval,omitempty"`
	// Maximum quota alert emails VMS will send per hour
	QuotaEmailHourlyLimit int32 `json:"quota_email_hourly_limit,omitempty"`
}
