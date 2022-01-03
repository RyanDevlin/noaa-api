# Ch4MonthlyApi

All URIs are relative to *https://api.planetpulse.io/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getCh4**](Ch4MonthlyApi.md#getCh4) | **GET** /ch4 | Requests monthly CH4 measurements.
[**getCh4Monthly**](Ch4MonthlyApi.md#getCh4Monthly) | **GET** /ch4/monthly | Requests monthly CH4 measurements.


<a name="getCh4"></a>
# **getCh4**
> oneOf&lt;ServerRespCh4,ServerRespCh4Simple&gt; getCh4(year, month, gt, lt, gte, lte, simple, pretty, limit, offset, page)

Requests monthly CH4 measurements.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **year** | **Integer**| Return all CH4 measurements for a given year. | [optional] [default to null]
 **month** | **Integer**| Return all CH4 measurements for a given month. | [optional] [default to null]
 **gt** | **Float**| Return all CH4 measurements with an average ppb reading greater than the supplied value. | [optional] [default to null]
 **lt** | **Float**| Return all CH4 measurements with an average ppb reading less than the supplied value. | [optional] [default to null]
 **gte** | **Float**| Return all CH4 measurements with an average ppb reading greater than OR equal to the supplied value. | [optional] [default to null]
 **lte** | **Float**| Return all CH4 measurements with an average ppb reading less than OR equal to the supplied value. | [optional] [default to null]
 **simple** | **Boolean**| If true, a smaller, simplified version of each CH4 measurement will be returned. | [optional] [default to false]
 **pretty** | **Boolean**| If true, json responses are indented for readability. | [optional] [default to true]
 **limit** | **Integer**| Maximum number of items to return. | [optional] [default to 10]
 **offset** | **Integer**| Number of items to skip before returning the results. | [optional] [default to 0]
 **page** | **Integer**| Shifts the response data by offset + (limit * page). | [optional] [default to 1]

### Return type

[**oneOf&lt;ServerRespCh4,ServerRespCh4Simple&gt;**](../Models/oneOf&lt;ServerRespCh4,ServerRespCh4Simple&gt;.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getCh4Monthly"></a>
# **getCh4Monthly**
> oneOf&lt;ServerRespCh4,ServerRespCh4Simple&gt; getCh4Monthly(year, month, gt, lt, gte, lte, simple, pretty, limit, offset, page)

Requests monthly CH4 measurements.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **year** | **Integer**| Return all CH4 measurements for a given year. | [optional] [default to null]
 **month** | **Integer**| Return all CH4 measurements for a given month. | [optional] [default to null]
 **gt** | **Float**| Return all CH4 measurements with an average ppb reading greater than the supplied value. | [optional] [default to null]
 **lt** | **Float**| Return all CH4 measurements with an average ppb reading less than the supplied value. | [optional] [default to null]
 **gte** | **Float**| Return all CH4 measurements with an average ppb reading greater than OR equal to the supplied value. | [optional] [default to null]
 **lte** | **Float**| Return all CH4 measurements with an average ppb reading less than OR equal to the supplied value. | [optional] [default to null]
 **simple** | **Boolean**| If true, a smaller, simplified version of each CH4 measurement will be returned. | [optional] [default to false]
 **pretty** | **Boolean**| If true, json responses are indented for readability. | [optional] [default to true]
 **limit** | **Integer**| Maximum number of items to return. | [optional] [default to 10]
 **offset** | **Integer**| Number of items to skip before returning the results. | [optional] [default to 0]
 **page** | **Integer**| Shifts the response data by offset + (limit * page). | [optional] [default to 1]

### Return type

[**oneOf&lt;ServerRespCh4,ServerRespCh4Simple&gt;**](../Models/oneOf&lt;ServerRespCh4,ServerRespCh4Simple&gt;.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

