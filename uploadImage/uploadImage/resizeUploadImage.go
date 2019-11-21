package uploadImage

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strings"
)

// методы API
func uploadImage(c *gin.Context, path string, width int, crop []int) {
	// извлекаем файл из парамeтров post запроса
	form, _ := c.MultipartForm()
	var fileName string
	imgExt := "jpeg"
	// берем первое имя файла из присланного списка
	for key := range form.File {
		fileName = key
		// извлекаем расширение файла
		arr := strings.Split(fileName, ".")
		if len(arr) > 1 {
			imgExt = arr[len(arr)-1]
		}
		continue
	}
	// извлекаем содержание присланного файла по названию файла
	file, _, err := c.Request.FormFile(fileName)
	if err != nil {
		httpError(c, http.StatusBadRequest, fmt.Sprintf("UploadXml c.Request.FormFile error: %s", err.Error()))
		return
	}
	defer file.Close()

	// перекодируем файл в картинку
	var img image.Image
	switch imgExt {
	case "jpeg":
		img, err = jpeg.Decode(file)
	case "jpg":
		img, err = jpeg.Decode(file)
	case "png":
		img, err = png.Decode(file)
	case "gif":
		img, err = gif.Decode(file)
	default:
		err = errors.New("Unsupported file type")
	}
	if err != nil {
		httpError(c, http.StatusBadRequest, fmt.Sprintf("uploadImage jpeg.Decode error: %s", err.Error()))
		return
	}
	fullFileName := fmt.Sprintf("%s.%s", randomFilename(), imgExt)
	// открываем файл для сохранения картинки
	fileOnDisk, err := os.Create(fmt.Sprintf("%s/%s", path, fullFileName))
	if err != nil {
		httpError(c, http.StatusBadRequest, fmt.Sprintf("uploadImage os.Create err: %s", err))
		return
	}
	defer fileOnDisk.Close()

	// если необходимо обрезать
	if crop != nil && len(crop) == 2 {
		analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
		topCrop, _ := analyzer.FindBestCrop(img, crop[0], crop[1])
		type SubImager interface {
			SubImage(r image.Rectangle) image.Image
		}
		img = img.(SubImager).SubImage(topCrop)
	}

	// сжатие размеров картинки до минимума - 500 или фактический размер
	imgWidth := uint(minInt(width, img.Bounds().Max.X))
	resizedImg := resize.Resize(imgWidth, 0, img, resize.Lanczos3)

	// сохранение файла
	err = jpeg.Encode(fileOnDisk, resizedImg, nil)
	if err != nil {
		httpError(c, http.StatusBadRequest, fmt.Sprintf("jpeg.Encode err: %s", err))
		return
	}

	// возвращаем ссылку на файл
	httpSuccess(c, map[string]string{"file": fmt.Sprintf("%s/%s", strings.Replace(path, IMAGE_DIR, STAT_IMAGE_PATH, 1), fullFileName)})
}

// вариант загрузки файла без преобразований
func ResizeUploadImage(c *gin.Context) {
	// извлекаем id продукта и создаем директорию
	productId, _ := c.GetPostForm("product_id")
	{
		if len(productId) == 0 {
			productId = "default"
		}
		err := os.MkdirAll(fmt.Sprintf("%s/product/%s", IMAGE_DIR, productId), os.ModePerm) // создаем директорию с id продукта, если еще не создана
		if err != nil {
			httpError(c, http.StatusBadRequest, fmt.Sprintf("uploadImage os.MkdirAll error: %s", err))
			return
		}
	}
	path := fmt.Sprintf("%s/product/%s", IMAGE_DIR, productId)
	uploadImage(c, path,800, []int{10, 10})
}