package usecase_test

import (
	"errors"
	"testing"

	"github.com/JohnGrimm/goexpert-weather-api/internal/entity"
	"github.com/JohnGrimm/goexpert-weather-api/internal/infra/web"
	"github.com/JohnGrimm/goexpert-weather-api/internal/usecase"
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

func TestGetCEPUseCase(t *testing.T) {
	t.Run("valid cep", func(t *testing.T) {
		cep_address := "01001001"

		mockCEPRepo := new(MockCEPRepository)
		mockWeatherRepo := new(MockWeatherRepository)

		mockCEPRepo.On("IsValid", cep_address).Return(true)

		mockCEPRepo.On("Get", cep_address).Return([]byte(`{
	        "cep": "01001-001",
	        "logradouro": "Praça da Sé",
	        "complemento": "lado par",
	        "bairro": "Sé",
	        "localidade": "São Paulo",
	        "uf": "SP",
	        "ibge": "3550308",
	        "gia": "1004",
	        "ddd": "11",
	        "siafi": "7107"
	    }`), nil)

		mockCEPRepo.On("Convert", mock.Anything).Return(&entity.CEP{
			CEP:         "01001-001",
			Logradouro:  "Praça da Sé",
			Complemento: "lado par",
			Bairro:      "Sé",
			Localidade:  "São Paulo",
			UF:          "SP",
			IBGE:        "3550308",
			GIA:         "1004",
			DDD:         "11",
			SIAFI:       "7107",
		}, nil)

		apiKey := "test_api_key"
		webCEPHandler := web.NewWebCEPHandlerWithDeps(mockCEPRepo, mockWeatherRepo, apiKey)

		get_cep_dto := usecase.CEPInputDTO{
			CEP: cep_address,
		}
		getCEP := usecase.NewGetCEPUseCase(webCEPHandler.CEPRepository)
		cep_output, err := getCEP.Execute(get_cep_dto)
		assert.NoError(t, err)
		assert.Equal(t, cep_output.CEP, "01001-001")
		assert.Equal(t, cep_output.Logradouro, "Praça da Sé")
		assert.Equal(t, cep_output.Complemento, "lado par")
		assert.Equal(t, cep_output.Bairro, "Sé")
		assert.Equal(t, cep_output.Localidade, "São Paulo")
		assert.Equal(t, cep_output.UF, "SP")
		assert.Equal(t, cep_output.IBGE, "3550308")
		assert.Equal(t, cep_output.GIA, "1004")
		assert.Equal(t, cep_output.DDD, "11")
		assert.Equal(t, cep_output.SIAFI, "7107")
	})

	t.Run("invalid zipcode", func(t *testing.T) {
		cep_address := "0100100"

		mockCEPRepo := new(MockCEPRepository)
		mockWeatherRepo := new(MockWeatherRepository)

		mockCEPRepo.On("IsValid", cep_address).Return(false)

		mockCEPRepo.On("Get", cep_address).Return([]byte(`{}`), errors.New("invalid zipcode"))

		mockCEPRepo.On("Convert", mock.Anything).Return(&entity.CEP{
			CEP:         "",
			Logradouro:  "",
			Complemento: "",
			Bairro:      "",
			Localidade:  "",
			UF:          "",
			IBGE:        "",
			GIA:         "",
			DDD:         "",
			SIAFI:       "",
		}, nil)

		apiKey := "test_api_key"
		webCEPHandler := web.NewWebCEPHandlerWithDeps(mockCEPRepo, mockWeatherRepo, apiKey)

		get_cep_dto := usecase.CEPInputDTO{
			CEP: cep_address,
		}
		getCEP := usecase.NewGetCEPUseCase(webCEPHandler.CEPRepository)
		cep_output, err := getCEP.Execute(get_cep_dto)
		assert.Error(t, err)
		assert.Equal(t, cep_output.CEP, "")
		assert.Equal(t, cep_output.Logradouro, "")
		assert.Equal(t, cep_output.Complemento, "")
		assert.Equal(t, cep_output.Bairro, "")
		assert.Equal(t, cep_output.Localidade, "")
		assert.Equal(t, cep_output.UF, "")
		assert.Equal(t, cep_output.IBGE, "")
		assert.Equal(t, cep_output.GIA, "")
		assert.Equal(t, cep_output.DDD, "")
		assert.Equal(t, cep_output.SIAFI, "")
	})
}
