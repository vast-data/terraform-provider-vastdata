# NativeReplicationRemoteTarget

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Name** | **string** |  | [optional] [default to null]
**PeerName** | **string** | Name of remote peer | [optional] [default to null]
**State** | **string** | State of peer connectivity | [optional] [default to null]
**Guid** | **string** | unique identifier | [optional] [default to null]
**Password** | **string** | password for authentication | [optional] [default to null]
**LeadingVip** | **string** | A VIP belonging to the remote peer&#39;s replication VIP Pool, used for connecting to the remote peer. | [optional] [default to null]
**RemoteVipRange** | **string** | VIP range of the remote peer&#39;s replication VIP Pool | [optional] [default to null]
**PoolId** | **string** | The ID of the VIP pool on the local cluster configured with the replication role | [optional] [default to null]
**Version** | **string** | The VAST software version running on the local peer. | [optional] [default to null]
**RemoteVersion** | **string** | The VAST software version running on the remote peer. | [optional] [default to null]
**LastHeartBeat** | **string** | The time of the last successful message sent, arrived and acknowledged by the peer. | [optional] [default to null]
**SpaceLeft** | **string** | The logical capacity remaining available on the remote peer. | [optional] [default to null]
**PeerCertificate** | **string** | A certificate to use for authentication with the peer. | [optional] [default to null]
**Secret** | **string** | Not yet implemented | [optional] [default to null]
**Mss** | **int32** | Maximum segment size (MSS), in bytes, that the peer can receive in a single TCP segment. | [optional] [default to null]
**Health** | **string** | Reflects health of connection between peers. | [optional] [default to null]
**SecureMode** | **string** | Secure mode | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


