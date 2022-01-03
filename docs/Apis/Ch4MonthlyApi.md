# Ch4MonthlyApi

All URIs are relative to *https://api.planetpulse.io/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ch4Get**](Ch4MonthlyApi.md#ch4Get) | **GET** /ch4 | Requests monthly CH4 measurements.
[**ch4MonthlyGet**](Ch4MonthlyApi.md#ch4MonthlyGet) | **GET** /ch4/monthly | Requests monthly CH4 measurements.


<a name="ch4Get"></a>
# **ch4Get**
> oneOf&lt;array,array&gt; ch4Get(year, month, gt, lt, gte, lte, simple, pretty, limit, offset, page)

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

[**oneOf&lt;array,array&gt;**](../Models/oneOf&lt;array,array&gt;.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="ch4MonthlyGet"></a>
# **ch4MonthlyGet**
> oneOf&lt;array,array&gt; ch4MonthlyGet(year, month, gt, lt, gte, lte, simple, pretty, limit, offset, page)

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

[**oneOf&lt;array,array&gt;**](../Models/oneOf&lt;array,array&gt;.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

