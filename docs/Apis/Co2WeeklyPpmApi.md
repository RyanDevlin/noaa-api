# Co2WeeklyPpmApi

All URIs are relative to *https://api.planetpulse.io/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getCo2PPM**](Co2WeeklyPpmApi.md#getCo2PPM) | **GET** /co2/weekly/{ppm} | Requests a single weekly CO2 measurement by PPM.


<a name="getCo2PPM"></a>
# **getCo2PPM**
> oneOf&lt;ServerRespCo2,ServerRespCo2Simple&gt; getCo2PPM(ppm, simple, pretty)

Requests a single weekly CO2 measurement by PPM.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ppm** | **Float**| The average CO2 measurement to retrieve, in parts-per-million, taken at Mauna Loa Observatory. | [default to 0]
 **simple** | **Boolean**| If true, a smaller, simplified version of each CO2 measurement will be returned. | [optional] [default to false]
 **pretty** | **Boolean**| If true, json responses are indented for readability. | [optional] [default to true]

### Return type

[**oneOf&lt;ServerRespCo2,ServerRespCo2Simple&gt;**](../Models/Co2Resp.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

