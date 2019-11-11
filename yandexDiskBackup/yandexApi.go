package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func apiRequest(path, method string) (*http.Response, error) {
	client := http.Client{}
	url := fmt.Sprintf("%s/%s", ynxUrl, path)
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("OAuth %s", config.AuthToken))
	return client.Do(req)
}

// создание директории
func createFolder(path string) error {
	_, err := apiRequest(fmt.Sprintf("resources?path=%s", path), "PUT")
	return err
}

// загрузка файла
func uploadFile(localPath, remotePath string) error {
	// функция получения url для загрузки файла
	getUploadUrl := func(path string) (string, error) {
		res, err := apiRequest(fmt.Sprintf("resources/upload?path=%s&overwrite=true", path), "GET")
		if err != nil {
			return "", err
		}
		var resultJson struct {
			Href string `json:"href"`
		}
		err = json.NewDecoder(res.Body).Decode(&resultJson)
		if err != nil {
			return "", err
		}
		return resultJson.Href, err
	}

	// читаем локальный файл с диска
	data, err := os.Open(localPath)
	if err != nil {
		return err
	}
	// получем ссылку для загрузки файла
	href, err := getUploadUrl(remotePath)
	if err != nil {
		return err
	}
	defer data.Close()
	// загружаем файл по полученной ссылке методом PUT
	req, err := http.NewRequest("PUT", href, data)
	if err != nil {
		return err
	}
	// в header запроса добавляем токен
	req.Header.Add("Authorization", fmt.Sprintf("OAuth %s", config.AuthToken))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

// удаление файла
func deleteFile(path string) error {
	_, err := apiRequest(fmt.Sprintf("resources?path=%s&permanently=true", path), "DELETE")
	return err
}

// получение содержимого директории
func getResource(path string) (*Resource, error) {
	res, err := apiRequest(fmt.Sprintf("resources?path=%s&limit=50&sort=-created", path), "GET")
	if err != nil {
		return nil, err
	}

	var result *Resource
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
