package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
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

	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatal("Ошибка открытия файла базы данных:", err)
	}
	defer db.Close()

	log.Println("Соединение с базой данных установлено.")

	createTable := `
	CREATE TABLE IF NOT EXISTS items (
		indicator_id TEXT,
		indicator_value TEXT,
		country_id TEXT,
		country_value TEXT,
		country_iso3_code TEXT,
		date TEXT PRIMARY KEY NOT NULL UNIQUE,
		value INTEGER,
		unit TEXT,
		obs_status TEXT,
		decimal INTEGER,
		created_at TEXT DEFAULT (DATETIME('now'))
	)
	`

	if _, err := db.Exec(createTable); err != nil {
		log.Fatal("Ошибка создания таблицы:", err)
	}

	log.Println("Таблица для хранения данных создана.")

	insertStmt := `
	INSERT OR REPLACE INTO items (indicator_id, indicator_value, country_id, country_value, country_iso3_code, date, value, unit, obs_status, decimal)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	for _, item := range items {
		_, err := db.Exec(insertStmt,
			item.Indicator.Id, item.Indicator.Value,
			item.Country.Id, item.Country.Value,
			item.CountryISO3Code, item.Date,
			item.Value, item.Unit,
			item.ObsStatus, item.Decimal,
		)
		if err != nil {
			log.Fatal("Ошибка вставки объекта в базу данных:", err)
		}
	}

	log.Printf("Число сохраненных записей: %d\n", len(items))
	log.Println("Программа завершает работу.")
}
