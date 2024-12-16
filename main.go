package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

// Como não se abordou nada sobre o mecanismo de entrada do CEP, preferi não fazer coisas que não foram
// solicitadas no desáfio de modo a não comprometer o resultado da avalição. Obrigado!

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

func showResponse(ch1, ch2 chan map[string]interface{}) {
	select {
	case msg := <-ch1:
		for k, v := range msg {
			fmt.Printf("Source: %v\nResponse: %+v\n", k, v)
		}
	case msg := <-ch2:
		for k, v := range msg {
			fmt.Printf("Source: %v\nResponse: %+v\n", k, v)
		}
	case <-time.After(time.Second):
		println("Timeout reached!")
	}
}

func execute(ch chan map[string]interface{}, url string) {
	result, err := getCep(url)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	re := regexp.MustCompile(`//([^/]+\.br)`)
	match := re.FindStringSubmatch(url)
	ch <- map[string]interface{}{match[1]: result}
}

func main() {
	chBrasilAPI := make(chan map[string]interface{})
	chViaCep := make(chan map[string]interface{})
	cep := "01153000"
	brasilApiUrl := "https://brasilapi.com.br/api/cep/v1/" + cep
	viaCepUrl := "https://viacep.com.br/ws/" + cep + "/json/"

	go execute(chBrasilAPI, brasilApiUrl)
	go execute(chViaCep, viaCepUrl)

	showResponse(chBrasilAPI, chViaCep)
}
