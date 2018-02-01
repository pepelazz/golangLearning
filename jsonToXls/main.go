package main

import (
	"io/ioutil"
	"fmt"
	"log"
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize"
)

type (
	Route struct {
		From  string      `json:"from"`
		To    string      `json:"to"`
		Car   string      `json:"car"`
		Cargo []CargoItem `json:"cargo"`
		Date  string      `json:"date"`
	}

	CargoItem struct {
		Product string `json:"product"`
		Amount  int    `json:"amount"`
	}
)

func main() {

	// читаем исходные данные из json файла
	file, err := ioutil.ReadFile("./routes.json")
	if err != nil {
		log.Fatal("File error: ", err)
	}

	// перекладываем данные из json в массив структур Route
	var routes []Route
	json.Unmarshal(file, &routes)

	//TODO: произвести необходимую трансофрмацию

	// формируем отчет. Например, только грузы из Москвы
	//Записываем данные в excel файл

	// создаем новый excel файл
	xlsx := excelize.NewFile()
	// Создаем новый лист
	xlsx.NewSheet("report")
	xlsx.SetCellValue("report", "A1", "Откуда")
	xlsx.SetCellValue("report", "B1", "Куда")
	xlsx.SetCellValue("report", "C1", "Продукт")
	xlsx.SetCellValue("report", "D1", "Количество")
	i := 2
	for _, v := range routes {
		// если не из Москвы, то пропускаем обработку и переходим к следующей записи
		if v.From != "Москва" {
			continue
		}
		for _, c := range v.Cargo {
			xlsx.SetCellValue("report", fmt.Sprintf("A%v", i), v.From)
			xlsx.SetCellValue("report", fmt.Sprintf("B%v", i), v.To)
			xlsx.SetCellValue("report", fmt.Sprintf("C%v", i), c.Product)
			xlsx.SetCellValue("report", fmt.Sprintf("D%v", i), c.Amount)
			i++
		}
	}
	// Сохраняем файл
	err = xlsx.SaveAs("./report.xlsx")
	if err != nil {
		fmt.Println(err)
	}

}
