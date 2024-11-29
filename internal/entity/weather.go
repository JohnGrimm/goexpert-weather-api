package entity

import (
	"fmt"
	"strconv"
)

type Weather struct {
	Celcius    float64
	Fahrenheit float64
	Kelvin     float64
}

type WeatherDetails struct {
	Temp float64
}

type WeatherResponse struct {
	Main WeatherDetails
}

func NewWeather(celcius, fahrenheit, kelvin float64) *Weather {
	return &Weather{
		Celcius:    celcius,
		Fahrenheit: fahrenheit,
		Kelvin:     kelvin,
	}
}

func (w *Weather) MakeTemperatureConversions(weather_res_main_temp float64) {
	w.Celcius, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", weather_res_main_temp), 64)
	w.Fahrenheit, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", weather_res_main_temp*1.8+32), 64)
	w.Kelvin, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", weather_res_main_temp+273.15), 64)
}
