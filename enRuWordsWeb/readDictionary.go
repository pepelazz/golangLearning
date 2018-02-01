package main

import (
	"github.com/xuri/excelize"
	"log"
)

// функция чтения словаря из excel файла
func readDictionary() {
	xlsx, err := excelize.OpenFile("./data/dict.xlsx")
	if err != nil {
		log.Fatalf("Open file: %s", err)
	}

	rows := xlsx.GetRows("Sheet1")
	for i, row := range rows {
		// первую строку таблицы excel пропускаем (там заголовки столбцов)
		if i == 0 {
			continue
		}
		// если в строке болше одного элемента и первый элемент не равен нулю, то считаем что в строке есть пара слов
		if len(row) > 1 && len(row) > 1 {
			// добавляем в массив вновьсозданную структуру Word, с заполнеными значениями из excel строки
			dictonary = append(dictonary, Word{En: row[0], Ru: row[1]})
		}
	}
}