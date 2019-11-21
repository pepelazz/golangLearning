package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pepelazz/golangLearning/uploadImage/uploadImage"
	"log"
	"net/http"
)

func main() {
	r := gin.New()

	// вырубаем CORS
	r.Use(LiberalCORS)
	r.Static("/stat-img", "./image")
	r.Static("/static", "./webClient/dist")
	r.Static("/statics", "./webClient/dist/statics")
	r.StaticFile("/", "./webClient/dist/index.html")

	r.POST("/upload_image", uploadImage.SimpleUploadImage)
	r.POST("/upload_image_resize", uploadImage.ResizeUploadImage)
	r.POST("/get_all_image", uploadImage.GetAllImage)

	// на ненайденный url отправляем статический файл для запуска vuejs приложения
	r.NoRoute(func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./webClient/dist/index.html")
	})

	err := r.Run(":3083")
	if err != nil {
		log.Fatalf("run webserver: %s", err)
	}
}

// LiberalCORS is a very allowing CORS middleware.
func LiberalCORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	if c.Request.Method == "OPTIONS" {
		if len(c.Request.Header["Access-Control-Request-Headers"]) > 0 {
			c.Header("Access-Control-Allow-Headers", c.Request.Header["Access-Control-Request-Headers"][0])
		}
		c.AbortWithStatus(http.StatusOK)
	}
}