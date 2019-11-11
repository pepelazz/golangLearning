package main

import "fmt"

// удаляем файлы кроме двух последних
func removeOldBackupsOnServer(prj Project) error {
	path := "backups/" + prj.Name
	res, err := getResource(path)
	if err != nil {
		return err
	}
	for i, v := range res.Embedded.Items {
		fmt.Printf("%v %s %s\n", i, v.Name, v.Path)
		if i > 1 {
			err = deleteFile(v.Path)
			if err != nil {
				fmt.Printf("deleteFile err: %s\n", err)
			}
		}
	}
	return nil
}
