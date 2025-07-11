# ReplicationPeers

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int64** | The unique ID of the replication peer configuration. | [optional] [default to null]
**Guid** | **string** | The unique GUID of the replication peer configuration. | [optional] [default to null]
**Name** | **string** | The name of the replication peer configuration. | [optional] [default to null]
**Url** | **string** | Direct URL of the replication peer configuration. | [optional] [default to null]
**LeadingVip** | **string** | The virtual IP pool provided for the replication peer configuration. | [optional] [default to null]
**RemoteVipRange** | **string** | The range of virtual IPs that were reported by the peer. | [optional] [default to null]
**Version** | **string** | The version of the source. | [optional] [default to null]
**RemoteVersion** | **string** | The version of the remote peer. | [optional] [default to null]
**IsLocal** | **bool** | Specifies whether the source of the replication is local (this host is the source). | [optional] [default to null]
**PeerName** | **string** | The name of the peer cluster. | [optional] [default to null]
**SecureMode** | **string** | If true, the connection is secure. | [optional] [default to null]
**PoolId** | **int** | The ID of the replication virtual IP pool. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

