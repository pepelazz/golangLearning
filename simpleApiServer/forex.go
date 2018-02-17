/*
 * This library is provided without warranty under the MIT license
 * Created by Jacob Davis <jacob@1forge.com>
 */

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ConversionResult struct {
	Value     float32
	Text      string
	Timestamp int
}

type MarketStatus struct {
	MarketIsOpen bool `json:"market_is_open"`
}

type Quote struct {
	Symbol    string
	Bid       float32
	Ask       float32
	Price     float32
	Timestamp int
}

func fetch(query string, api_key string) []byte {
	response, e := http.Get("http://forex.1forge.com/1.0.3/" + query + "&api_key=" + api_key)

	if e != nil {
		log.Fatal(e)
	}

	defer response.Body.Close()
	body, e := ioutil.ReadAll(response.Body)

	if e != nil {
		log.Fatal(e)
	}

	return body
}

func GetSymbols(api_key string) []string {
	result := fetch("symbols?cache=false", api_key)

	symbol_list := []string{}

	e := json.Unmarshal(result, &symbol_list)

	if e != nil {
		log.Fatal(e)
	}

	return symbol_list
}

func GetQuotes(symbols []string, api_key string) []Quote {
	result := fetch("quotes?pairs="+strings.Join(symbols, ","), api_key)

	quotes := []Quote{}

	e := json.Unmarshal(result, &quotes)

	if e != nil {
		log.Fatal(e)
	}

	return quotes
}

func Convert(from string, to string, quantity int, api_key string) ConversionResult {

	result := fetch("convert?from="+from+"&to="+to+"&quantity="+strconv.Itoa(quantity), api_key)

	conversion_result := ConversionResult{}

	e := json.Unmarshal(result, &conversion_result)

	if e != nil {
		log.Fatal(e)
	}

	return conversion_result
}

func MarketIsOpen(api_key string) bool {
	result := fetch("market_status?cache=false", api_key)

	market_status := MarketStatus{}

	e := json.Unmarshal(result, &market_status)

	if e != nil {
		log.Fatal(e)
	}

	return market_status.MarketIsOpen
}
