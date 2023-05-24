# Host

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Guid** | **string** |  | [optional] [default to null]
**VastInstallInfo** | [***interface{}**](interface{}.md) |  | [optional] [default to null]
**Ip** | **string** | internal IP Bond Address | [optional] [default to null]
**Ip1** | **string** | 1st internal IP Address | [optional] [default to null]
**Ip2** | **string** | 2nd internal IP Address | [optional] [default to null]
**Ipv6** | **string** | External IPv6 Address | [optional] [default to null]
**HostLabel** | **string** | Host label, used to label container, e.g. 11.0.0.1-4000 | [optional] [default to null]
**Auto** | **bool** | auto or manually host | [optional] [default to null]
**Hostname** | **string** | Hostname | [default to null]
**NodeType** | **string** | Node type: C-Node/D-Node | [default to null]
**DboxUid** | **string** | Unique h/w identifier | [optional] [default to null]
**Loopback** | **bool** | Loopback (single node) installation | [optional] [default to null]
**PerfCheck** | **bool** | Check cluster performance | [optional] [default to null]
**Dmsetup** | **bool** | Mock NVMeoF devices with dmsetup devices | [optional] [default to null]
**EnableDr** | **bool** | enable data reduction | [optional] [default to null]
**HalfSystem** | **bool** | Is half system | [optional] [default to null]
**DeepStripe** | **bool** | Is deep stripe system | [optional] [default to null]
**DrHashSize** | **int32** | DR hash size in buckets | [optional] [default to null]
**NvramSize** | **int32** | NVRAM size for mocked devices | [optional] [default to null]
**DriveSize** | **int32** | Drive size for mocked devices | [optional] [default to null]
**SshUser** | **string** | SSH User name | [optional] [default to null]
**State** | **string** | Node state | [optional] [default to null]
**PlatformRdmaPort** | **int32** |  | [optional] [default to null]
**HwInfo** | [***interface{}**](interface{}.md) | General hardware info related to the host | [optional] [default to null]
**PlatformTcpPort** | **int32** |  | [optional] [default to null]
**DataRdmaPort** | **int32** |  | [optional] [default to null]
**DataTcpPort** | **int32** |  | [optional] [default to null]
**BoxRdmaPort** | **int32** |  | [optional] [default to null]
**InstallState** | **string** | Node installation state | [optional] [default to null]
**SwVersion** | **string** | Node s/w version | [optional] [default to null]
**OsVersion** | **string** | Node OS version | [optional] [default to null]
**Build** | **string** | Node build | [optional] [default to null]
**Cluster** | **string** | Parent Cluster | [optional] [default to null]
**Url** | **string** |  | [optional] [default to null]
**UpgradeState** | **string** | Host upgrade state | [optional] [default to null]
**NodeGuid** | **string** | Hosts sibling node guid | [optional] [default to null]
**SingleNic** | **string** | Host has single NIC | [optional] [default to null]
**NetType** | **string** |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


