package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(LiberalCORS)
	r.GET("/ping", pong)
	r.GET("/user", getUser)
	r.POST("/user", postUser)
	r.POST("/getSymbols", getSymbols)
	r.POST("/getQuotes", getQuotes)
	r.Run(":9090")
}

var apiKey = "0mq7fJfyh6nFRi2frX0eL1h8OADmQMht"

func getSymbols(c *gin.Context) {
	symbol_list := GetSymbols(apiKey)
	c.JSON(200, gin.H{
		"symbols": symbol_list,
	})
}

func getQuotes(c *gin.Context) {

	s := struct {
		Symbol string `json:"symbol"`
	}{}

	c.BindJSON(&s)

	if len(s.Symbol) > 0 {
		// symbols := []string{"EURUSD", "AUDJPY", "GBPCHF"}
		symbols := []string{s.Symbol}
		quotes := GetQuotes(symbols, apiKey)

		for i := range quotes {
			quotes[i].Bid = quotes[i].Bid - 0.02
			quotes[i].Ask = quotes[i].Ask + 0.02
		}

		c.JSON(200, gin.H{
			"quotes": quotes,
		})
	}
}

func pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func getUser(c *gin.Context) {
	id := c.Query("id")
	user := struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}{id, "Вася"}

	c.JSON(200, gin.H{
		"user": user,
	})
}

func postUser(c *gin.Context) {
	user := struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}{}
	c.BindJSON(&user)

	user.Name = "Петя"

	c.JSON(200, gin.H{
		"user": user,
	})
}

func LiberalCORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	if c.Request.Method == "OPTIONS" {
		if len(c.Request.Header["Access-Control-Request-Headers"]) > 0 {
			c.Header("Access-Control-Allow-Headers", c.Request.Header["Access-Control-Request-Headers"][0])
		}
		c.AbortWithStatus(http.StatusOK)
	}
}
