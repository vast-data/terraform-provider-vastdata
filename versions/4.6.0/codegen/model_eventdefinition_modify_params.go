/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type EventdefinitionModifyParams struct {
	// The severity of an alarm triggered by this event. INFO means no alarm is triggered.
	Severity string `json:"severity,omitempty"`
	// For 'Object Modified' alarms, the value of the monitored property at which to trigger an alarm.
	TriggerOn []ErrorUnknown `json:"trigger_on,omitempty"`
	// For 'Object Modified' alarms, the value of the monitored property at which to trigger off an alarm. 
	TriggerOff []ErrorUnknown `json:"trigger_off,omitempty"`
	// For rate alarms, the The time frame over which to monitor the property.
	TimeFrame string `json:"time_frame,omitempty"`
	// Comma separated list of email recipients for alarms
	EmailRecipients []string `json:"email_recipients,omitempty"`
	// The URL of the API endpoint of an external application to trigger on alarms.
	WebhookUrl string `json:"webhook_url,omitempty"`
	// The HTTP method to invoke with the webhook trigger.
	WebhookMethod string `json:"webhook_method,omitempty"`
	// The payload, if required, to send with a POST command. You can use the $event variable to include the event message.
	WebhookData string `json:"webhook_data,omitempty"`
	// The URL parameters to send with a GET command.
	WebhookParams string `json:"webhook_params,omitempty"`
	// Set to true to disable alert actions for the event definition.
	DisableActions bool `json:"disable_actions,omitempty"`
	// Set to true to enable events, alarms and actions.
	Enabled bool `json:"enabled,omitempty"`
	// 
	Internal bool `json:"internal,omitempty"`
	// Set to true for only alarms to trigger configured actions such as email and webhook. Set to false for all events of the definition to trigger the configured actions.
	AlarmOnly bool `json:"alarm_only,omitempty"`
}
