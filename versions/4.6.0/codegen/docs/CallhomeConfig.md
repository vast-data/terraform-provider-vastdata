# CallhomeConfig

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Name** | **string** |  | [optional] [default to null]
**Guid** | **string** |  | [optional] [default to null]
**BundleEnabled** | **bool** | Enabled/disabled periodic bundle callhome | [optional] [default to null]
**BundleInterval** | **int32** | Interval for periodic bundle | [optional] [default to null]
**LogEnabled** | **bool** | Enabled/disabled periodic log callhome | [optional] [default to null]
**LogInterval** | **int32** | Interval for periodic log | [optional] [default to null]
**Customer** | **string** | Customer/Company name | [optional] [default to null]
**Site** | **string** | Site name | [optional] [default to null]
**Location** | **string** | Site location | [optional] [default to null]
**ProxyScheme** | **string** | http, https, socks5, socks5h | [optional] [default to null]
**ProxyHost** | **string** | Proxy IP/hostname | [optional] [default to null]
**ProxyPort** | **string** | Proxy Port | [optional] [default to null]
**ProxyUsername** | **string** | Proxy username | [optional] [default to null]
**ProxyPassword** | **string** | Proxy password | [optional] [default to null]
**TestMode** | **bool** | enable/disable test mode | [optional] [default to null]
**VerifySsl** | **bool** | Enable/disable ssl certificate verification | [optional] [default to null]
**SupportChannel** | **bool** | Enable/disable the vast support channel | [optional] [default to null]
**CloudEnabled** | **bool** | Is Cloud reporting modules enabled | [optional] [default to null]
**CloudRegistered** | **bool** | Is the cluster registered with VAST Cloud Services | [optional] [default to null]
**CloudApiKey** | **string** | The API Key for the Cloud API | [optional] [default to null]
**CloudApiDomain** | **string** | The domain name for the Cloud API | [optional] [default to null]
**CloudSubdomain** | **string** | The Customer&#39;s subdomain in Vast Cloud Services | [optional] [default to null]
**CustomerId** | **string** | The id issued to the customer | [optional] [default to null]
**MaxUploadConcurrency** | **int32** | Maximum upload concurrency | [optional] [default to null]
**Obfuscated** | **bool** | If true, callhome data is obfuscated | [optional] [default to null]
**Aggregated** | **bool** | If true, send aggregated callhome logs, otherwise upload logs from each node | [optional] [default to null]
**UploadViaVms** | **bool** | If true, upload non-aggregated Callhome Bundle via VMS (requires proxy). Otherwise, upload from each node. | [optional] [default to null]
**CompressMethod** | **string** | Compression method for callhome bundles (by default zstd) | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


