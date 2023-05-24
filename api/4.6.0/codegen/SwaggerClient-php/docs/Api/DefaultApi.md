# Swagger\Client\DefaultApi

All URIs are relative to *https://example.com/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**viewsGet**](DefaultApi.md#viewsGet) | **GET** /views/ | 




# **viewsGet**
> viewsGet($name$tenant_id)



Return a list of VastData views

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');


$apiInstance = new Swagger\Client\Api\DefaultApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client()
);
$name = array("name_example"); // string | The name of the view
$tenant_id = array(new \Swagger\Client\Model\Int()); // Int | The tenant id related to this view


try {
    $apiInstance->viewsGet($name$tenant_id);
} catch (Exception $e) {
    echo 'Exception when calling DefaultApi->viewsGet: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| The name of the view |
 **tenant_id** | [**Int**](../Model/.md)| The tenant id related to this view | [optional]


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)



