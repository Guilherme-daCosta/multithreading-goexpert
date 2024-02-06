package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BrasilAPIResponse struct {
	URL    string    `json:"url"`
	Result BrasilAPI `json:"result"`
}

type BrasilAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCEPResponse struct {
	URL    string `json:"url"`
	Result ViaCEP `json:"result"`
}

type ViaCEP struct {
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

func main() {
	resultBrasilAPI := make(chan BrasilAPIResponse, 2)
	resultViaCEP := make(chan ViaCEPResponse, 2)
	CEP := "01153000"

	go func() {
		url := "https://brasilapi.com.br/api/cep/v1/" + CEP
		resp, err := http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		body := BrasilAPI{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		if err != nil {
			return
		}

		response := BrasilAPIResponse{
			URL:    url,
			Result: body,
		}

		resultBrasilAPI <- response
	}()

	go func() {
		url := "http://viacep.com.br/ws/" + CEP + "/json/"
		resp, err := http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		body := ViaCEP{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		if err != nil {
			return
		}

		response := ViaCEPResponse{
			URL:    url,
			Result: body,
		}

		resultViaCEP <- response
	}()

	select {
	case res := <-resultBrasilAPI:
		fmt.Printf("BrasilAPI: %+v\n", res)
	case res := <-resultViaCEP:
		fmt.Printf("ViaCEP: %+v\n", res)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout")
	}
}
