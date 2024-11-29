package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JohnGrimm/goexpert-weather-api/internal/entity"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCEPRepository struct {
	mock.Mock
}

func (m *MockCEPRepository) Get(cep string) ([]byte, error) {
	args := m.Called(cep)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCEPRepository) Convert(data []byte) (*entity.CEP, error) {
	args := m.Called(data)
	return args.Get(0).(*entity.CEP), args.Error(1)
}

func (m *MockCEPRepository) IsValid(cep string) bool {
	args := m.Called(cep)
	return args.Bool(0)
}

type MockWeatherRepository struct {
	mock.Mock
}

func (m *MockWeatherRepository) Get(localidade, apiKey string) ([]byte, error) {
	args := m.Called(localidade, apiKey)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockWeatherRepository) ConvertToWeather(data *entity.WeatherResponse) (*entity.Weather, error) {
	args := m.Called(data)
	return args.Get(0).(*entity.Weather), args.Error(1)
}

func (m *MockWeatherRepository) ConvertToWeatherResponse(data []byte) (*entity.WeatherResponse, error) {
	args := m.Called(data)
	return args.Get(0).(*entity.WeatherResponse), args.Error(1)
}

func TestCEPHandler(t *testing.T) {
	mockCEPRepo := new(MockCEPRepository)
	mockWeatherRepo := new(MockWeatherRepository)

	apiKey := "test_api_key"
	handler := NewWebCEPHandlerWithDeps(mockCEPRepo, mockWeatherRepo, apiKey)

	router := chi.NewRouter()
	router.Get("/cep/{cep}", handler.Get)

	cep := "01001001"
	localidade := "São Paulo"
	weatherResponse := &entity.WeatherResponse{
		Main: entity.WeatherDetails{Temp: 298.15},
	}

	weather := &entity.Weather{
		Celcius:    25.0,
		Fahrenheit: 77.0,
		Kelvin:     298.15,
	}

	cepData := []byte(`{"localidade": "São Paulo"}`)
	mockCEPRepo.On("Get", cep).Return(cepData, nil)
	mockCEPRepo.On("Convert", cepData).Return(&entity.CEP{Localidade: localidade}, nil)
	mockCEPRepo.On("IsValid", cep).Return(true)

	mockWeatherRepo.On("Get", localidade, apiKey).Return([]byte(`{"main":{"temp":298.15}}`), nil)
	mockWeatherRepo.On("ConvertToWeatherResponse", []byte(`{"main":{"temp":298.15}}`)).Return(weatherResponse, nil)
	mockWeatherRepo.On("ConvertToWeather", weatherResponse).Return(weather, nil)

	req, err := http.NewRequest("GET", "/cep/"+cep, nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var responseWeather *entity.Weather
	err = json.Unmarshal(rr.Body.Bytes(), &responseWeather)
	assert.NoError(t, err)

	assert.Equal(t, weather.Celcius, responseWeather.Celcius)
	assert.Equal(t, weather.Fahrenheit, responseWeather.Fahrenheit)
	assert.Equal(t, weather.Kelvin, responseWeather.Kelvin)

	mockCEPRepo.AssertExpectations(t)
	mockWeatherRepo.AssertExpectations(t)
}
