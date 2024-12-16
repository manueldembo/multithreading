package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Como não se abordou nada sobre o mecanismo de entrada do CEP, preferi não fazer coisas que não foram
// solicitadas no desáfio de modo a não comprometer o resultado da avalição. Obrigado!

func main() {
	chBrasilAPI := make(chan map[string]interface{})
	chViaCep := make(chan map[string]interface{})
	cep := "01153000"
	brasilApiUrl := "https://brasilapi.com.br/api/cep/v1/" + cep
	viaCep := "https://viacep.com.br/ws/" + cep + "/json/"

	go func() {
		result, err := getCep(brasilApiUrl)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		chBrasilAPI <- map[string]interface{}{"brasilapi.com.br": result}
	}()

	go func() {
		result, err := getCep(viaCep)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		chViaCep <- map[string]interface{}{"viacep.com.br": result}
	}()

	select {
	case msg := <-chBrasilAPI:
		for k, v := range msg {
			fmt.Printf("Source: %v\nAnswer: %+v\n", k, v)
		}
	case msg := <-chViaCep:
		for k, v := range msg {
			fmt.Printf("Source: %v\nAnswer: %+v\n", k, v)
		}
	case <-time.After(time.Second):
		println("Timeout reached!")
	}
}

func getCep(url string) (interface{}, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result interface{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
