# goexpert-weather-api
Projeto do Laboratório "Deploy com Cloud Run" do treinamento GoExpert(FullCycle).



## O desafio
Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin). Esse sistema deverá ser publicado no Google Cloud Run.



## Como rodar o projeto: make
``` shell
# build the container image
make build

# push the container image and deploy to GCP Cloud Run
make deploy

# run locally
make run
```



## Como rodar o projeto: manual
``` shell
## 1. Clone o repo

## 2. Crie o .env
cp .env.example .env

## 3. Coloque sua api-key como valor na variável OPEN_WEATHERMAP_API_KEY no .env

## 4. Baixe compose, se estiver up
docker-compose down

## 5. Remover a imagem antiga, se existir
docker image rm -f gcr.io/full-cycle-go-expert/goexpert-weather-api:v1

## 6. Suba o compose 
docker-compose up -d

## 7. Faça as chamadas
## retorno 200
echo -n "422: "; curl -s "http://localhost:8080/cep/1234567"
echo -n "404: "; curl -s "http://localhost:8080/cep/12345678"
echo -n "200: "; curl -s "http://localhost:8080/cep/12212250"
```



## Funcionalidades da Linguagem Utilizadas
- context
- net/http
- encoding/json
- testing
- testify



## Requisitos: sistema
- [x] O sistema deve receber um CEP válido de 8 digitos
- [x] O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso, deverá retornar as temperaturas e formata-lás em: Celsius, Fahrenheit, Kelvin.
- [x] O sistema deve responder adequadamente nos seguintes cenários:
    - Em caso de sucesso:
        - [x] Código HTTP: 200
        - [x] Response Body: { "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }
    - Em caso de falha, caso o CEP não seja válido (com formato correto):
        - [x] Código HTTP: 422
        - [x] Mensagem: invalid zipcode
    - Em caso de falha, caso o CEP não seja encontrado:
        - [x] Código HTTP: 404
        - [x] Mensagem: can not find zipcode
- [x] Deverá ser realizado o [deploy no Google Cloud Run](https://goexpert-weather-api-t3n5zezwwq-uc.a.run.app).



## Requisitos: entrega
- [x] O código-fonte completo da implementação.
- [x] Testes automatizados demonstrando o funcionamento.
- [x] Utilize docker e docker-compose para que possamos realizar os testes de sua aplicação.
- [x] Deploy realizado no Google Cloud Run (free tier) e [endereço ativo para ser acessado](https://goexpert-weather-api-t3n5zezwwq-uc.a.run.app).
