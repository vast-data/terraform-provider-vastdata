/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ReplicationTargetModifyParams struct {
	// Name
	Name string `json:"name,omitempty"`
	// If configured, replication traffic is routed via proxies. Separate with commas. Format: http://<username>:<password>@<IP>:<port>
	Proxies string `json:"proxies,omitempty"`
	// Access key of a valid key pair for accessing the named S3 bucket
	AccessKey string `json:"access_key,omitempty"`
	// The secret key of a valid key pair for accessing the destination S3 bucket
	SecretKey string `json:"secret_key,omitempty"`
	// The S3 bucket name of an existing S3 bucket that you want to configure as the replication target
	BucketName string `json:"bucket_name,omitempty"`
	// For custom S3 buckets (not AWS), the protocol to use to connect to the bucket. Can be http or https.
	HttpProtocol string `json:"http_protocol,omitempty"`
	// custom bucket url
	CustomBucketUrl string `json:"custom_bucket_url,omitempty"`
	// If the target is an AWS S3 bucket, use this parameter to specify the AWS region of the bucket
	AwsRegion string `json:"aws_region,omitempty"`
	// Not yet implemented
	AwsAccountId string `json:"aws_account_id,omitempty"`
	// Not yet implemented
	AwsRole string `json:"aws_role,omitempty"`
}