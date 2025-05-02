package main

import (
	"log"
	"os"
)

func main() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Ошибка открытия лог-файла:", err)
	}
	log.SetOutput(file)

	log.Println("Запуск утилиты...")
}
