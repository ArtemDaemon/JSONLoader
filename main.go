package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type Item struct {
	Indicator       Subitem `json:"indicator"`
	Country         Subitem `json:"country"`
	CountryISO3Code string  `json:"countryiso3code"`
	Date            string  `json:"date"`
	Value           uint    `json:"value"`
	Unit            string  `json:"unit"`
	ObsStatus       string  `json:"obs_status"`
	Decimal         uint    `json:"decimal"`
}

type Subitem struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

func main() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Ошибка открытия лог-файла:", err)
	}
	log.SetOutput(file)

	log.Println("Запуск утилиты...")

	resp, err := http.Get("https://api.worldbank.org/v2/countries/USA/indicators/SP.POP.TOTL?per_page=5000&format=json")
	if err != nil {
		log.Fatal("Ошибка получения JSON-файла:", err)
	}
	defer resp.Body.Close()
	log.Println("JSON успешно получен.")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Ошибка чтения JSON-файла:", err)
	}

	var topLevel []json.RawMessage
	if err := json.Unmarshal(body, &topLevel); err != nil {
		log.Fatal("Ошибка парсинга JSON-файла:", err)
	}

	if len(topLevel) < 2 {
		log.Fatal("Ожидалось минимум 2 элемента в JSON-массиве")
	}

	// Пропускаем первый элемент (summary)
	var items []Item
	if err := json.Unmarshal(topLevel[1], &items); err != nil {
		log.Fatal("Ошибка парсинга JSON-файла:", err)
	}

	log.Println("JSON успешно распарсен.")
}
