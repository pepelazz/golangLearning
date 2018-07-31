package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"strconv"
)

func startWebServer() {
	r := gin.New()

	// вырубаем CORS
	r.Use(LiberalCORS)

	r.POST("/api/userListGet", apiGetUserList)
	r.POST("/api/userListSave", apiSaveUserList)
	r.POST("/api/userRemove", apiUserRemove)

	// на ненайденный url отправляем статический файл для запуска vuejs приложения
	r.NoRoute(func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./index.html")
	})

	r.Run(fmt.Sprintf(":%v", HTTP_PORT))
}

// методы API
func apiSaveUserList(c *gin.Context) {
	res := []User{}
	// извлекаем json-параметры запроса
	if err := c.BindJSON(&res); err != nil {
		httpError(c, http.StatusOK, fmt.Sprintf("post json params error: %s\n", err))
		return
	}

	//userList = res

	for _, u := range res {
		err := u.SaveToDB()
		if err != nil {
			fmt.Printf("u.SaveToDB err: %s\n", err)
		}
	}

	httpSuccess(c, nil)
}

func apiGetUserList(c *gin.Context) {

	userList, err := getUserListFromDB(ENCODE_TYPE_GOB)
	if err != nil {
		fmt.Printf("err %s\n", err)
		httpError(c, http.StatusOK, fmt.Sprintf("user list encoding error: %s\n", err))
		return
	}

	userList[0].UserAccount.GetAvatar()

	httpSuccess(c, userList)
}

func apiUserRemove(c *gin.Context) {
	res := struct {
		Id int `json:"id"`
	}{}

	if err := c.BindJSON(&res); err != nil {
		httpError(c, http.StatusOK, fmt.Sprintf("post json params error: %s\n", err))
		return
	}

	localDb.Delete(USER_DB_BUCKET, strconv.Itoa(res.Id))

	httpSuccess(c, nil)
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

func httpError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"ok":      false,
		"message": message,
	})
	c.Abort()
}

func httpSuccess(c *gin.Context, res interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"result": res,
	})
}
