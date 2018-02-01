package main

import (
	"fmt"
	"net/http"
	"html/template"
	"log"
	"math/rand"
	"time"
	"strings"
	"os/exec"
	"runtime"
)

type (
	// структура описывающая пару слов: ангийский и русский варианты
	Word struct {
		En string
		Ru string
	}

	// структура, в которую буждем сохранять ответ пользователя
	Answer struct {
		Ru          string // ответ пользователя
		En          string // английское слово, которое надо перевести
		IsRight     bool   // верный или неверный ответ
		RightAnswer string // правильный ответ
	}

	// структру данных для html шаблона
	DataForTemplate struct {
		CurrentWord        Word     // текущее слово
		UserAnswers        []Answer // список ответов пользователя
		RightCnt, WrongCnt int      // счетчики правильных и неправильных ответов
	}
)

var (
	// массив, в который будем накапливать ответы пользователя
	userAnswers = []Answer{}
	// набор html шаблонов, которые загружаются из директории static
	tmpl = template.Must(template.ParseGlob("static/*.html"))
	// словарь: массив с парами слов
	dictonary = []Word{}
	// текущее слово. Объявляем переменную глобально, чтобы можно было иметь доступ из разных функций
	currentWord Word
	//счетчики правильных и неправильных ответов
	rightCnt, wrongCnt int
)

// точка запуска программы
func main() {
	// создаем генератор случайных чисел при запуске. Если этого не сделать, то rand будет выдавать всегда одну и ту же последовательность
	rand.Seed(time.Now().UnixNano())

	// читаем данные из excel файла - заполняем массив слов в словаре
	readDictionary()

	// прописываем роутинг: url и функция, которая вызывается при переходе на этот url
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/enRu", enRu)

	// открываем приложение в браузере
	go openBrowser("http://localhost:9090")

	// стартуем сервер на порту 9090. Открыть можно localhost:9090
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// рендеринг главной страницы
func mainPage(w http.ResponseWriter, r *http.Request) {
	// просто отображаем страницу index.html без всяких данных
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

// рендеринг страницы с тренажером
func enRu(w http.ResponseWriter, r *http.Request) {
	// вариант когда пришел запрос типа POST. Это происходит в том случае если пользваотель написал ответ и нажал кнопку "отправить" или enter на клавиатуре
	if r.Method == "POST" {
		// обрабатываем форму с ответом
		r.ParseForm()
		// извлекаем значение поля "answer"
		userAnswer := r.Form.Get("answer")

		// проверяем ответ пользователя.
		// В результате выполнения функции получаем объект Answer c заполненными полями и проставленным занчением - правильный ответ или нет
		answer := checkUserAnswer(userAnswer)
		// добавляем полученный объект Answer в массив с ответами.
		// Добавляем его в начало массива - для этого создаем массив из одного элемента и объекдиням его с существующим массивом ответов. При добавлении массива к массиву нужно ставить ...
		userAnswers = append([]Answer{answer}, userAnswers...)

	}

	// формируем новый вопрос
	{
		// случайным образом выбираем индекс в диапозоне [0, длина словаря-1]
		index := rand.Intn(len(dictonary) - 1)
		// по индексу находим слово в массиве
		currentWord = dictonary[index]
	}

	// формируем данные для html шаблона. Передаем туда текущее выбранное слово, массив ответов и итоговые счетчики
	data := DataForTemplate{CurrentWord: currentWord, UserAnswers: userAnswers, RightCnt: rightCnt, WrongCnt: wrongCnt}
	// рендерим html страницу - передавая в объект http.ResponseWriter имя html шаблона (в данном случае "index.html") и объект с данными
	err := tmpl.ExecuteTemplate(w, "enRu.html", data)
	if err != nil {
		fmt.Println(err)
	}
}

// функция проверки ответа пользователя
func checkUserAnswer(userAnswer string) Answer {
	// создаем объект типа Answer и сразу при создании заполняем его поля, кроме поля IsRight
	res := Answer{Ru: userAnswer, En: currentWord.En, RightAnswer: currentWord.Ru}
	// проверяем содержит ли ответ пользователя правильный вариант
	if strings.Contains(userAnswer, currentWord.Ru) {
		// увеличиваем счетчик правильных ответов
		rightCnt++
		// заполняем поле IsRight
		res.IsRight = true
	} else {
		// увеличиваем счетчик неверных ответов
		wrongCnt++
		// заполняем поле IsRight
		res.IsRight = false
	}

	return res
}

// openBrowser tries to open the URL in a browser,
// and returns whether it succeed in doing so.
func openBrowser(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
