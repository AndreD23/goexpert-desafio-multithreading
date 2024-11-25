package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CepData struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ViaCEP struct {
	CepData
}

type BrasilAPI struct {
	CepData
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

func main() {
	cep := "05187010"

	cepViaCep := make(chan *CepData)
	cepBrasilAPI := make(chan *CepData)

	go BuscaViaCep(cep, cepViaCep)
	go BuscaBrasilAPI(cep, cepBrasilAPI)

	select {
	case mgs1 := <-cepViaCep:
		fmt.Println("ViaCEP", mgs1)
	case msg2 := <-cepBrasilAPI:
		fmt.Println("BrasilAPI", msg2)
	case <-time.After(time.Second * 1):
		fmt.Println("Timeout")
	}
}

func fetchCepData(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(res, target)
	if err != nil {
		return err
	}

	return nil
}

func BuscaViaCep(cep string, cepData chan<- *CepData) {
	url := "https://viacep.com.br/ws/" + cep + "/json/"
	var data ViaCEP
	err := fetchCepData(url, &data)
	if err != nil {
		panic(err)
	}
	cepData <- &data.CepData
}

func BuscaBrasilAPI(cep string, cepData chan<- *CepData) {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep
	var data BrasilAPI
	err := fetchCepData(url, &data)
	if err != nil {
		panic(err)
	}

	data.CepData = CepData{
		Cep:        data.Cep,
		Logradouro: data.Street,
		Bairro:     data.Neighborhood,
		Localidade: data.City,
		Uf:         data.State,
	}

	cepData <- &data.CepData
}
