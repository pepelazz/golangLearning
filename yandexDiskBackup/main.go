package main

import (
	"fmt"
	"time"
)

type (
	Project struct {
		Name         string
		Type         string
		Path         string
		DockerPgName string
		ExistFile    string
	}
	Resource struct {
		Name       string `json:"name"`
		Path       string `json:"path"`
		Created    string `json:"created"`
		ResourceId string `json:"resource_id"`
		Type       string `json:"type"`
		MimeType   string `json:"mime_type"`
		Embedded   struct {
			Items []Resource `json:"items"`
			Path  string     `json:"path"`
		} `json:"_embedded"`
	}
)

var config Config
var ynxUrl = "https://cloud-api.yandex.net/v1/disk"

func main() {
	// читаем инфу из конфига
	readConfig()

	for _, prj := range config.ProjectList {
		// создаем папку на яндексе для проекта
		err := createFolder("backups/" + prj.Name)
		if err != nil {
			// ошибку считаем не критичной, потому что если папка уже есть на диске, то буде сообщение об ошибке. Но на всякий случай выводим в консоль
			fmt.Printf("error createFolder %s", prj.Name)
		}
		// получаем адрес файла с бэкапом
		fileName, err := getBackupFile(prj)
		if err != nil {
			fmt.Printf("error getBackupFile %s err:%s", fileName, err)
			continue
		}
		// копируем файл на сервер
		fullPath := prj.Path + "/" + fileName
		// в случае если идет обработка уже существующего файла, то fullPath будет соответствовать тому пути, который указан в конфиге
		if len(prj.ExistFile) > 0 {
			fullPath = prj.ExistFile
		}
		err = uploadFile(fullPath, "backups/"+prj.Name+"/"+fileName)
		if err != nil {
			fmt.Printf("error uploadFile %s err:%s", fullPath, err)
			continue
		}
		fmt.Printf("appName: %s uploaded to Yandex Disk\n", prj.Name)

		// в случае если файл бэкапа был создан программно - удаляем файл на сервере с таймаутом
		if len(prj.ExistFile) == 0 {
			time.Sleep(1 * time.Minute)
			removeBackupFile(fullPath)

			// удаляем старые файлы на яндекс диске
			err = removeOldBackupsOnServer(prj)
			if err != nil {
				fmt.Printf("error removeBackupFile %s err:%s", prj.Name, err)
				continue
			}
		}
	}
}
