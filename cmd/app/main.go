package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/JohnGrimm/goexpert-weather-api/internal/infra/repo"
	"github.com/JohnGrimm/goexpert-weather-api/internal/infra/web"
	"github.com/JohnGrimm/goexpert-weather-api/internal/infra/web/webserver"
)

func ConfigureServer() *webserver.WebServer {
	webserver := webserver.NewWebServer(":8080")

	cepRepo := repo.NewCEPRepository()
	weatherRepo := repo.NewWeatherRepository(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})

	open_weathermap_api_key := os.Getenv("OPEN_WEATHERMAP_API_KEY")
	if open_weathermap_api_key == "" {
		log.Fatal("Please, provide the OPEN_WEATHERMAP_API_KEY environment variable; Make sure you provide a valid api-key, otherwise it will not be possible to get and convert weather data")
	}

	webCEPHandler := web.NewWebCEPHandlerWithDeps(cepRepo, weatherRepo, os.Getenv("OPEN_WEATHERMAP_API_KEY"))

	webStatusHandler := web.NewWebStatusHandler()

	// CEP
	webserver.AddHandler("GET /cep/{cep}", webCEPHandler.Get)

	// Health check
	webserver.AddHandler("GET /status", webStatusHandler.Get)

	return webserver
}

func main() {
	webserver := ConfigureServer()
	fmt.Println("Starting web server on port", ":8080")
	webserver.Start()
}
