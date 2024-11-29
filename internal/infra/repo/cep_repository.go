package repo

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/JohnGrimm/goexpert-weather-api/internal/entity"
)

type CEPRepository struct{}

func NewCEPRepository() *CEPRepository {
	return &CEPRepository{}
}

func (r *CEPRepository) IsValid(cep_address string) bool {
	check, _ := regexp.MatchString("^[0-9]{8}$", cep_address)
	return (len(cep_address) == 8 && cep_address != "" && check)
}

func (r *CEPRepository) Get(cep_address string) ([]byte, error) {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://viacep.com.br/ws/%s/json", cep_address), nil)
	if err != nil {
		log.Printf("Fail to create the request: %v", err)
		return nil, err
	}

	res, err := client.Do(req)
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

	return resp_json, nil
}

func (r *CEPRepository) Convert(cep_response []byte) (*entity.CEP, error) {
	var cep entity.CEP
	err := json.Unmarshal(cep_response, &cep)
	if err != nil {
		log.Printf("Fail to decode the response: %v", err)
		return nil, err
	}
	return &cep, nil
}
