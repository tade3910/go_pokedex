package pokemap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	Next     string
	Previous string
}

type jsonResult struct {
	Count    int
	Next     *string
	Previous *string
	Results  []struct {
		Name string
		Url  string
	}
}

func (currConfig *Config) GetMap(prev bool) (jsonResult, error) {
	errorJson, result := jsonResult{}, jsonResult{}
	url := currConfig.Next
	if prev {
		if currConfig.Previous == "" {
			return errorJson, errors.New("there is no previous map")
		}
		url = currConfig.Previous
	}
	res, err := http.Get(url)
	if err != nil {
		return errorJson, err
	}
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		errorMessage := fmt.Sprintf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		return errorJson, errors.New(errorMessage)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return errorJson, err
	}
	return result, nil
}

func GetNewConfig() Config {
	return Config{
		Next:     "https://pokeapi.co/api/v2/location/",
		Previous: "",
	}
}
