# Co2WeeklyApi

All URIs are relative to *https://api.planetpulse.io/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**co2Get**](Co2WeeklyApi.md#co2Get) | **GET** /co2 | Requests weekly CO2 measurements.
[**co2WeeklyGet**](Co2WeeklyApi.md#co2WeeklyGet) | **GET** /co2/weekly | Requests weekly CO2 measurements.


<a name="co2Get"></a>
# **co2Get**
> oneOf&lt;array,array&gt; co2Get(year, month, gt, lt, gte, lte, simple, pretty, limit, offset, page)

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

[**oneOf&lt;array,array&gt;**](../Models/oneOf&lt;array,array&gt;.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="co2WeeklyGet"></a>
# **co2WeeklyGet**
> oneOf&lt;array,array&gt; co2WeeklyGet(year, month, gt, lt, gte, lte, simple, pretty, limit, offset, page)

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

[**oneOf&lt;array,array&gt;**](../Models/oneOf&lt;array,array&gt;.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

