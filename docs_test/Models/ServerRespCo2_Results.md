# ServerRespCo2_Results
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Year** | **Integer** | The year this measurement was taken. | [optional] [default to null]
**Month** | **Integer** | The month this measurement was taken. | [optional] [default to null]
**Day** | **Integer** | The day representing the start of the week for this measurement. Measurements are taken hourly and averaged together over a week. | [optional] [default to null]
**DateDecimal** | **Float** | A decimal representation of the week this measurement was taken. | [optional] [default to null]
**Average** | **Integer** | The average gas measurement recorded for the week. | [optional] [default to null]
**NumDays** | **Integer** | The number of days measurements were taken to compute the weekly average. | [optional] [default to null]
**OneYearAgo** | **Float** | The CO2 mole fraction in dry air (in parts-per-million) exactly 365 days prior to this measurement. | [optional] [default to null]
**TenYearsAgo** | **Float** | The CO2 mole fraction in dry air (in parts-per-million) exactly 10*365 days + 3 days (for leap years) prior to this measurement. | [optional] [default to null]
**IncSincePreIndustrial** | **Float** | The CO2 mole fraction difference in dry air (in parts-per-million) between this measurement and measurements from 1800. | [optional] [default to null]
**Timestamp** | **Integer** | The unix timestamp of when the measurement was recorded. This is always recorded as a date with the hh:mm:ss portion of the timestamp set to zero. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

