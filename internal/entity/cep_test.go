package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCEP(t *testing.T) {
	cep := NewCEP(
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

	assert.Equal(t, cep.CEP, "01001001")
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
