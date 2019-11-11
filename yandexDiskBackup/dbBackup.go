package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func getBackupFile(prj Project) (string, error) {
	// вариант, когда файл с бэкапом базы уже есть на диске
	if len(prj.ExistFile) > 0 {
		// извлекаем имя файла из полного пути
		arr := strings.Split(prj.ExistFile, "/")
		return arr[len(arr)-1], nil
	}
	// далее вариант создания бэкапа базы из докер контейнера
	// формируем название файла с учетом даты создания бэкапа
	fileName := fmt.Sprintf("%s_dump_%s.zip", prj.Name, time.Now().Format("2006_01_02"))
	// формируем и исполняем команду в bush
	cmd := exec.Command("sh", "-c", strings.Join([]string{
		"cd " + prj.Path,
		fmt.Sprintf("docker exec -t %s pg_dumpall -c -U postgres  > %s_dump", prj.DockerPgName, prj.Name),
		fmt.Sprintf("zip %s %s_dump", fileName, prj.Name), // архивируем бэкап
		fmt.Sprintf("rm %s_dump", prj.Name),               // удаляем бэкап
	}, ";"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// проверяем что в результате выполнения команды в bash на диске появился нужный нам файл
	fullPath := fmt.Sprintf("%s/%s", prj.Path, fileName)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", errors.New(fmt.Sprintf("file %s not created", fullPath))
	}
	return fileName, nil
}

func removeBackupFile(path string) {
	cmd := exec.Command("sh", "-c", strings.Join([]string{
		fmt.Sprintf("rm %s", path),
	}, ";"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("removeBackupFile err %s\n", err)
	}
}
