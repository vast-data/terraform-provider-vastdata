# CallhomeConfigModifyParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BundleEnabled** | **bool** | Set to true to enable periodic sending of bundles to the Support server | [optional] [default to null]
**BundleInterval** | **int32** | The frequency for sending bundles to the Support server | [optional] [default to null]
**LogEnabled** | **bool** | Set to true to enable system state data to be logged to the Support server | [optional] [default to null]
**LogInterval** | **int32** | The frequency for sending system state data to the Support server. | [optional] [default to null]
**Customer** | **string** | Company name | [optional] [default to null]
**Site** | **string** | Site name | [optional] [default to null]
**Location** | **string** | Site location | [optional] [default to null]
**ProxyHost** | **string** | Proxy IP/hostname | [optional] [default to null]
**ProxyPort** | **string** | Proxy Port | [optional] [default to null]
**ProxyUsername** | **string** | Proxy username | [optional] [default to null]
**ProxyPassword** | **string** | Proxy password | [optional] [default to null]
**TestMode** | **bool** | Set to true to enable test mode | [optional] [default to null]
**VerifySsl** | **bool** | Set to true to enable SSL verification. Set to false to disable. VAST Cluster recognizes SSL certificates from a large range of widely recognized certificate authorities (CAs). VAST Cluster may not recognize an SSL certificate signed by your own in-house CA. | [optional] [default to null]
**ProxyScheme** | **string** |  | [optional] [default to null]
**SupportChannel** | **bool** | Set to true to enable the VAST Support channel. | [optional] [default to null]
**CloudEnabled** | **bool** | Set to true to enable reporting to VAST Cloud Services | [optional] [default to null]
**CloudApiKey** | **string** | Cloud Services API key | [optional] [default to null]
**CloudApiDomain** | **string** |  Cloud Services API domain name | [optional] [default to null]
**CloudSubdomain** | **string** | Cloud Services subdomain, unique per customer, common to all reporting clusters | [optional] [default to null]
**CustomerId** | **string** | The ID issued to the customer | [optional] [default to null]
**MaxUploadConcurrency** | **int32** | The maximum number of parts of a file to upload simultaneously. | [optional] [default to null]
**Obfuscated** | **bool** | If true, call home data is obfuscated. | [optional] [default to null]
**Aggregated** | **bool** | If true, send aggregated callhome logs, otherwise upload logs from each node | [optional] [default to null]
**UploadViaVms** | **bool** | If true, upload non-aggregated Callhome Bundle via VMS (requires proxy). Otherwise, upload from each node. | [optional] [default to null]
**CompressMethod** | **string** | Compression method for callhome bundles (by default zstd) | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


