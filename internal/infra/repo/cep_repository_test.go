package repo

import (
	"testing"

	"github.com/JohnGrimm/goexpert-weather-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestValidateCep(t *testing.T) {
	t.Run("valid_cep", func(t *testing.T) {
		cep := entity.NewCEP(
			"01001001",
			"Praça da Sé",
			"lado par",
			"Sé",
			"São Paulo",
			"SP",
			"3550308",
			"1004",
			"11",
			"7107",
		)
		repo := NewCEPRepository()
		validation := repo.IsValid(cep.CEP)
		assert.True(t, validation)
	})

	t.Run("invalid_cep", func(t *testing.T) {
		cep := entity.NewCEP(
			"01001-001",
			"Praça da Sé",
			"lado par",
			"Sé",
			"São Paulo",
			"SP",
			"3550308",
			"1004",
			"11",
			"7107",
		)
		repo := NewCEPRepository()
		validation := repo.IsValid(cep.CEP)
		assert.False(t, validation)
	})
}

func TestGetCEP(t *testing.T) {
	cep := entity.NewCEP(
		"01001-001",
		"Praça da Sé",
		"lado par",
		"Sé",
		"São Paulo",
		"SP",
		"3550308",
		"1004",
		"11",
		"7107",
	)
	repo := NewCEPRepository()
	cep_res, err := repo.Get(cep.CEP)
	assert.NoError(t, err)
	assert.Contains(t, string(cep_res), "01001-001")
	assert.Contains(t, string(cep_res), "Sé")
	assert.Contains(t, string(cep_res), "São Paulo")
}

func TestConvertCEP(t *testing.T) {
	resp_json := []byte(`{
		"cep":"01001-001",
		"logradouro":"Praça da Sé",
		"complemento":"lado par",
		"bairro":"Sé",
		"localidade":"São Paulo",
		"uf":"SP",
		"ibge":"3550308",
		"gia":"1004",
		"ddd":"11",
		"siafi":"7107"
	}`)

	repo := NewCEPRepository()
	cep, err := repo.Convert(resp_json)
	assert.NoError(t, err)
	assert.Equal(t, cep.CEP, "01001-001")
	assert.Equal(t, cep.Logradouro, "Praça da Sé")
	assert.Equal(t, cep.Complemento, "lado par")
	assert.Equal(t, cep.Bairro, "Sé")
	assert.Equal(t, cep.Localidade, "São Paulo")
	assert.Equal(t, cep.UF, "SP")
	assert.Equal(t, cep.IBGE, "3550308")
	assert.Equal(t, cep.GIA, "1004")
	assert.Equal(t, cep.DDD, "11")
	assert.Equal(t, cep.SIAFI, "7107")
}
