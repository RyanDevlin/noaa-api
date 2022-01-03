# Documentation for Planet Pulse

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *https://api.planetpulse.io/v1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*Co2WeeklyApi* | [**getCo2**](Apis/Co2WeeklyApi.md#getco2) | **GET** / | Requests weekly CO2 measurements.
*Co2WeeklyApi* | [**getCo2**](Apis/Co2WeeklyApi.md#getco2) | **GET** /co2 | Requests weekly CO2 measurements.
*Co2WeeklyApi* | [**getCo2Weekly**](Apis/Co2WeeklyApi.md#getco2weekly) | **GET** /co2/weekly | Requests weekly CO2 measurements.
*Co2WeeklyIncreaseApi* | [**getCo2WeeklyIncrease**](Apis/Co2WeeklyIncreaseApi.md#getco2weeklyincrease) | **GET** /co2/weekly/increase | Requests weekly CO2 measurements by increase in ppm since 1800.
*Co2WeeklyPpmApi* | [**getCo2PPM**](Apis/Co2WeeklyPpmApi.md#getco2ppm) | **GET** /co2/weekly/{ppm} | Requests a single weekly CO2 measurement by PPM.
*Ch4MonthlyApi* | [**getCh4**](Apis/Ch4MonthlyApi.md#getch4) | **GET** /ch4 | Requests monthly CH4 measurements.
*Ch4MonthlyApi* | [**getCh4Monthly**](Apis/Ch4MonthlyApi.md#getch4monthly) | **GET** /ch4/monthly | Requests monthly CH4 measurements.
*Ch4MonthlyTrendApi* | [**getCh4MonthlyTrend**](Apis/Ch4MonthlyTrendApi.md#getch4monthlytrend) | **GET** /ch4/monthly/trend | Requests monthly CH4 measurements.
*HeatlhApi* | [**getServerHealth**](Apis/HeatlhApi.md#getserverhealth) | **GET** /health | An endpoint to perform a server health check.


<a name="documentation-for-models"></a>
## Documentation for Models
 - [Co2Resp](./Models/Co2Resp.md)
 - [Ch4Resp](./Models/Ch4Resp.md)
 - [ErrorResp](./Models/ErrorResp.md)
 - [ServerRespCh4](./Models/ServerRespCh4.md)
 - [ServerRespCh4Simple](./Models/ServerRespCh4Simple.md)
 - [ServerRespCh4Simple_Results](./Models/ServerRespCh4Simple_Results.md)
 - [ServerRespCh4_Results](./Models/ServerRespCh4_Results.md)
 - [ServerRespCo2](./Models/ServerRespCo2.md)
 - [ServerRespCo2Simple](./Models/ServerRespCo2Simple.md)
 - [ServerRespCo2Simple_Results](./Models/ServerRespCo2Simple_Results.md)
 - [ServerRespCo2_Results](./Models/ServerRespCo2_Results.md)
 - [ServerRespError](./Models/ServerRespError.md)
 - [ServerRespHealth](./Models/ServerRespHealth.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
