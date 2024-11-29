package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeTemperatureConversions(t *testing.T) {
	weather := NewWeather(1, 1, 1)
	weather.MakeTemperatureConversions(27.1)
	assert.Equal(t, weather.Celcius, 27.1)
	assert.Equal(t, weather.Fahrenheit, weather.Celcius*1.8+32)
	assert.Equal(t, weather.Kelvin, weather.Celcius+273.15)
}
