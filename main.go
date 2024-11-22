package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	returnViaCep, err := BuscaViaCep(cep)
	if err != nil {
		panic(err)
	}
	fmt.Println("Retorno ViaCEP", returnViaCep)

	returnBuscaBrasilAPI, err := BuscaBrasilAPI(cep)
	if err != nil {
		panic(err)
	}
	fmt.Println("Retorno BrasilAPI", returnBuscaBrasilAPI)
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

func BuscaViaCep(cep string) (*ViaCEP, error) {
	url := "https://viacep.com.br/ws/" + cep + "/json/"
	var data ViaCEP
	err := fetchCepData(url, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func BuscaBrasilAPI(cep string) (*BrasilAPI, error) {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep
	var data BrasilAPI
	err := fetchCepData(url, &data)
	if err != nil {
		return nil, err
	}

	data.CepData = CepData{
		Cep:        data.Cep,
		Logradouro: data.Street,
		Bairro:     data.Neighborhood,
		Localidade: data.City,
		Uf:         data.State,
	}
	return &data, nil
}
