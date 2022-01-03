# Co2WeeklyApi

All URIs are relative to *https://api.planetpulse.io/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getCo2**](Co2WeeklyApi.md#getCo2) | **GET** /co2 | Requests weekly CO2 measurements.
[**getCo2Weekly**](Co2WeeklyApi.md#getCo2Weekly) | **GET** /co2/weekly | Requests weekly CO2 measurements.


<a name="getCo2"></a>
# **getCo2**
> oneOf&lt;ServerRespCo2,ServerRespCo2Simple&gt; getCo2(year, month, gt, lt, gte, lte, simple, pretty, limit, offset, page)

Requests weekly CO2 measurements.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **year** | **Integer**| Return all CO2 measurements for a given year. | [optional] [default to null]
 **month** | **Integer**| Return all CO2 measurements for a given month. | [optional] [default to null]
 **gt** | **Float**| Return all CO2 measurements with a ppm reading greater than the supplied value. | [optional] [default to null]
 **lt** | **Float**| Return all CO2 measurements with a ppm reading less than the supplied value. | [optional] [default to null]
 **gte** | **Float**| Return all CO2 measurements with a ppm reading greater than OR equal to the supplied value. | [optional] [default to null]
 **lte** | **Float**| Return all CO2 measurements with a ppm reading less than OR equal to the supplied value. | [optional] [default to null]
 **simple** | **Boolean**| If true, a smaller, simplified version of each CO2 measurement will be returned. | [optional] [default to false]
 **pretty** | **Boolean**| If true, json responses are indented for readability. | [optional] [default to true]
 **limit** | **Integer**| Maximum number of items to return. | [optional] [default to 10]
 **offset** | **Integer**| Number of items to skip before returning the results. | [optional] [default to 0]
 **page** | **Integer**| Shifts the response data by offset + (limit * page). | [optional] [default to 1]

### Return type

[**oneOf&lt;ServerRespCo2,ServerRespCo2Simple&gt;**](../Models/oneOf&lt;ServerRespCo2,ServerRespCo2Simple&gt;.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getCo2Weekly"></a>
# **getCo2Weekly**
> oneOf&lt;ServerRespCo2,ServerRespCo2Simple&gt; getCo2Weekly(year, month, gt, lt, gte, lte, simple, pretty, limit, offset, page)

Requests weekly CO2 measurements.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **year** | **Integer**| Return all CO2 measurements for a given year. | [optional] [default to null]
 **month** | **Integer**| Return all CO2 measurements for a given month. | [optional] [default to null]
 **gt** | **Float**| Return all CO2 measurements with a ppm reading greater than the supplied value. | [optional] [default to null]
 **lt** | **Float**| Return all CO2 measurements with a ppm reading less than the supplied value. | [optional] [default to null]
 **gte** | **Float**| Return all CO2 measurements with a ppm reading greater than OR equal to the supplied value. | [optional] [default to null]
 **lte** | **Float**| Return all CO2 measurements with a ppm reading less than OR equal to the supplied value. | [optional] [default to null]
 **simple** | **Boolean**| If true, a smaller, simplified version of each CO2 measurement will be returned. | [optional] [default to false]
 **pretty** | **Boolean**| If true, json responses are indented for readability. | [optional] [default to true]
 **limit** | **Integer**| Maximum number of items to return. | [optional] [default to 10]
 **offset** | **Integer**| Number of items to skip before returning the results. | [optional] [default to 0]
 **page** | **Integer**| Shifts the response data by offset + (limit * page). | [optional] [default to 1]

### Return type

[**oneOf&lt;ServerRespCo2,ServerRespCo2Simple&gt;**](../Models/Co2Responses.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

