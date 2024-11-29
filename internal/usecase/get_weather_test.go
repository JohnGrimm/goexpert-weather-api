package usecase_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/JohnGrimm/goexpert-weather-api/internal/infra/repo"
	"github.com/JohnGrimm/goexpert-weather-api/internal/usecase"
	"github.com/stretchr/testify/assert"
)

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func getWeatherJSON() []byte {
	return []byte(
		`{
			"coord": {
				"lon": -46.6361,
				"lat": -23.5475
			},
			"weather": [
				{
					"id": 803,
					"main": "Clouds",
					"description": "broken clouds",
					"icon": "04d"
				}
			],
			"base": "stations",
			"main": {
				"temp": 21.1,
				"feels_like": 21.35,
				"temp_min": 19.75,
				"temp_max": 24.14,
				"pressure": 1024,
				"humidity": 80
			},
			"visibility": 10000,
			"wind": {
				"speed": 3.6,
				"deg": 140
			},
			"clouds": {
				"all": 75
			},
			"dt": 1716126286,
			"sys": {
				"type": 1,
				"id": 8394,
				"country": "BR",
				"sunrise": 1716111334,
				"sunset": 1716150647
			},
			"timezone": -10800,
			"id": 3448439,
			"name": "São Paulo",
			"cod": 200
		}`,
	)
}
func NewMockClient(json_type string) *MockClient {
	var weather_json []byte

	if json_type == "invalid" {
		weather_json = getWeatherJSON()
		return &MockClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 500,
					Body:       io.NopCloser(bytes.NewReader([]byte(``))),
					Header:     make(http.Header),
				}, errors.New("fail to get weather")
			},
		}
	}

	return &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			weather_json = getWeatherJSON()
			mockResponse := weather_json
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(mockResponse))),
				Header:     make(http.Header),
			}, nil
		},
	}
}

func TestGetWeather(t *testing.T) {
	mock_client_with_valid_json := NewMockClient("valid")
	weather_repository_with_valid_json := repo.NewWeatherRepository(mock_client_with_valid_json)
	get_weather_with_valid_json := usecase.NewGetWeatherUseCase(weather_repository_with_valid_json)

	mock_client_with_invalid_json := NewMockClient("invalid")
	weather_repository_with_invalid_json := repo.NewWeatherRepository(mock_client_with_invalid_json)
	get_weather_with_invalid_json := usecase.NewGetWeatherUseCase(weather_repository_with_invalid_json)

	t.Run("valid weather", func(t *testing.T) {
		get_weather_dto := usecase.WeatherInputDTO{
			Localidade: "São Paulo",
			ApiKey:     "test_api_key",
		}

		weather_output, err := get_weather_with_valid_json.Execute(get_weather_dto)
		assert.NoError(t, err)
		assert.IsType(t, &weather_output, &usecase.WeatherOutputDTO{})
	})

	t.Run("missing Localidade field", func(t *testing.T) {
		get_weather_dto := usecase.WeatherInputDTO{
			ApiKey: "test_api_key",
		}

		weather_output, err := get_weather_with_valid_json.Execute(get_weather_dto)
		assert.EqualError(t, err, "missing input: Localidade")
		assert.IsType(t, &weather_output, &usecase.WeatherOutputDTO{})
	})

	t.Run("missing ApiKey field", func(t *testing.T) {
		get_weather_dto := usecase.WeatherInputDTO{
			Localidade: "São Paulo",
		}

		weather_output, err := get_weather_with_valid_json.Execute(get_weather_dto)
		assert.EqualError(t, err, "missing input: ApiKey")
		assert.IsType(t, &weather_output, &usecase.WeatherOutputDTO{})
	})

	t.Run("fail to get weather", func(t *testing.T) {
		get_weather_dto := usecase.WeatherInputDTO{
			Localidade: "goexpert",
			ApiKey:     "test_api_key",
		}

		weather_output, err := get_weather_with_invalid_json.Execute(get_weather_dto)
		assert.EqualError(t, err, "fail to get weather")
		assert.Equal(t, weather_output.Celcius, float64(0))
		assert.Equal(t, weather_output.Fahrenheit, float64(0))
		assert.Equal(t, weather_output.Kelvin, float64(0))
	})
}
