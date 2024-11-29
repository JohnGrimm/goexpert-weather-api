package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/JohnGrimm/goexpert-weather-api/internal/entity"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type WeatherRepository struct {
	client HTTPClient
}

func NewWeatherRepository(client HTTPClient) *WeatherRepository {
	return &WeatherRepository{
		client: client,
	}
}

func (w *WeatherRepository) Get(localidade string, api_key string) ([]byte, error) {
	localidade = strings.Replace(localidade, " ", "%20", -1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf(
			"http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&temperature.unit=celsius",
			localidade,
			api_key,
		),
		nil,
	)

	if err != nil {
		log.Printf("Fail to create the request: %v", err)
		return nil, err
	}

	res, err := w.client.Do(req)
	if err != nil {
		log.Printf("Fail to make the request: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	ctx_err := ctx.Err()
	if ctx_err != nil {
		if <-ctx.Done(); true {
			err := ctx.Err()
			log.Printf("Max timeout reached: %v", err)
			return nil, err
		}
	}

	resp_json, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Fail to read the response: %v", err)
		return nil, err
	}

	if strings.Contains(string(resp_json), "Invalid API key") {
		return nil, err
	}

	return resp_json, nil
}

func (w *WeatherRepository) ConvertToWeatherResponse(weather_response []byte) (*entity.WeatherResponse, error) {
	var weather_res entity.WeatherResponse
	err := json.Unmarshal(weather_response, &weather_res)
	if err != nil {
		log.Printf("Fail to decode the response: %v", err)
		return nil, err
	}

	return &weather_res, nil
}

func (w *WeatherRepository) ConvertToWeather(weather_response *entity.WeatherResponse) (*entity.Weather, error) {
	weather := entity.Weather{}
	weather.MakeTemperatureConversions(weather_response.Main.Temp)
	return &weather, nil
}
