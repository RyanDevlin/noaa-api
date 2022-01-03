# Co2WeeklyIncreaseApi

All URIs are relative to *https://api.planetpulse.io/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getCo2WeeklyIncrease**](Co2WeeklyIncreaseApi.md#getCo2WeeklyIncrease) | **GET** /co2/weekly/increase | Requests weekly CO2 measurements by increase in ppm since 1800.


<a name="getCo2WeeklyIncrease"></a>
# **getCo2WeeklyIncrease**
> oneOf&lt;ServerRespCo2,ServerRespCo2Simple&gt; getCo2WeeklyIncrease(gt, lt, gte, lte, simple, pretty, limit, offset, page)

Requests weekly CO2 measurements by increase in ppm since 1800.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **gt** | **Float**| Return all CO2 measurements with a ppm increase since 1800 greater than the supplied value. | [optional] [default to null]
 **lt** | **Float**| Return all CO2 measurements with a ppm increase since 1800 less than the supplied value. | [optional] [default to null]
 **gte** | **Float**| Return all CO2 measurements with a ppm reading greater than OR equal to the supplied value. | [optional] [default to null]
 **lte** | **Float**| Return all CO2 measurements with a ppm reading less than OR equal to the supplied value. | [optional] [default to null]
 **simple** | **Boolean**| If true, a smaller, simplified version of each CO2 measurement will be returned. | [optional] [default to false]
 **pretty** | **Boolean**| If true, json responses are indented for readability. | [optional] [default to true]
 **limit** | **Integer**| Maximum number of items to return. | [optional] [default to 10]
 **offset** | **Integer**| Number of items to skip before returning the results. | [optional] [default to 0]
 **page** | **Integer**| Shifts the response data by offset + (limit * page). | [optional] [default to 1]

### Return type

[**oneOf&lt;ServerRespCo2,ServerRespCo2Simple&gt;**](../Models/Co2Resp.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

