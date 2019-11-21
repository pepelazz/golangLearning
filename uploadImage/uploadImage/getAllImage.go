package uploadImage

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetAllImage(c *gin.Context) {
	type Photo struct {
		Url  string `json:"url"`
		Size int64  `json:"size"`
	}
	photoList := []Photo{}
	err := filepath.Walk("./image",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				// заменяем название директории image на stat-img для того чтобы корректно отрабатывал раутинг
				url := "stat-img" + strings.TrimPrefix(path, "image")
				// для версии windows заменяем слэыши в пути файла на обратные
				if runtime.GOOS == "windows" {
					url = strings.Replace(url, "\\", "/", -1)
				}
				photoList = append(photoList, Photo{url, info.Size()})
			}
			return nil
		})
	if err != nil {
		fmt.Printf("err %s\n", err)
	}
	httpSuccess(c, photoList)
}
