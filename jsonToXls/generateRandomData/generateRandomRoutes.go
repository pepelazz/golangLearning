package generateRandomData

import (
	"time"
	"math/rand"
	"github.com/Pallinder/go-randomdata"
	"github.com/icrowley/fake"
	"github.com/pepelazz/go-bot-utils"
	"io/ioutil"
	"encoding/json"
)

type City string
type Commodity string
type Car string

type Route struct {
	From City `json:"from"`
	To   City `json:"to"`
}

type CargoUnit struct {
	Product string `json:"product"`
	Amount  int    `json:"amount"`
}

type Tour struct {
	Route
	Car  `json:"car"`
	Cargo []CargoUnit `json:"cargo"`
	Date  string `json:"date"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	err := fake.SetLang("ru")
	goBotUtils.CheckErr(err, "fake.SetLang")
	res, err := json.Marshal(generateRoutes())
	goBotUtils.CheckErr(err, "json.Marshal")

	err = ioutil.WriteFile("routes.json", res, 0644)
	goBotUtils.CheckErr(err, "Write file")

}

func generateRoutes() (res []Tour) {

	cityList1 := []City{"Москва", "Санкт-Петербург"}
	cityList2 := []City{"Тула", "Воронеж", "Иваново", "Владимир", "Брянск", "Смоленск", "Тверь", "Калуга", "Кастрома", "Чебоксары", "Курск", "Орел"}

	routes := []Route{}

	productList := make([]string, 50)
	for i := range productList {
		productList[i] = fake.ProductName()
	}

	for _, v1 := range cityList1 {
		for _, v2 := range cityList2 {
			routes = append(routes, Route{v1, v2}, Route{v1, v2})
		}
	}

	amountList := []int{20, 30, 45, 55, 15}

	i := 0
	for {
		tour := Tour{}
		tour.Route = routes[rand.Intn(len(routes)-1)]
		tour.Cargo = make([]CargoUnit, rand.Intn(3) + 2)
		for i := range tour.Cargo {
			tour.Cargo[i] = CargoUnit{productList[rand.Intn(len(productList)-1)], amountList[rand.Intn(len(amountList)-1)]}
		}
		tour.Car = Car(randomdata.StringNumberExt(2, "-", 3))
		tour.Date = randate().Format("02-01-2006")

		res = append(res, tour)
		i++
		if i == 1000 {
			break
		}
	}

	return
}


func randate() time.Time {
	min := time.Date(2017, 10, 1, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2018, 1, 31, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min

	return time.Unix(sec, 0)
}